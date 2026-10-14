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
	authGithubApp, err := expandAuthConfigGithubApp(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.GithubAppConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authGithubApp.Name, authGithubApp, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authGithubApp.Name)

	return resourceRancher2AuthConfigGithubAppRead(d, meta)
}

func resourceRancher2AuthConfigGithubAppRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Getting Auth Config %s: %w", d.Id(), err)
	}

	authGithubApp, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigGithubApp(d, authGithubApp.(*managementClient.GithubAppConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigGithubAppUpdate(d *schema.ResourceData, meta interface{}) error {
	authGithubApp, err := expandAuthConfigGithubApp(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.GithubAppConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authGithubApp, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigGithubAppRead(d, meta)
}

func resourceRancher2AuthConfigGithubAppDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Getting Auth Config %s: %w", d.Id(), err)
	}

	if auth.Enabled {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", d.Id(), err)
		}
	}

	d.SetId("")
	return nil
}
