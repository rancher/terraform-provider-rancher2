package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigADFS() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigADFSCreate,
		Read:   resourceRancher2AuthConfigADFSRead,
		Update: resourceRancher2AuthConfigADFSUpdate,
		Delete: resourceRancher2AuthConfigADFSDelete,

		Schema: authConfigADFSFields(),
	}
}

func resourceRancher2AuthConfigADFSCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigADFSName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigADFSName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigADFSName, auth.Name)

	authADFS, err := expandAuthConfigADFS(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigADFSName, err)
	}

	// Checking if other auth config is enabled
	if authADFS.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigADFSName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigADFSName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.ADFSConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authADFS, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigADFSName, err)
	}

	return resourceRancher2AuthConfigADFSRead(d, meta)
}

func resourceRancher2AuthConfigADFSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigADFSName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigADFSName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigADFSName)
			d.SetId("")
			return nil
		}
		return err
	}

	authADFS, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigADFS(d, authADFS.(*managementClient.ADFSConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigADFSUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigADFSName)

	return resourceRancher2AuthConfigADFSCreate(d, meta)
}

func resourceRancher2AuthConfigADFSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigADFSName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigADFSName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigADFSName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigADFSName, err)
		}
	}

	d.SetId("")
	return nil
}
