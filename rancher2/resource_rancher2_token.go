package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2Token() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2TokenCreate,
		Read:   resourceRancher2TokenRead,
		Update: resourceRancher2TokenUpdate,
		Delete: resourceRancher2TokenDelete,

		Schema: tokenFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2TokenCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Token")
	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return err
	}
	token, err := expandToken(d, patch)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newToken, err := client.Token.Create(token)
	if err != nil {
		return err
	}

	err = flattenToken(d, newToken, patch)
	if err != nil {
		return err
	}

	return resourceRancher2TokenRead(d, meta)
}

func resourceRancher2TokenRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Token ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	token, err := client.Token.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Token ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	renew := d.Get("renew").(bool)
	if (!*token.Enabled || token.Expired) && renew {
		d.Set("renew", false)
	}

	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return err
	}
	err = flattenToken(d, token, patch)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2TokenUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRancher2TokenRead(d, meta)
}

func resourceRancher2TokenDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Token ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	token, err := client.Token.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Token ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Token.Delete(token)
	if err != nil {
		return fmt.Errorf("Error removing Token: %s", err)
	}

	d.SetId("")
	return nil
}
