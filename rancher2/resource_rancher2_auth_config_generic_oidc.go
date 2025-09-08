package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigGenericOIDC() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigGenericOIDCCreate,
		Read:   resourceRancher2AuthConfigGenericOIDCRead,
		Update: resourceRancher2AuthConfigGenericOIDCUpdate,
		Delete: resourceRancher2AuthConfigGenericOIDCDelete,

		Schema: authConfigGenericOIDCFields(),
	}
}

func resourceRancher2AuthConfigGenericOIDCCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGenericOIDCName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigGenericOIDCName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigGenericOIDCName)

	authOIDC, err := expandAuthConfigGenericOIDC(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigGenericOIDCName, err)
	}

	// Checking if other auth config is enabled
	if authOIDC.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigGenericOIDCName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigGenericOIDCName, err)
		}
	}

	newAuth := &managementClient.OIDCConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOIDC, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigGenericOIDCName, err)
	}

	return resourceRancher2AuthConfigGenericOIDCRead(d, meta)
}

func resourceRancher2AuthConfigGenericOIDCRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigGenericOIDCName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGenericOIDCName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGenericOIDCName)
			d.SetId("")
			return nil
		}
		return err
	}

	authOIDC, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigGenericOIDC(d, authOIDC.(*managementClient.OIDCConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigGenericOIDCUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigGenericOIDCName)
	return resourceRancher2AuthConfigGenericOIDCCreate(d, meta)
}

func resourceRancher2AuthConfigGenericOIDCDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigGenericOIDCName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGenericOIDCName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGenericOIDCName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigGenericOIDCName, err)
		}
	}

	d.SetId("")
	return nil
}
