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
	authKeyCloak, err := expandAuthConfigKeyCloak(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.KeyCloakConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authKeyCloak.Name, authKeyCloak, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authKeyCloak.Name)

	return resourceRancher2AuthConfigKeyCloakRead(d, meta)
}

func resourceRancher2AuthConfigKeyCloakRead(d *schema.ResourceData, meta interface{}) error {
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
	authKeyCloak, err := expandAuthConfigKeyCloak(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.KeyCloakConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authKeyCloak, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigKeyCloakRead(d, meta)
}

func resourceRancher2AuthConfigKeyCloakDelete(d *schema.ResourceData, meta interface{}) error {
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
