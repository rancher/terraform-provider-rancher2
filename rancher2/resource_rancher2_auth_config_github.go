package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigGithub() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigGithubCreate,
		Read:   resourceRancher2AuthConfigGithubRead,
		Update: resourceRancher2AuthConfigGithubUpdate,
		Delete: resourceRancher2AuthConfigGithubDelete,

		Schema: authConfigGithubFields(),
	}
}

func resourceRancher2AuthConfigGithubCreate(d *schema.ResourceData, meta interface{}) error {
	authGithub, err := expandAuthConfigGithub(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.GithubConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authGithub.Name, authGithub, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authGithub.Name)

	return resourceRancher2AuthConfigGithubRead(d, meta)
}

func resourceRancher2AuthConfigGithubRead(d *schema.ResourceData, meta interface{}) error {
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

	authGithub, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigGithub(d, authGithub.(*managementClient.GithubConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigGithubUpdate(d *schema.ResourceData, meta interface{}) error {
	authGithub, err := expandAuthConfigGithub(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.GithubConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authGithub, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigGithubRead(d, meta)
}

func resourceRancher2AuthConfigGithubDelete(d *schema.ResourceData, meta interface{}) error {
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
