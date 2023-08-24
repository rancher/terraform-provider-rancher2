package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigOKTA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigOKTACreate,
		ReadContext:   resourceRancher2AuthConfigOKTARead,
		UpdateContext: resourceRancher2AuthConfigOKTAUpdate,
		DeleteContext: resourceRancher2AuthConfigOKTADelete,

		Schema: authConfigOKTAFields(),
	}
}

func resourceRancher2AuthConfigOKTACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOKTAName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigOKTAName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigOKTAName)

	authOKTA, err := expandAuthConfigOKTA(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigOKTAName, err)
	}

	// Checking if other auth config is enabled
	if authOKTA.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigOKTAName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigOKTAName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.OKTAConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOKTA, newAuth)
	if err != nil {
		return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigOKTAName, err)
	}

	return resourceRancher2AuthConfigOKTARead(ctx, d, meta)
}

func resourceRancher2AuthConfigOKTARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigOKTAName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOKTAName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOKTAName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authOKTA, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigOKTA(d, authOKTA.(*managementClient.OKTAConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigOKTAUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigOKTAName)

	return resourceRancher2AuthConfigOKTACreate(ctx, d, meta)
}

func resourceRancher2AuthConfigOKTADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigOKTAName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOKTAName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOKTAName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigOKTAName, err)
		}
	}

	d.SetId("")
	return nil
}
