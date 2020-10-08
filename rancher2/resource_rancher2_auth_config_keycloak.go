package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigKeyCloak() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigKeyCloakCreate,
		Read:   resourceRancher2AuthConfigKeyCloakRead,
		Update: resourceRancher2AuthConfigKeyCloakUpdate,
		Delete: resourceRancher2AuthConfigKeyCloakDelete,

		Schema: authConfigKeyCloakFields(),
	}
}

func resourceRancher2AuthConfigKeyCloakCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigKeyCloakName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigKeyCloakName)

	authKeyCloak, err := expandAuthConfigKeyCloak(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigKeyCloakName, err)
	}

	// Checking if other auth config is enabled
	if authKeyCloak.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigKeyCloakName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigKeyCloakName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.KeyCloakConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authKeyCloak, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigKeyCloakName, err)
	}

	return resourceRancher2AuthConfigKeyCloakRead(d, meta)
}

func resourceRancher2AuthConfigKeyCloakRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigKeyCloakName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigKeyCloakName)
			d.SetId("")
			return nil
		}
		return err
	}

	authKeyCloak, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigKeyCloak(d, authKeyCloak.(*managementClient.KeyCloakConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigKeyCloakUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigKeyCloakName)

	return resourceRancher2AuthConfigKeyCloakCreate(d, meta)
}

func resourceRancher2AuthConfigKeyCloakDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigKeyCloakName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigKeyCloakName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigKeyCloakName, err)
		}
	}

	d.SetId("")
	return nil
}
