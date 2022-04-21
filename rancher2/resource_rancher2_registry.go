package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2Registry() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2RegistryCreate,
		Read:   resourceRancher2RegistryRead,
		Update: resourceRancher2RegistryUpdate,
		Delete: resourceRancher2RegistryDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2RegistryImport,
		},

		Schema: registryFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2RegistryCreate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		return err
	}

	registry := expandRegistry(d)

	log.Printf("[INFO] Creating Registry %s on Project ID %s", name, projectID)

	newRegistry, err := meta.(*Config).CreateRegistry(registry)
	if err != nil {
		return err
	}

	err = flattenRegistry(d, newRegistry)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    registryStateRefreshFunc(meta, d.Id(), projectID, d.Get("namespace_id").(string)),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for registry (%s) to be created: %s", d.Id(), waitErr)
	}

	return resourceRancher2RegistryRead(d, meta)
}

func resourceRancher2RegistryRead(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Refreshing Registry ID %s", id)

	registry, err := meta.(*Config).GetRegistry(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Registry ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenRegistry(d, registry)
}

func resourceRancher2RegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Updating Registry ID %s", id)

	registry, err := meta.(*Config).GetRegistry(id, projectID, namespaceID)
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"description": d.Get("description").(string),
		"registries":  expandRegistryCredential(d.Get("registries").([]interface{})),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newRegistry, err := meta.(*Config).UpdateRegistry(registry, update)
	if err != nil {
		return err
	}

	err = flattenRegistry(d, newRegistry)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    registryStateRefreshFunc(meta, id, projectID, namespaceID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for registry (%s) to be updated: %s", id, waitErr)
	}

	return resourceRancher2RegistryRead(d, meta)
}

func resourceRancher2RegistryDelete(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Deleting Registry ID %s", id)

	registry, err := meta.(*Config).GetRegistry(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Registry ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = meta.(*Config).DeleteRegistry(registry)
	if err != nil {
		return fmt.Errorf("Error removing Registry: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    registryStateRefreshFunc(meta, id, projectID, namespaceID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for registry (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func registryStateRefreshFunc(meta interface{}, id string, projectID string, namespaceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := meta.(*Config).GetRegistry(id, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "exists", nil
	}
}
