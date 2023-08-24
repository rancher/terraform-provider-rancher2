package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigGithub() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigGithubCreate,
		ReadContext:   resourceRancher2AuthConfigGithubRead,
		UpdateContext: resourceRancher2AuthConfigGithubUpdate,
		DeleteContext: resourceRancher2AuthConfigGithubDelete,

		Schema: authConfigGithubFields(),
	}
}

func resourceRancher2AuthConfigGithubCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigGithubName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigGithubName, auth.Name)

	authGithub, err := expandAuthConfigGithub(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigGithubName, err)
	}

	// Checking if other auth config is enabled
	if authGithub.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigGithubName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigGithubName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.GithubConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authGithub, newAuth)
	if err != nil {
		return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigGithubName, err)
	}

	return resourceRancher2AuthConfigGithubRead(ctx, d, meta)
}

func resourceRancher2AuthConfigGithubRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigGithubName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGithubName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authGithub, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigGithub(d, authGithub.(*managementClient.GithubConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigGithubUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigGithubName)

	return resourceRancher2AuthConfigGithubCreate(ctx, d, meta)
}

func resourceRancher2AuthConfigGithubDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigGithubName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGithubName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigGithubName, err)
		}
	}

	d.SetId("")
	return nil
}
