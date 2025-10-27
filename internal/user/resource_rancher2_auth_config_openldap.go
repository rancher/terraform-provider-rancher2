package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

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

	auth, err := client.AuthConfig.ByID(AuthConfigOpenLdapName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigOpenLdapName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", AuthConfigOpenLdapName)

	authOpenLdap, err := expandAuthConfigOpenLdap(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigOpenLdapName, err)
	}

	// Checking if other auth config is enabled
	if authOpenLdap.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigOpenLdapName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigOpenLdapName, err)
		}
		// Updated auth config
		authResource := &norman.Resource{
			ID:      auth.ID,
			Type:    auth.Type,
			Links:   auth.Links,
			Actions: auth.Actions,
		}
		testInput := &managementClient.OpenLdapTestAndApplyInput{
			LdapConfig: authOpenLdap,
			Password:   d.Get("test_password").(string),
			Username:   d.Get("test_username").(string),
		}
		err = client.APIBaseClient.Action(managementClient.OpenLdapConfigType, "testAndApply", authResource, testInput, nil)
		if err != nil {
			return err
		}
	} else {
		if len(auth.Actions["disable"]) > 0 {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigOpenLdapName, err)
			}
		}
		newAuth := &managementClient.OpenLdapConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOpenLdap, newAuth)
		if err != nil {
			return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigOpenLdapName, err)
		}
	}

	return resourceRancher2AuthConfigOpenLdapRead(d, meta)
}

func resourceRancher2AuthConfigOpenLdapRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigOpenLdapName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOpenLdapName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOpenLdapName)
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
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigOpenLdapName)

	return resourceRancher2AuthConfigOpenLdapCreate(d, meta)
}

func resourceRancher2AuthConfigOpenLdapDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigOpenLdapName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigOpenLdapName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigOpenLdapName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigOpenLdapName, err)
		}
	}

	d.SetId("")
	return nil
}
