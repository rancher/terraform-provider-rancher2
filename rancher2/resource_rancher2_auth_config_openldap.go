package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigOpenLdap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigOpenLdapCreate,
		ReadContext:   resourceRancher2AuthConfigOpenLdapRead,
		UpdateContext: resourceRancher2AuthConfigOpenLdapUpdate,
		DeleteContext: resourceRancher2AuthConfigOpenLdapDelete,

		Schema: authConfigOpenLdapFields(),
	}
}

func resourceRancher2AuthConfigOpenLdapCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOpenLdapName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigOpenLdapName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigOpenLdapName)

	authOpenLdap, err := expandAuthConfigOpenLdap(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigOpenLdapName, err)
	}

	// Checking if other auth config is enabled
	if authOpenLdap.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigOpenLdapName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigOpenLdapName, err)
		}
		// Updated auth config
		authResource := &norman.Resource{
			ID:      auth.ID,
			Type:    auth.Type,
			Links:   auth.Links,
			Actions: auth.Actions,
		}
		testInput := &managementClient.OpenLdapTestAndApplyInput{
			LdapConfig: authOpenLdap,
			Password:   d.Get("test_password").(string),
			Username:   d.Get("test_username").(string),
		}
		err = client.APIBaseClient.Action(managementClient.OpenLdapConfigType, "testAndApply", authResource, testInput, nil)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		if len(auth.Actions["disable"]) > 0 {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigOpenLdapName, err)
			}
		}
		newAuth := &managementClient.OpenLdapConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOpenLdap, newAuth)
		if err != nil {
			return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigOpenLdapName, err)
		}
	}

	return resourceRancher2AuthConfigOpenLdapRead(ctx, d, meta)
}

func resourceRancher2AuthConfigOpenLdapRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigOpenLdapName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOpenLdapName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOpenLdapName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authOpenLdap, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigOpenLdap(d, authOpenLdap.(*managementClient.LdapConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigOpenLdapUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigOpenLdapName)

	return resourceRancher2AuthConfigOpenLdapCreate(ctx, d, meta)
}

func resourceRancher2AuthConfigOpenLdapDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigOpenLdapName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOpenLdapName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOpenLdapName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigOpenLdapName, err)
		}
	}

	d.SetId("")
	return nil
}
