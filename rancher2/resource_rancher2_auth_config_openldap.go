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
	authOpenLdap, err := expandAuthConfigOpenLdap(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	// Create the Auth Config disabled, so that enabling it goes through testAndApply.
	disabled := *authOpenLdap
	disabled.Enabled = false
	_, err = meta.(*Config).CreateAuthConfigIfMissing(authOpenLdap.Name, &disabled, &managementClient.OpenLdapConfig{})
	if err != nil {
		return err
	}

	d.SetId(authOpenLdap.Name)

	err = resourceRancher2AuthConfigOpenLdapApply(d, meta, authOpenLdap)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigOpenLdapRead(d, meta)
}

func resourceRancher2AuthConfigOpenLdapApply(d *schema.ResourceData, meta interface{}, authOpenLdap *managementClient.LdapConfig) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", d.Id(), err)
	}

	err = checkAuthConfigType(d.Id(), auth, authOpenLdap)
	if err != nil {
		return err
	}

	err = meta.(*Config).CheckAuthConfigCanBeEnabled(d.Id(), authOpenLdap)
	if err != nil {
		return err
	}

	if authOpenLdap.Enabled {
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
				return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", d.Id(), err)
			}
		}
		newAuth := &managementClient.OpenLdapConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authOpenLdap, newAuth, managementClient.AuthConfigFieldDisplayName)
		if err != nil {
			return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceRancher2AuthConfigOpenLdapRead(d *schema.ResourceData, meta interface{}) error {
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
	log.Printf("[INFO] Updating Auth Config %s", d.Id())

	authOpenLdap, err := expandAuthConfigOpenLdap(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	err = resourceRancher2AuthConfigOpenLdapApply(d, meta, authOpenLdap)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigOpenLdapRead(d, meta)
}

func resourceRancher2AuthConfigOpenLdapDelete(d *schema.ResourceData, meta interface{}) error {
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
