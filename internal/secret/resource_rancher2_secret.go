package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2Secret() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2SecretCreate,
		Read:   resourceRancher2SecretRead,
		Update: resourceRancher2SecretUpdate,
		Delete: resourceRancher2SecretDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2SecretImport,
		},

		Schema: secretFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2SecretCreate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := meta.(*Config).ProjectExist(projectID)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		secret := expandSecret(d)

		log.Printf("[INFO] Creating Secret %s on Project ID %s", name, projectID)

		newSecret, err := meta.(*Config).CreateSecret(secret)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = flattenSecret(d, newSecret)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = resourceRancher2SecretRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2SecretRead(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Refreshing Secret ID %s", id)

	secret, err := meta.(*Config).GetSecret(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Secret ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenSecret(d, secret)
}

func resourceRancher2SecretUpdate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Updating Secret ID %s", id)

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		secret, err := meta.(*Config).GetSecret(id, projectID, namespaceID)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		update := map[string]interface{}{
			"description": d.Get("description").(string),
			"data":        toMapString(d.Get("data").(map[string]interface{})),
			"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":      toMapString(d.Get("labels").(map[string]interface{})),
		}

		newSecret, err := meta.(*Config).UpdateSecret(secret, update)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = flattenSecret(d, newSecret)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = resourceRancher2SecretRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2SecretDelete(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Deleting Secret ID %s", id)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		secret, err := meta.(*Config).GetSecret(id, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Secret ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		err = meta.(*Config).DeleteSecret(secret)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] Error removing Secret: %s", err))
		}

		d.SetId("")
		return nil
	})
}
