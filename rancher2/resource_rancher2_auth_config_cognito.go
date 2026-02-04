package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigCognito() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigCognitoCreate,
		Read:   resourceRancher2AuthConfigCognitoRead,
		Update: resourceRancher2AuthConfigCognitoUpdate,
		Delete: resourceRancher2AuthConfigCognitoDelete,

		Schema: authConfigCognitoFields(),
	}
}

func resourceRancher2AuthConfigCognitoCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigCognitoName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigCognitoName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigCognitoName)

	authOIDC, err := expandAuthConfigCognito(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigCognitoName, err)
	}

	// Checking if other auth config is enabled
	if authOIDC.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigCognitoName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigCognitoName, err)
		}
	}

	newAuth := &managementClient.CognitoConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOIDC, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigCognitoName, err)
	}

	return resourceRancher2AuthConfigCognitoRead(d, meta)
}

func resourceRancher2AuthConfigCognitoRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigCognitoName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigCognitoName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigCognitoName)
			d.SetId("")
			return nil
		}
		return err
	}

	authOIDC, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigCognito(d, authOIDC.(*managementClient.CognitoConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigCognitoUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigCognitoName)
	return resourceRancher2AuthConfigCognitoCreate(d, meta)
}

func resourceRancher2AuthConfigCognitoDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigCognitoName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigCognitoName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigCognitoName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigCognitoName, err)
		}
	}

	d.SetId("")
	return nil
}
