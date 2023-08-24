package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigADFS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigADFSCreate,
		ReadContext:   resourceRancher2AuthConfigADFSRead,
		UpdateContext: resourceRancher2AuthConfigADFSUpdate,
		DeleteContext: resourceRancher2AuthConfigADFSDelete,

		Schema: authConfigADFSFields(),
	}
}

func resourceRancher2AuthConfigADFSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigADFSName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigADFSName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigADFSName, auth.Name)

	authADFS, err := expandAuthConfigADFS(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigADFSName, err)
	}

	// Checking if other auth config is enabled
	if authADFS.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigADFSName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigADFSName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.ADFSConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authADFS, newAuth)
	if err != nil {
		return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigADFSName, err)
	}

	return resourceRancher2AuthConfigADFSRead(ctx, d, meta)
}

func resourceRancher2AuthConfigADFSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigADFSName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigADFSName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigADFSName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authADFS, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigADFS(d, authADFS.(*managementClient.ADFSConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigADFSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigADFSName)

	return resourceRancher2AuthConfigADFSCreate(ctx, d, meta)
}

func resourceRancher2AuthConfigADFSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigADFSName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigADFSName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigADFSName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigADFSName, err)
		}
	}

	d.SetId("")
	return nil
}
