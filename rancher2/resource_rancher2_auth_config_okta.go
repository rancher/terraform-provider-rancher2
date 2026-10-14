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
	authOKTA, err := expandAuthConfigOKTA(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.OKTAConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authOKTA.Name, authOKTA, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authOKTA.Name)

	return resourceRancher2AuthConfigOKTARead(d, meta)
}

func resourceRancher2AuthConfigOKTARead(d *schema.ResourceData, meta interface{}) error {
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
	authOKTA, err := expandAuthConfigOKTA(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.OKTAConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authOKTA, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigOKTARead(d, meta)
}

func resourceRancher2AuthConfigOKTADelete(d *schema.ResourceData, meta interface{}) error {
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
