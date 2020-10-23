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
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigActiveDirectoryName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigActiveDirectoryName, auth.Name)

	authActiveDirectory, err := expandAuthConfigActiveDirectory(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
	}

	log.Printf("[INFO] +++ %v", authActiveDirectory)

	// Checking if other auth config is enabled
	if authActiveDirectory.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigActiveDirectoryName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
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
			return err
		}
	} else {
		if len(auth.Actions["disable"]) > 0 {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
			}
		}
		newAuth := &managementClient.ActiveDirectoryConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authActiveDirectory, newAuth)
		if err != nil {
			return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
		}
	}

	return resourceRancher2AuthConfigActiveDirectoryRead(d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigActiveDirectoryName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigActiveDirectoryName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigActiveDirectoryName)
			d.SetId("")
			return nil
		}
		return err
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
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigActiveDirectoryName)

	return resourceRancher2AuthConfigActiveDirectoryCreate(d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigActiveDirectoryName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigActiveDirectoryName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigActiveDirectoryName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigActiveDirectoryName, err)
		}
	}

	d.SetId("")
	return nil
}
