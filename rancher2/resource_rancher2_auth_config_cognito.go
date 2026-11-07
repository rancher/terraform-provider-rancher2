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
	authCognito, err := expandAuthConfigCognito(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.CognitoConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authCognito.Name, authCognito, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authCognito.Name)

	return resourceRancher2AuthConfigCognitoRead(d, meta)
}

func resourceRancher2AuthConfigCognitoRead(d *schema.ResourceData, meta interface{}) error {
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

	authCognito, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigCognito(d, authCognito.(*managementClient.CognitoConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigCognitoUpdate(d *schema.ResourceData, meta interface{}) error {
	authCognito, err := expandAuthConfigCognito(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.CognitoConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authCognito, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigCognitoRead(d, meta)
}

func resourceRancher2AuthConfigCognitoDelete(d *schema.ResourceData, meta interface{}) error {
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
