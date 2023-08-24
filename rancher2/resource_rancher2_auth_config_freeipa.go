package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigFreeIpa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigFreeIpaCreate,
		ReadContext:   resourceRancher2AuthConfigFreeIpaRead,
		UpdateContext: resourceRancher2AuthConfigFreeIpaUpdate,
		DeleteContext: resourceRancher2AuthConfigFreeIpaDelete,

		Schema: authConfigFreeIpaFields(),
	}
}

func resourceRancher2AuthConfigFreeIpaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigFreeIpaName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigFreeIpaName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigFreeIpaName, auth.Name)

	authFreeIpa, err := expandAuthConfigFreeIpa(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigFreeIpaName, err)
	}

	// Checking if other auth config is enabled
	if authFreeIpa.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigFreeIpaName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigFreeIpaName, err)
		}
		// Updated auth config
		authResource := &norman.Resource{
			ID:      auth.ID,
			Type:    auth.Type,
			Links:   auth.Links,
			Actions: auth.Actions,
		}
		testInput := &managementClient.FreeIpaTestAndApplyInput{
			LdapConfig: authFreeIpa,
			Password:   d.Get("test_password").(string),
			Username:   d.Get("test_username").(string),
		}
		err = client.APIBaseClient.Action(managementClient.FreeIpaConfigType, "testAndApply", authResource, testInput, nil)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		if len(auth.Actions["disable"]) > 0 {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigFreeIpaName, err)
			}
		}
		newAuth := &managementClient.FreeIpaConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authFreeIpa, newAuth)
		if err != nil {
			return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigFreeIpaName, err)
		}
	}

	return resourceRancher2AuthConfigFreeIpaRead(ctx, d, meta)
}

func resourceRancher2AuthConfigFreeIpaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigFreeIpaName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigFreeIpaName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigFreeIpaName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authFreeIpa, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigFreeIpa(d, authFreeIpa.(*managementClient.LdapConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigFreeIpaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigFreeIpaName)

	return resourceRancher2AuthConfigFreeIpaCreate(ctx, d, meta)
}

func resourceRancher2AuthConfigFreeIpaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigFreeIpaName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigFreeIpaName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigFreeIpaName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigFreeIpaName, err)
		}
	}

	d.SetId("")
	return nil
}
