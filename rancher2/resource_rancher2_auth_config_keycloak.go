package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigKeyCloak() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigKeyCloakCreate,
		ReadContext:   resourceRancher2AuthConfigKeyCloakRead,
		UpdateContext: resourceRancher2AuthConfigKeyCloakUpdate,
		DeleteContext: resourceRancher2AuthConfigKeyCloakDelete,

		Schema: authConfigKeyCloakFields(),
	}
}

func resourceRancher2AuthConfigKeyCloakCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigKeyCloakName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigKeyCloakName)

	authKeyCloak, err := expandAuthConfigKeyCloak(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigKeyCloakName, err)
	}

	// Checking if other auth config is enabled
	if authKeyCloak.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigKeyCloakName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigKeyCloakName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.KeyCloakConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authKeyCloak, newAuth)
	if err != nil {
		return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigKeyCloakName, err)
	}

	return resourceRancher2AuthConfigKeyCloakRead(ctx, d, meta)
}

func resourceRancher2AuthConfigKeyCloakRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigKeyCloakName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigKeyCloakName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authKeyCloak, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigKeyCloak(d, authKeyCloak.(*managementClient.KeyCloakConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigKeyCloakUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigKeyCloakName)

	return resourceRancher2AuthConfigKeyCloakCreate(ctx, d, meta)
}

func resourceRancher2AuthConfigKeyCloakDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigKeyCloakName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigKeyCloakName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigKeyCloakName, err)
		}
	}

	d.SetId("")
	return nil
}
