package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigOKTA() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigOKTACreate,
		Read:   resourceRancher2AuthConfigOKTARead,
		Update: resourceRancher2AuthConfigOKTAUpdate,
		Delete: resourceRancher2AuthConfigOKTADelete,

		Schema: authConfigOKTAFields(),
	}
}

func resourceRancher2AuthConfigOKTACreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOKTAName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigOKTAName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigOKTAName)

	authOKTA, err := expandAuthConfigOKTA(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigOKTAName, err)
	}

	// Checking if other auth config is enabled
	if authOKTA.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigOKTAName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigOKTAName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.OKTAConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOKTA, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigOKTAName, err)
	}

	return resourceRancher2AuthConfigOKTARead(d, meta)
}

func resourceRancher2AuthConfigOKTARead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigOKTAName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOKTAName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOKTAName)
			d.SetId("")
			return nil
		}
		return err
	}

	authOKTA, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigOKTA(d, authOKTA.(*managementClient.OKTAConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigOKTAUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigOKTAName)

	return resourceRancher2AuthConfigOKTACreate(d, meta)
}

func resourceRancher2AuthConfigOKTADelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigOKTAName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOKTAName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOKTAName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigOKTAName, err)
		}
	}

	d.SetId("")
	return nil
}
