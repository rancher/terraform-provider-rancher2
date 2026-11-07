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
	authGenericOIDC, err := expandAuthConfigGenericOIDC(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.GenericOIDCConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authGenericOIDC.Name, authGenericOIDC, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authGenericOIDC.Name)

	return resourceRancher2AuthConfigGenericOIDCRead(d, meta)
}

func resourceRancher2AuthConfigGenericOIDCRead(d *schema.ResourceData, meta interface{}) error {
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

	authGenericOIDC, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigGenericOIDC(d, authGenericOIDC.(*managementClient.GenericOIDCConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigGenericOIDCUpdate(d *schema.ResourceData, meta interface{}) error {
	authGenericOIDC, err := expandAuthConfigGenericOIDC(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.GenericOIDCConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authGenericOIDC, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigGenericOIDCRead(d, meta)
}

func resourceRancher2AuthConfigGenericOIDCDelete(d *schema.ResourceData, meta interface{}) error {
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
