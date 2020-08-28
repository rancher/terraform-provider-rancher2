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
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigGithubName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigGithubName, auth.Name)

	authGithub, err := expandAuthConfigGithub(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigGithubName, err)
	}

	// Checking if other auth config is enabled
	if authGithub.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigGithubName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigGithubName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.GithubConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authGithub, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigGithubName, err)
	}

	return resourceRancher2AuthConfigGithubRead(d, meta)
}

func resourceRancher2AuthConfigGithubRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigGithubName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGithubName)
			d.SetId("")
			return nil
		}
		return err
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
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigGithubName)

	return resourceRancher2AuthConfigGithubCreate(d, meta)
}

func resourceRancher2AuthConfigGithubDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigGithubName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGithubName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGithubName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigGithubName, err)
		}
	}

	d.SetId("")
	return nil
}
