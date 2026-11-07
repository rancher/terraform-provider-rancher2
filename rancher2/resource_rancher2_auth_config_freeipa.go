package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigFreeIpa() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigFreeIpaCreate,
		Read:   resourceRancher2AuthConfigFreeIpaRead,
		Update: resourceRancher2AuthConfigFreeIpaUpdate,
		Delete: resourceRancher2AuthConfigFreeIpaDelete,

		Schema: authConfigFreeIpaFields(),
	}
}

func resourceRancher2AuthConfigFreeIpaCreate(d *schema.ResourceData, meta interface{}) error {
	authFreeIpa, err := expandAuthConfigFreeIpa(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	// Create the Auth Config disabled, so that enabling it goes through testAndApply.
	disabled := *authFreeIpa
	disabled.Enabled = false
	_, err = meta.(*Config).CreateAuthConfigIfMissing(authFreeIpa.Name, &disabled, &managementClient.FreeIpaConfig{})
	if err != nil {
		return err
	}

	d.SetId(authFreeIpa.Name)

	err = resourceRancher2AuthConfigFreeIpaApply(d, meta, authFreeIpa)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigFreeIpaRead(d, meta)
}

func resourceRancher2AuthConfigFreeIpaApply(d *schema.ResourceData, meta interface{}, authFreeIpa *managementClient.LdapConfig) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", d.Id(), err)
	}

	err = checkAuthConfigType(d.Id(), auth, authFreeIpa)
	if err != nil {
		return err
	}

	err = meta.(*Config).CheckAuthConfigCanBeEnabled(d.Id(), authFreeIpa)
	if err != nil {
		return err
	}

	if authFreeIpa.Enabled {
		authResource := &norman.Resource{
			ID:      auth.ID,
			Type:    auth.Type,
			Links:   auth.Links,
			Actions: auth.Actions,
		}
		testInput := &managementClient.FreeIpaTestAndApplyInput{
			LdapConfig: authFreeIpa,
			Password:   d.Get("test_password").(string),
			Username:   d.Get("test_username").(string),
		}
		err = client.APIBaseClient.Action(managementClient.FreeIpaConfigType, "testAndApply", authResource, testInput, nil)
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
		newAuth := &managementClient.FreeIpaConfig{}
		err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authFreeIpa, newAuth, managementClient.AuthConfigFieldDisplayName)
		if err != nil {
			return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceRancher2AuthConfigFreeIpaRead(d *schema.ResourceData, meta interface{}) error {
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

	authFreeIpa, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigFreeIpa(d, authFreeIpa.(*managementClient.LdapConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigFreeIpaUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", d.Id())

	authFreeIpa, err := expandAuthConfigFreeIpa(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	err = resourceRancher2AuthConfigFreeIpaApply(d, meta, authFreeIpa)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigFreeIpaRead(d, meta)
}

func resourceRancher2AuthConfigFreeIpaDelete(d *schema.ResourceData, meta interface{}) error {
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
