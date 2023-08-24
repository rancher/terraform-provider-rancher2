package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigPing() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigPingCreate,
		ReadContext:   resourceRancher2AuthConfigPingRead,
		UpdateContext: resourceRancher2AuthConfigPingUpdate,
		DeleteContext: resourceRancher2AuthConfigPingDelete,

		Schema: authConfigPingFields(),
	}
}

func resourceRancher2AuthConfigPingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigPingName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigPingName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigPingName)

	authPing, err := expandAuthConfigPing(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigPingName, err)
	}

	// Checking if other auth config is enabled
	if authPing.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigPingName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigPingName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.PingConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authPing, newAuth)
	if err != nil {
		return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigPingName, err)
	}

	return resourceRancher2AuthConfigPingRead(ctx, d, meta)
}

func resourceRancher2AuthConfigPingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigPingName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigPingName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigPingName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authPing, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigPing(d, authPing.(*managementClient.PingConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigPingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigPingName)

	return resourceRancher2AuthConfigPingCreate(ctx, d, meta)
}

func resourceRancher2AuthConfigPingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigPingName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigPingName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigPingName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigPingName, err)
		}
	}

	d.SetId("")
	return nil
}
