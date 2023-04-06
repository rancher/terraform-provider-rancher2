package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2GlobalRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2GlobalRoleCreate,
		Read:   resourceRancher2GlobalRoleRead,
		Update: resourceRancher2GlobalRoleUpdate,
		Delete: resourceRancher2GlobalRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2GlobalRoleImport,
		},

		Schema: globalRoleFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2GlobalRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		globalRole := expandGlobalRole(d)

		log.Printf("[INFO] Creating global role")

		newGlobalRole, err := client.GlobalRole.Create(globalRole)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		d.SetId(newGlobalRole.ID)

		err = resourceRancher2GlobalRoleRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2GlobalRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing global role ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		globalRole, err := client.GlobalRole.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] global role ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenGlobalRole(d, globalRole); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2GlobalRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating global role ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		globalRole, err := client.GlobalRole.ByID(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}

		update := map[string]interface{}{
			"description":    d.Get("description").(string),
			"name":           d.Get("name").(string),
			"newUserDefault": d.Get("new_user_default").(bool),
			"rules":          expandPolicyRules(d.Get("rules").([]interface{})),
			"annotations":    toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":         toMapString(d.Get("labels").(map[string]interface{})),
		}

		if _, err = client.GlobalRole.Update(globalRole, update); err != nil {
			return resource.NonRetryableError(err)
		}

		if err = resourceRancher2GlobalRoleRead(d, meta); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2GlobalRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting global role ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		globalRole, err := client.GlobalRole.ByID(id)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Global role ID %s not found.", id)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if !globalRole.Builtin {
			if err = client.GlobalRole.Delete(globalRole); err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Error removing global role: %s", err))
			}
		}

		d.SetId("")
		return nil
	})
}
