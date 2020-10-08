package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigPing() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigPingCreate,
		Read:   resourceRancher2AuthConfigPingRead,
		Update: resourceRancher2AuthConfigPingUpdate,
		Delete: resourceRancher2AuthConfigPingDelete,

		Schema: authConfigPingFields(),
	}
}

func resourceRancher2AuthConfigPingCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigPingName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigPingName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigPingName)

	authPing, err := expandAuthConfigPing(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigPingName, err)
	}

	// Checking if other auth config is enabled
	if authPing.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigPingName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigPingName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.PingConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authPing, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigPingName, err)
	}

	return resourceRancher2AuthConfigPingRead(d, meta)
}

func resourceRancher2AuthConfigPingRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigPingName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigPingName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigPingName)
			d.SetId("")
			return nil
		}
		return err
	}

	authPing, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigPing(d, authPing.(*managementClient.PingConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigPingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigPingName)

	return resourceRancher2AuthConfigPingCreate(d, meta)
}

func resourceRancher2AuthConfigPingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigPingName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigPingName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigPingName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigPingName, err)
		}
	}

	d.SetId("")
	return nil
}
