package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ServiceAccountToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ServiceAccountTokenCreate,
		Read:   resourceRancher2ServiceAccountTokenRead,
		Update: resourceRancher2ServiceAccountTokenUpdate,
		Delete: resourceRancher2ServiceAccountTokenDelete,

		Schema: serviceAccountTokenFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2ServiceAccountTokenCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Service Account Token")
	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return err
	}

	token, err := expandToken(d, patch)
	if err != nil {
		return err
	}

	client, err := doUserLogin(d, meta)
	if err != nil {
		return err
	}
	defer doUserLogout(d, client)

	newToken, err := client.Token.Create(token)
	if err != nil {
		return err
	}

	err = flattenToken(d, newToken, patch)
	if err != nil {
		return err
	}

	return resourceRancher2ServiceAccountTokenRead(d, meta)
}

func resourceRancher2ServiceAccountTokenRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Token ID %s", d.Id())
	client, err := doUserLogin(d, meta)
	if err != nil {
		return err
	}
	defer doUserLogout(d, client)

	token, err := client.Token.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Token ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	renew := d.Get("renew").(bool)
	if (!*token.Enabled || token.Expired) && renew {
		d.Set("renew", false)
	}

	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return err
	}
	err = flattenToken(d, token, patch)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2ServiceAccountTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRancher2TokenRead(d, meta)
}

func resourceRancher2ServiceAccountTokenDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Token ID %s", d.Id())
	id := d.Id()
	client, err := doUserLogin(d, meta)
	if err != nil {
		return err
	}
	defer doUserLogout(d, client)

	token, err := client.Token.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Token ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Token.Delete(token)
	if err != nil {
		return fmt.Errorf("[ERROR] Deleting Token: %s", err)
	}

	d.SetId("")
	return nil
}

func doUserLogin(d *schema.ResourceData, meta interface{}) (*managementClient.Client, error) {
	if client, err := getManagementClientForTempToken(d, meta); err != nil || client != nil {
		return client, err
	}

	log.Printf("[DEBUG] Creating Temp API Token for User %s", d.Get("username").(string))
	tempTokenID, tempTokenValue, err := DoUserLogin(meta.(*Config).URL, d.Get("username").(string), d.Get("password").(string), "0", "Temp Terraform API token", meta.(*Config).CACerts, meta.(*Config).Insecure)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Login with %s user: %v", d.Get("username").(string), err)
	}

	d.Set("temp_token_id", tempTokenID)
	d.Set("temp_token", tempTokenValue)

	return newManagementClient(meta, tempTokenValue)
}

func getManagementClientForTempToken(d *schema.ResourceData, meta interface{}) (*managementClient.Client, error) {
	if v, ok := d.Get("temp_token").(string); ok && len(v) > 0 {
		client, err := newManagementClient(meta, v)
		if err != nil {
			d.Set("temp_token_id", "")
			d.Set("temp_token", "")
			if !IsNotFound(err) && !IsUnauthorized(err) && !IsForbidden(err) {
				return nil, err
			}
			return nil, nil
		}

		return client, nil
	}

	return nil, nil
}

func newManagementClient(meta interface{}, token string) (*managementClient.Client, error) {
	options := meta.(*Config).CreateClientOpts()
	options.URL = options.URL + rancher2ClientAPIVersion
	options.TokenKey = token
	return managementClient.NewClient(options)
}

func doUserLogout(d *schema.ResourceData, client *managementClient.Client) error {
	if v, ok := d.Get("temp_token_id").(string); ok && len(v) > 0 {
		log.Printf("[DEBUG] Deleting Temp API Token for User %s", d.Id())
		existing := &norman.Resource{
			ID: v,
			Actions: map[string]string{
				"logout": client.Opts.URL + "/tokens?action=logout",
			},
		}

		err := client.APIBaseClient.Action("token", "logout", existing, map[string]interface{}{}, nil)
		if err != nil {
			return err
		}

		d.Set("temp_token_id", "")
		d.Set("temp_token", "")
	}

	return nil
}
