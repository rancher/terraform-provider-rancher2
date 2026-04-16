package rancher2

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2OIDCClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2OIDCClientCreate,
		Read:   resourceRancher2OIDCClientRead,
		Update: resourceRancher2OIDCClientUpdate,
		Delete: resourceRancher2OIDCClientDelete,

		Schema: oidcClientFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func oidcClientFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"token_expiration_seconds": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			Description:  "The duration (in seconds) before an access token and ID token expires. Reducing this will invalidate existing tokens.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"refresh_token_expiration_seconds": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			Description:  "The duration (in seconds) a refresh token remains valid before expiration. Reducing this will invalidate existing tokens.",
			ValidateFunc: validation.IntAtLeast(1),
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

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		log.Printf("[ERROR] getting a client when creating OIDC Client: %s", err)
		return err
	}

	updatedClient, err := client.OIDCClient.Create(oidcClient)
	if err != nil {
		log.Printf("[ERROR] creating OIDCClient: %s", err)
		return err
	}

	d.SetId(updatedClient.ID)

	return resourceRancher2OIDCClientRead(d, meta)
}

func resourceRancher2OIDCClientRead(d *schema.ResourceData, meta any) error {
	log.Printf("[INFO] Refreshing OIDCClient ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		log.Printf("[ERROR] getting a client when reading OIDC Client: %s", err)
		return err
	}
	oidcClient, err := client.OIDCClient.ByID(d.Id())
	if err != nil {
		log.Printf("[ERROR] reading OIDC Client %s: %s", d.Id(), err)
		return err
	}

	return flattenOIDCClient(d, oidcClient)
}

func resourceRancher2OIDCClientUpdate(d *schema.ResourceData, meta any) error {
	log.Printf("[INFO] Updating OIDCClient ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		log.Printf("[ERROR] getting a client when updating OIDC Client: %s", err)
		return err
	}

	oidcClient, err := client.OIDCClient.ByID(d.Id())
	if err != nil {
		log.Printf("[ERROR] getting OIDC Client %s for update: %s", d.Id(), err)
		return err
	}

	update := map[string]any{
		"description":                      d.Get("description"),
		"redirectURIs":                     d.Get("redirect_uris"),
		"token_expiration_seconds":         d.Get("token_expiration_seconds"),
		"refresh_token_expiration_seconds": d.Get("refresh_token_expiration_seconds"),
	}

	_, err = client.OIDCClient.Update(oidcClient, update)
	if err != nil {
		log.Printf("[ERROR] updating OIDC Client %s: %s", d.Id(), err)
		return fmt.Errorf("updating OIDC Client %s: %w", d.Id(), err)
	}

	return resourceRancher2OIDCClientRead(d, meta)
}

func resourceRancher2OIDCClientDelete(d *schema.ResourceData, meta any) error {
	log.Printf("[INFO] Deleting OIDCClient ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		log.Printf("[ERROR] getting a client when deleting OIDC Client: %s", err)
		return err
	}

	oidcClient, err := client.OIDCClient.ByID(d.Id())
	if err != nil {
		log.Printf("[ERROR] getting OIDC Client %s for deletion: %s", d.Id(), err)
		return nil
	}

	if err := client.OIDCClient.Delete(oidcClient); err != nil {
		log.Printf("[ERROR] deleting OIDC Client %s: %s", d.Id(), err)
		return err
	}

	return nil
}

func expandOIDCClient(in *schema.ResourceData) (*managementClient.OIDCClient, error) {
	obj := &managementClient.OIDCClient{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding OidcClient: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.GetOk("description"); ok {
		obj.Description = v.(string)
	}

	v := in.Get("redirect_uris")
	obj.RedirectURIs = toArrayString(v.([]any))

	v, ok := in.GetOk("token_expiration_seconds")
	if ok {
		// This should be safe because the field declared as an integer.
		tokenExpirationSeconds, _ := v.(int)
		obj.TokenExpirationSeconds = int64(tokenExpirationSeconds)
	} else {
		obj.TokenExpirationSeconds = 0
	}

	v, ok = in.GetOk("refresh_token_expiration_seconds")
	if ok {
		// This should be safe because the field declared as an integer.
		refreshTokenExpirationSeconds, _ := v.(int)
		obj.RefreshTokenExpirationSeconds = int64(refreshTokenExpirationSeconds)
	} else {
		obj.RefreshTokenExpirationSeconds = 0
	}

	if v, ok := in.Get("annotations").(map[string]any); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]any); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func flattenOIDCClient(d *schema.ResourceData, in *managementClient.OIDCClient) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening OIDCClient: Input config is nil")
	}

	if in.ID != "" {
		d.SetId(in.ID)
	}

	return errors.Join(
		d.Set("description", in.Description),
		d.Set("redirect_uris", in.RedirectURIs),
		d.Set("token_expiration_seconds", int(in.TokenExpirationSeconds)),
		d.Set("refresh_token_expiration_seconds", int(in.RefreshTokenExpirationSeconds)),
		d.Set("annotations", toMapInterface(in.Annotations)),
		d.Set("labels", toMapInterface(in.Labels)),
	)
}
