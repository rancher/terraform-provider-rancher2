package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRancher2App() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AppCreate,
		Read:   resourceRancher2AppRead,
		Update: resourceRancher2AppUpdate,
		Delete: resourceRancher2AppDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2AppImport,
		},

		Schema: appFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2AppCreate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		return err
	}

	app := expandApp(d)

	log.Printf("[INFO] Creating App %s on Project ID %s", name, projectID)

	newApp, err := meta.(*Config).CreateApp(app)
	if err != nil {
		return err
	}

	err = flattenApp(d, newApp)
	if err != nil {
		return err
	}

	return resourceRancher2AppRead(d, meta)
}

func resourceRancher2AppRead(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()

	log.Printf("[INFO] Refreshing App ID %s", id)

	app, err := meta.(*Config).GetApp(id, projectID)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] App ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenApp(d, app)
}

func resourceRancher2AppUpdate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()

	log.Printf("[INFO] Updating App ID %s", id)

	app, err := meta.(*Config).GetApp(id, projectID)
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":            d.Get("name").(string),
		"targetNamespace": d.Get("target_namespace").(string),
		"externalId":      d.Get("external_id").(string),
		"annotations":     toMapString(d.Get("annotations").(map[string]interface{})),
		"answers":         toMapString(d.Get("answers").(map[string]interface{})),
		"description":     d.Get("description").(string),
		"labels":          toMapString(d.Get("labels").(map[string]interface{})),
		"valuesYaml":      d.Get("values_yaml").(string),
	}

	newApp, err := meta.(*Config).UpdateApp(app, update)
	if err != nil {
		return err
	}

	err = flattenApp(d, newApp)
	if err != nil {
		return err
	}

	return resourceRancher2AppRead(d, meta)
}

func resourceRancher2AppDelete(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()

	log.Printf("[INFO] Deleting App ID %s", id)

	app, err := meta.(*Config).GetApp(id, projectID)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] App ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = meta.(*Config).DeleteApp(app)
	if err != nil {
		return fmt.Errorf("Error removing App: %s", err)
	}

	log.Printf("[DEBUG] Waiting for App (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    appStateRefreshFunc(meta, id, projectID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for App (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// appStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher App.
func appStateRefreshFunc(meta interface{}, appID, projectID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := meta.(*Config).GetApp(appID, projectID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
