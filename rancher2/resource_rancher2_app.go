package rancher2

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
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
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		return err
	}

	err = resourceRancher2AppGetVersion(d, meta)
	if err != nil {
		return err
	}

	app, err := expandApp(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating App %s on Project ID %s", name, projectID)

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return err
	}

	newApp, err := client.App.Create(app)
	if err != nil {
		return err
	}

	if d.Get("wait").(bool) {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"no"},
			Refresh:    appTransitionRefreshFunc(client, newApp.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr := stateConf.WaitForState()
		if waitErr != nil {
			client.App.Delete(newApp)
			return fmt.Errorf("[ERROR] waiting for app (%s) to finish transitioning: %s", newApp.ID, waitErr)
		}
		stateConf = &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    appStateRefreshFunc(client, newApp.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr = stateConf.WaitForState()
		if waitErr != nil {
			client.App.Delete(newApp)
			return fmt.Errorf("[ERROR] waiting for app (%s) to be active: %s", newApp.ID, waitErr)
		}
	}
	d.SetId(newApp.ID)

	return resourceRancher2AppRead(d, meta)
}

func resourceRancher2AppRead(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)
	id := d.Id()

	log.Printf("[INFO] Refreshing App ID %s", id)

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Project ID %s not found.", projectID)
			d.SetId("")
			return nil
		}
		return err
	}

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return err
	}

	app, err := client.App.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] App ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenApp(d, app)
}

func resourceRancher2AppUpdate(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)
	id := d.Id()

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return err
	}

	app, err := client.App.ByID(id)
	if err != nil {
		return err
	}

	if d.HasChange("description") || d.HasChange("annotations") || d.HasChange("labels") {
		log.Printf("[INFO] Updating App ID %s", id)

		app.Description = d.Get("description").(string)
		app.Annotations = toMapString(d.Get("annotations").(map[string]interface{}))
		app.Labels = toMapString(d.Get("labels").(map[string]interface{}))
		_, err := client.App.Replace(app)
		if err != nil {
			return err
		}
	}

	// Rollback or upgrade app
	if d.HasChange("revision_id") {
		revID := d.Get("revision_id").(string)
		log.Printf("[INFO] Rollbacking App ID %s to %s", id, revID)

		rollback := &projectClient.RollbackRevision{
			ForceUpgrade: d.Get("force_upgrade").(bool),
			RevisionID:   revID,
		}
		err = client.App.ActionRollback(app, rollback)
		if err != nil {
			return err
		}
	} else if d.HasChange("answers") || d.HasChange("catalog_name") || d.HasChange("template_name") || d.HasChange("template_version") || d.HasChange("values_yaml") {
		log.Printf("[INFO] Upgrading App ID %s", id)

		values, err := Base64Decode(d.Get("values_yaml").(string))
		if err != nil {
			return err
		}

		upgrade := &projectClient.AppUpgradeConfig{
			Answers:      toMapString(d.Get("answers").(map[string]interface{})),
			ExternalID:   expandAppExternalID(d),
			ForceUpgrade: d.Get("force_upgrade").(bool),
			ValuesYaml:   values,
		}

		err = client.App.ActionUpgrade(app, upgrade)
		if err != nil {
			return err
		}
	}

	if d.Get("wait").(bool) {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"yes"},
			Target:     []string{"no"},
			Refresh:    appTransitionRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr := stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf("[ERROR] waiting for app (%s) to finish transitioning: %s", id, waitErr)
		}
		stateConf = &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    appStateRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr = stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf(
				"[ERROR] waiting for app (%s) to be updated: %s", id, waitErr)
		}
	}

	return resourceRancher2AppRead(d, meta)
}

func resourceRancher2AppDelete(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)
	id := d.Id()

	log.Printf("[INFO] Deleting App ID %s", id)

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return err
	}

	app, err := client.App.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] App ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.App.Delete(app)
	if err != nil {
		return fmt.Errorf("[ERROR] removing App: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    appStateRefreshFunc(client, id),
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

func resourceRancher2AppGetVersion(d *schema.ResourceData, meta interface{}) error {
	catalogName := d.Get("catalog_name").(string)
	appName := d.Get("template_name").(string)
	appVersion := d.Get("template_version").(string)

	if len(appVersion) > 0 {
		return nil
	}

	if !strings.Contains(catalogName, ":") {
		catalogName = "cattle-global-data:" + catalogName
	}

	appID := catalogName + "-" + appName

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	template, err := client.Template.ByID(appID)
	if err != nil {
		return err
	}

	appVersion, err = getLatestVersion(template.VersionLinks)
	if err != nil {
		return err
	}

	d.Set("template_version", appVersion)

	return nil
}

// appStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher App.
func appStateRefreshFunc(client *projectClient.Client, appID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.App.ByID(appID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, obj.State, nil
	}
}

// appTransitionRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher App.
func appTransitionRefreshFunc(client *projectClient.Client, appID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.App.ByID(appID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "no", nil
			}
			return nil, "", err
		}
		return obj, obj.Transitioning, nil
	}
}
