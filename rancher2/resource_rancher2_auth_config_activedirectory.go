package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigActiveDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AuthConfigActiveDirectoryCreate,
		ReadContext:   resourceRancher2AuthConfigActiveDirectoryRead,
		UpdateContext: resourceRancher2AuthConfigActiveDirectoryUpdate,
		DeleteContext: resourceRancher2AuthConfigActiveDirectoryDelete,

		Schema: authConfigActiveDirectoryFields(),
	}
}

func resourceRancher2AuthConfigActiveDirectoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigActiveDirectoryName)
	if err != nil {
		return diag.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigActiveDirectoryName, auth.Name)

	authActiveDirectory, err := expandAuthConfigActiveDirectory(d)
	if err != nil {
		return diag.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
	}

	log.Printf("[INFO] +++ %v", authActiveDirectory)

	// Checking if other auth config is enabled
	if authActiveDirectory.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigActiveDirectoryName)
		if err != nil {
			return diag.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
		}
		// Updated auth config
		authResource := &norman.Resource{
			ID:      auth.ID,
			Type:    auth.Type,
			Links:   auth.Links,
			Actions: auth.Actions,
		}
		testInput := &managementClient.ActiveDirectoryTestAndApplyInput{
			ActiveDirectoryConfig: authActiveDirectory,
			Enabled:               authActiveDirectory.Enabled,
			Password:              d.Get("test_password").(string),
			Username:              d.Get("test_username").(string),
		}
		err = client.APIBaseClient.Action(managementClient.ActiveDirectoryConfigType, "testAndApply", authResource, testInput, nil)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		if len(auth.Actions["disable"]) > 0 {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
			}
		}
		newAuth := &managementClient.ActiveDirectoryConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authActiveDirectory, newAuth)
		if err != nil {
			return diag.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
		}
	}

	return resourceRancher2AuthConfigActiveDirectoryRead(ctx, d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigActiveDirectoryName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigActiveDirectoryName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigActiveDirectoryName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	authActiveDirectory, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenAuthConfigActiveDirectory(d, authActiveDirectory.(*managementClient.ActiveDirectoryConfig))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2AuthConfigActiveDirectoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigActiveDirectoryName)

	return resourceRancher2AuthConfigActiveDirectoryCreate(ctx, d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigActiveDirectoryName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigActiveDirectoryName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigActiveDirectoryName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return diag.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
		}
	}

	d.SetId("")
	return nil
}
