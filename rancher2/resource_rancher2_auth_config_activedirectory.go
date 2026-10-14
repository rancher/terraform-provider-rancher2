package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigActiveDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigActiveDirectoryCreate,
		Read:   resourceRancher2AuthConfigActiveDirectoryRead,
		Update: resourceRancher2AuthConfigActiveDirectoryUpdate,
		Delete: resourceRancher2AuthConfigActiveDirectoryDelete,

		Schema: authConfigActiveDirectoryFields(),
	}
}

func resourceRancher2AuthConfigActiveDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	authActiveDirectory, err := expandAuthConfigActiveDirectory(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	// Create the Auth Config disabled, so that enabling it goes through testAndApply.
	disabled := *authActiveDirectory
	disabled.Enabled = false
	_, err = meta.(*Config).CreateAuthConfigIfMissing(authActiveDirectory.Name, &disabled, &managementClient.ActiveDirectoryConfig{})
	if err != nil {
		return err
	}

	d.SetId(authActiveDirectory.Name)

	err = resourceRancher2AuthConfigActiveDirectoryApply(d, meta, authActiveDirectory)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigActiveDirectoryRead(d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryApply(d *schema.ResourceData, meta interface{}, authActiveDirectory *managementClient.ActiveDirectoryConfig) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", d.Id(), err)
	}

	err = checkAuthConfigType(d.Id(), auth, authActiveDirectory)
	if err != nil {
		return err
	}

	err = meta.(*Config).CheckAuthConfigCanBeEnabled(d.Id(), authActiveDirectory)
	if err != nil {
		return err
	}

	if authActiveDirectory.Enabled {
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
			return err
		}
	} else {
		if len(auth.Actions["disable"]) > 0 {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", d.Id(), err)
			}
		}
		newAuth := &managementClient.ActiveDirectoryConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authActiveDirectory, newAuth, managementClient.AuthConfigFieldDisplayName)
		if err != nil {
			return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceRancher2AuthConfigActiveDirectoryRead(d *schema.ResourceData, meta interface{}) error {
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

	authActiveDirectory, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigActiveDirectory(d, authActiveDirectory.(*managementClient.ActiveDirectoryConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigActiveDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", d.Id())

	authActiveDirectory, err := expandAuthConfigActiveDirectory(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	err = resourceRancher2AuthConfigActiveDirectoryApply(d, meta, authActiveDirectory)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigActiveDirectoryRead(d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
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
