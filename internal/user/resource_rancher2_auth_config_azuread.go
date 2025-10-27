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
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigAzureADName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigAzureADName, err)
	}

	log.Printf("[INFO] Creating Auth Config AzureAD %s", auth.Name)

	authAzureAD, err := expandAuthConfigAzureAD(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigAzureADName, err)
	}

	// Checking if other auth config is enabled
	if authAzureAD.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigAzureADName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigAzureADName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.AzureADConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authAzureAD, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigAzureADName, err)
	}

	return resourceRancher2AuthConfigAzureADRead(d, meta)
}

func resourceRancher2AuthConfigAzureADRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigAzureADName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigAzureADName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigAzureADName)
			d.SetId("")
			return nil
		}
		return err
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
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigAzureADName)

	return resourceRancher2AuthConfigAzureADCreate(d, meta)
}

func resourceRancher2AuthConfigAzureADDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigAzureADName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigAzureADName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigAzureADName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigAzureADName, err)
		}
	}

	d.SetId("")
	return nil
}
