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
	authADFS, err := expandAuthConfigADFS(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config: %s", err)
	}

	newAuth := &managementClient.ADFSConfig{}
	err = meta.(*Config).CreateOrUpdateAuthConfig(authADFS.Name, authADFS, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	d.SetId(authADFS.Name)

	return resourceRancher2AuthConfigADFSRead(d, meta)
}

func resourceRancher2AuthConfigADFSRead(d *schema.ResourceData, meta interface{}) error {
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
	authADFS, err := expandAuthConfigADFS(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", d.Id(), err)
	}

	newAuth := &managementClient.ADFSConfig{}
	err = meta.(*Config).UpdateAuthConfigByName(d.Id(), authADFS, newAuth, managementClient.AuthConfigFieldDisplayName)
	if err != nil {
		return err
	}

	return resourceRancher2AuthConfigADFSRead(d, meta)
}

func resourceRancher2AuthConfigADFSDelete(d *schema.ResourceData, meta interface{}) error {
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
