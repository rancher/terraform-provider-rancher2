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
	authPing, err := expandAuthConfigPing(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.PingConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authPing.Name, authPing, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authPing.Name)

	return resourceRancher2AuthConfigPingRead(d, meta)
}

func resourceRancher2AuthConfigPingRead(d *schema.ResourceData, meta interface{}) error {
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
	authPing, err := expandAuthConfigPing(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.PingConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authPing, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigPingRead(d, meta)
}

func resourceRancher2AuthConfigPingDelete(d *schema.ResourceData, meta interface{}) error {
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
