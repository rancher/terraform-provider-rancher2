package rancher2

import (
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	corev1 "k8s.io/api/core/v1"
)

const oidcSecretsNamespace = "cattle-oidc-client-secrets"

func resourceRancher2OIDCClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2OIDCClientCreate,
		Read:   resourceRancher2OIDCClientRead,
		Update: resourceRancher2OIDCClientUpdate,
		Delete: resourceRancher2OIDCClientDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2OIDCClientImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: oidcClientFields(),
	}
}

func oidcClientFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"token_expiration_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The duration (in seconds) before an access token and ID token expire.",
		},
		"refresh_token_expiration_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The duration (in seconds) a refresh token remains valid before it expires.",
		},
		"redirect_uris": {
			Description: "List of allowed redirect_uris for this client.",
			Required:    true,
			Type:        schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "OIDCClient description",
		},
		"client_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The Client ID for OIDC Access.",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The Client Secret for OIDC Access.",
		},
	}

	maps.Copy(s, commonAnnotationLabelFields())

	return s
}

func resourceRancher2OIDCClientCreate(d *schema.ResourceData, meta any) error {
	log.Printf("[INFO] Creating OIDCClient")
	oidcClient, err := expandOIDCClient(d)
	if err != nil {
		return err
	}

	client, err := meta.(managementClientGetter).ManagementClient()
	if err != nil {
		return fmt.Errorf("getting a client when creating OIDC Client: %w", err)
	}

	createdClient, err := client.OIDCClient.Create(oidcClient)
	if err != nil {
		return fmt.Errorf("creating OIDCClient: %w", err)
	}

	log.Printf("[INFO] Created OIDCClient ID %s", createdClient.ID)

	d.SetId(createdClient.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"secret created"},
		Refresh:    oidcClientStateRefreshFunc(client, createdClient.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("[ERROR] waiting for OIDCClient (%s) to be ready: %w", createdClient.ID, err)
	}

	return resourceRancher2OIDCClientRead(d, meta)
}

func oidcClientStateRefreshFunc(client *managementClient.Client, clientID string) resource.StateRefreshFunc {
	return func() (any, string, error) {
		oidcClient, err := client.OIDCClient.ByID(clientID)
		if err != nil {
			if IsNotFound(err) {
				return nil, "", fmt.Errorf("reading OIDC Client %s: %w", clientID, err)
			}
			return nil, "reading client", fmt.Errorf("reading OIDC Client %s: %w", clientID, err)
		}

		if oidcClient.Status.ClientSecrets == nil || len(oidcClient.Status.ClientSecrets) == 0 {
			return oidcClient, "pending", nil
		}

		return oidcClient, "secret created", nil
	}
}

func resourceRancher2OIDCClientRead(d *schema.ResourceData, meta any) error {
	log.Printf("[INFO] Refreshing OIDCClient ID %s", d.Id())

	client, err := meta.(managementClientGetter).ManagementClient()
	if err != nil {
		return fmt.Errorf("getting a client when reading OIDC Client: %w", err)
	}

	oidcClient, err := client.OIDCClient.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading OIDC Client %s: %w", d.Id(), err)
	}

	var clientSecret *corev1.Secret
	if oidcClient.Status.ClientID != "" {
		clientSecret, err = meta.(secretFetcher).SecretByName(rancher2DefaultLocalClusterID, oidcSecretsNamespace, oidcClient.Status.ClientID)
		if err != nil && !IsNotFound(err) && !IsForbidden(err) {
			return fmt.Errorf("[ERROR] getting OIDC Client Secret: %s: %w", oidcClient.Status.ClientID, err)
		}
	}

	return flattenOIDCClient(d, oidcClient, clientSecret)
}

func resourceRancher2OIDCClientUpdate(d *schema.ResourceData, meta any) error {
	log.Printf("[INFO] Updating OIDCClient ID %s", d.Id())

	client, err := meta.(managementClientGetter).ManagementClient()
	if err != nil {
		return fmt.Errorf("getting a client when updating OIDC Client: %w", err)
	}

	oidcClient, err := client.OIDCClient.ByID(d.Id())
	if err != nil {
		return fmt.Errorf("getting OIDC Client %s for update: %w", d.Id(), err)
	}

	update := map[string]any{
		"labels":                        toMapString(d.Get("labels").(map[string]any)),
		"annotations":                   toMapString(d.Get("annotations").(map[string]any)),
		"description":                   d.Get("description"),
		"redirectURIs":                  toArrayString(d.Get("redirect_uris").([]any)),
		"tokenExpirationSeconds":        d.Get("token_expiration_seconds"),
		"refreshTokenExpirationSeconds": d.Get("refresh_token_expiration_seconds"),
	}

	_, err = client.OIDCClient.Update(oidcClient, update)
	if err != nil {
		return fmt.Errorf("updating OIDC Client %s: %w", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"secret created"},
		Refresh:    oidcClientStateRefreshFunc(client, oidcClient.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("[ERROR] waiting for OIDCClient (%s) to be ready: %w", oidcClient.ID, err)
	}

	return resourceRancher2OIDCClientRead(d, meta)
}

func resourceRancher2OIDCClientDelete(d *schema.ResourceData, meta any) error {
	log.Printf("[INFO] Deleting OIDCClient ID %s", d.Id())

	client, err := meta.(managementClientGetter).ManagementClient()
	if err != nil {
		return fmt.Errorf("getting a client when deleting OIDC Client: %w", err)
	}

	oidcClient, err := client.OIDCClient.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting OIDC Client %s for deletion: %w", d.Id(), err)
	}

	err = client.OIDCClient.Delete(oidcClient)
	if err != nil {
		if !IsNotFound(err) {
			return fmt.Errorf("deleting OIDC Client %s: %w", d.Id(), err)
		}
	}

	d.SetId("")

	return nil
}

func resourceRancher2OIDCClientImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRancher2OIDCClientRead(d, meta)
	if err != nil || d.Id() == "" {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}

type managementClientGetter interface {
	ManagementClient() (*managementClient.Client, error)
}

type secretFetcher interface {
	SecretByName(cluster, namespace, secretName string) (*corev1.Secret, error)
}
