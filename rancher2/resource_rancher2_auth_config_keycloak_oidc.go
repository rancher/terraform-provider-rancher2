package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigKeyCloakOIDC() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigKeyCloakOIDCCreate,
		Read:   resourceRancher2AuthConfigKeyCloakOIDCRead,
		Update: resourceRancher2AuthConfigKeyCloakOIDCUpdate,
		Delete: resourceRancher2AuthConfigKeyCloakOIDCDelete,

		Schema: authConfigKeyCloakOIDCFields(),
	}
}

func resourceRancher2AuthConfigKeyCloakOIDCCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakOIDCName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigKeyCloakOIDCName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigKeyCloakOIDCName)

	authOIDC, err := expandAuthConfigKeyCloakOIDC(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigKeyCloakOIDCName, err)
	}

	// Checking if other auth config is enabled
	if authOIDC.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigKeyCloakOIDCName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigKeyCloakOIDCName, err)
		}
	}

	newAuth := &managementClient.KeyCloakOIDCConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOIDC, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigKeyCloakOIDCName, err)
	}

	return resourceRancher2AuthConfigKeyCloakOIDCRead(d, meta)
}

func resourceRancher2AuthConfigKeyCloakOIDCRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigKeyCloakOIDCName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakOIDCName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigKeyCloakOIDCName)
			d.SetId("")
			return nil
		}
		return err
	}

	authOIDC, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigKeyCloakOIDC(d, authOIDC.(*managementClient.KeyCloakOIDCConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigKeyCloakOIDCUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigKeyCloakOIDCName)
	return resourceRancher2AuthConfigKeyCloakOIDCCreate(d, meta)
}

func resourceRancher2AuthConfigKeyCloakOIDCDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigKeyCloakOIDCName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigKeyCloakOIDCName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigKeyCloakOIDCName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigKeyCloakOIDCName, err)
		}
	}

	d.SetId("")
	return nil
}
