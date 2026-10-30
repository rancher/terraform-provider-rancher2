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
	authKeyCloakOIDC, err := expandAuthConfigKeyCloakOIDC(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.KeyCloakOIDCConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authKeyCloakOIDC.Name, authKeyCloakOIDC, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authKeyCloakOIDC.Name)

	return resourceRancher2AuthConfigKeyCloakOIDCRead(d, meta)
}

func resourceRancher2AuthConfigKeyCloakOIDCRead(d *schema.ResourceData, meta interface{}) error {
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

	authKeyCloakOIDC, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigKeyCloakOIDC(d, authKeyCloakOIDC.(*managementClient.KeyCloakOIDCConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigKeyCloakOIDCUpdate(d *schema.ResourceData, meta interface{}) error {
	authKeyCloakOIDC, err := expandAuthConfigKeyCloakOIDC(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.KeyCloakOIDCConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authKeyCloakOIDC, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigKeyCloakOIDCRead(d, meta)
}

func resourceRancher2AuthConfigKeyCloakOIDCDelete(d *schema.ResourceData, meta interface{}) error {
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
