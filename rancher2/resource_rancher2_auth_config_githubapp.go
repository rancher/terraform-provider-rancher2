package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigGithubApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigGithubAppCreate,
		Read:   resourceRancher2AuthConfigGithubAppRead,
		Update: resourceRancher2AuthConfigGithubAppUpdate,
		Delete: resourceRancher2AuthConfigGithubAppDelete,
		Schema: authConfigGithubAppFields(),
	}
}

func resourceRancher2AuthConfigGithubAppCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get ManagementClient %s", err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubAppName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigGithubAppName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigGithubAppName, auth.Name)

	authGithubApp, err := expandAuthConfigGithubApp(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigGithubAppName, err)
	}

	// Checking if other auth config is enabled
	if authGithubApp.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigGithubAppName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking if an Auth Config other than %s is enabled: %s", AuthConfigGithubAppName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.GithubAppConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authGithubApp, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigGithubAppName, err)
	}

	return resourceRancher2AuthConfigGithubAppRead(d, meta)
}

func resourceRancher2AuthConfigGithubAppRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigGithubAppName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get ManagementClient %s", err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubAppName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGithubAppName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Getting Auth Config %s By ID: %w", AuthConfigGithubAppName, err)
	}

	authGithubApp, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return fmt.Errorf("[ERROR] Getting Auth Config %s: %w", AuthConfigGithubAppName, err)
	}

	err = flattenAuthConfigGithubApp(d, authGithubApp.(*managementClient.GithubAppConfig))
	if err != nil {
		return fmt.Errorf("[ERROR] Flattening GitHub app: %w", err)
	}

	return nil
}

func resourceRancher2AuthConfigGithubAppUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigGithubAppName)

	return resourceRancher2AuthConfigGithubAppCreate(d, meta)
}

func resourceRancher2AuthConfigGithubAppDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigGithubAppName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get ManagementClient %s", err)
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubAppName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGithubAppName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Getting the %s AuthConfig: %w", AuthConfigGithubAppName, err)
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigGithubAppName, err)
		}
	}

	d.SetId("")
	return nil
}
