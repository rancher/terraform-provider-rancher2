package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2AuthConfigGoogleOauth() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigGoogleOauthCreate,
		Read:   resourceRancher2AuthConfigGoogleOauthRead,
		Update: resourceRancher2AuthConfigGoogleOauthUpdate,
		Delete: resourceRancher2AuthConfigGoogleOauthDelete,

		Schema: authConfigGoogleOauthFields(),
	}
}

func resourceRancher2AuthConfigGoogleOauthCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGoogleOauthName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigGoogleOauthName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigGoogleOauthName, auth.Name)

	authGoogleOauth, err := expandAuthConfigGoogleOauth(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigGoogleOauthName, err)
	}

	// Checking if other auth config is enabled
	if authGoogleOauth.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigGoogleOauthName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigGoogleOauthName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.GoogleOauthConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authGoogleOauth, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigGoogleOauthName, err)
	}

	return resourceRancher2AuthConfigGoogleOauthRead(d, meta)
}

func resourceRancher2AuthConfigGoogleOauthRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigGoogleOauthName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGoogleOauthName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGoogleOauthName)
			d.SetId("")
			return nil
		}
		return err
	}

	authGoogleOauth, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigGoogleOauth(d, authGoogleOauth.(*managementClient.GoogleOauthConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigGoogleOauthUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigGoogleOauthName)

	return resourceRancher2AuthConfigGoogleOauthCreate(d, meta)
}

func resourceRancher2AuthConfigGoogleOauthDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigGoogleOauthName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGoogleOauthName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigGoogleOauthName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigGoogleOauthName, err)
		}
	}

	d.SetId("")
	return nil
}
