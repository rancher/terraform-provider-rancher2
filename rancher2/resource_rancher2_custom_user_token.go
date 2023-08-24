package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

/*
The Rancher API can only create tokens for the user owning the token used to authenticate the API call. To acquire
the initial token to connect to the Rancher API a username+password login is needed. The lifespan of this token is
controlled by Rancher.

The following steps are used to create a Terraform controlled Rancher API token:
1) Acquire a temporary Rancher API token for the custom user by doing a username+password login
2) Create a Terraform controlled token by authenticating using the (temporary) token from step 1
3) Logout the temporary token acquired in step 1
*/

func resourceRancher2CustomUserToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2CustomUserTokenCreate,
		ReadContext:   resourceRancher2CustomUserTokenRead,
		UpdateContext: resourceRancher2CustomUserTokenUpdate,
		DeleteContext: resourceRancher2CustomUserTokenDelete,

		Schema: customUserTokenFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2CustomUserTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating Custom User Account Token")
	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	token, err := expandToken(d, patch)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := doUserLogin(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	defer doUserLogout(d, client)

	newToken, err := client.Token.Create(token)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenToken(d, newToken, patch)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRancher2CustomUserTokenRead(ctx, d, meta)
}

func resourceRancher2CustomUserTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Token ID %s", d.Id())
	client, err := doUserLogin(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	defer doUserLogout(d, client)

	token, err := client.Token.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Token ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	renew := d.Get("renew").(bool)
	if (!*token.Enabled || token.Expired) && renew {
		d.Set("renew", false)
	}

	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return diag.FromErr(err)
	}
	err = flattenToken(d, token, patch)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2CustomUserTokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRancher2TokenRead(ctx, d, meta)
}

func resourceRancher2CustomUserTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Token ID %s", d.Id())
	id := d.Id()
	client, err := doUserLogin(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	defer doUserLogout(d, client)

	token, err := client.Token.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Token ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.Token.Delete(token)
	if err != nil {
		return diag.Errorf("[ERROR] Deleting Token: %s", err)
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
