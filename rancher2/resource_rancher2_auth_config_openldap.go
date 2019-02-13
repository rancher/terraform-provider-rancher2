package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const OpenLdapConfigName = "openldap"

//Schemas

func authConfigOpenLdapFields() map[string]*schema.Schema {
	return authConfigLdapFields()
}

// Flatteners

func flattenAuthConfigOpenLdap(d *schema.ResourceData, in *managementClient.LdapConfig) error {
	d.SetId(OpenLdapConfigName)

	err := d.Set("name", OpenLdapConfigName)
	if err != nil {
		return err
	}
	err = d.Set("type", managementClient.OpenLdapConfigType)
	if err != nil {
		return err
	}

	err = flattenAuthConfigLdap(d, in)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAuthConfigOpenLdap(in *schema.ResourceData) (*managementClient.LdapConfig, error) {
	obj, err := expandAuthConfigLdap(in)
	if err != nil {
		return nil, err
	}

	obj.Name = OpenLdapConfigName
	obj.Type = managementClient.OpenLdapConfigType

	return obj, nil
}

func resourceRancher2AuthConfigOpenLdap() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigOpenLdapCreate,
		Read:   resourceRancher2AuthConfigOpenLdapRead,
		Update: resourceRancher2AuthConfigOpenLdapUpdate,
		Delete: resourceRancher2AuthConfigOpenLdapDelete,

		Schema: authConfigOpenLdapFields(),
	}
}

func resourceRancher2AuthConfigOpenLdapCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(OpenLdapConfigName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", OpenLdapConfigName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", OpenLdapConfigName)

	authOpenLdap, err := expandAuthConfigOpenLdap(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", OpenLdapConfigName, err)
	}

	authOpenLdapTestAndApply := managementClient.OpenLdapTestAndApplyInput{
		LdapConfig: authOpenLdap,
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
	}

	err = client.Post(auth.Actions["testAndApply"], authOpenLdapTestAndApply, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Posting Auth Config %s: %s", OpenLdapConfigName, err)
	}

	return resourceRancher2AuthConfigOpenLdapRead(d, meta)
}

func resourceRancher2AuthConfigOpenLdapRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", OpenLdapConfigName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(OpenLdapConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", OpenLdapConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	authOpenLdap, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigOpenLdap(d, authOpenLdap.(*managementClient.LdapConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigOpenLdapUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", OpenLdapConfigName)

	return resourceRancher2AuthConfigOpenLdapCreate(d, meta)
}

func resourceRancher2AuthConfigOpenLdapDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", OpenLdapConfigName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(OpenLdapConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", OpenLdapConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", OpenLdapConfigName, err)
		}
	}

	d.SetId("")
	return nil
}
