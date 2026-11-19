package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigAzureAD() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigAzureADCreate,
		Read:   resourceRancher2AuthConfigAzureADRead,
		Update: resourceRancher2AuthConfigAzureADUpdate,
		Delete: resourceRancher2AuthConfigAzureADDelete,

		Schema: authConfigAzureADFields(),
	}
}

func resourceRancher2AuthConfigAzureADCreate(d *schema.ResourceData, meta interface{}) error {
	authAzureAD, err := expandAuthConfigAzureAD(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.AzureADConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authAzureAD.Name, authAzureAD, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authAzureAD.Name)

	return resourceRancher2AuthConfigAzureADRead(d, meta)
}

func resourceRancher2AuthConfigAzureADRead(d *schema.ResourceData, meta interface{}) error {
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

	authAzureAD, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigAzureAD(d, authAzureAD.(*managementClient.AzureADConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigAzureADUpdate(d *schema.ResourceData, meta interface{}) error {
	authAzureAD, err := expandAuthConfigAzureAD(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.AzureADConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authAzureAD, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigAzureADRead(d, meta)
}

func resourceRancher2AuthConfigAzureADDelete(d *schema.ResourceData, meta interface{}) error {
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
