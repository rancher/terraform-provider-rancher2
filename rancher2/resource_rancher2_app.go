package rancher2

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

func resourceRancher2App() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AppCreate,
		ReadContext:   resourceRancher2AppRead,
		UpdateContext: resourceRancher2AppUpdate,
		DeleteContext: resourceRancher2AppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2AppImport,
		},

		Schema: appFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2AppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceRancher2AppGetVersion(ctx, d, meta)
	if diagnostics.HasError() {
		return diagnostics
	}

	app, err := expandApp(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating App %s on Project ID %s", name, projectID)

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	newApp, err := client.App.Create(app)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("wait").(bool) {
		stateConf := &retry.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"no"},
			Refresh:    appTransitionRefreshFunc(client, newApp.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr := stateConf.WaitForStateContext(ctx)
		if waitErr != nil {
			client.App.Delete(newApp)
			return diag.Errorf("[ERROR] waiting for app (%s) to finish transitioning: %s", newApp.ID, waitErr)
		}
		stateConf = &retry.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    appStateRefreshFunc(client, newApp.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr = stateConf.WaitForStateContext(ctx)
		if waitErr != nil {
			client.App.Delete(newApp)
			return diag.Errorf("[ERROR] waiting for app (%s) to be active: %s", newApp.ID, waitErr)
		}
	}
	d.SetId(newApp.ID)

	return resourceRancher2AppRead(ctx, d, meta)
}

func resourceRancher2AppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	app, err := client.App.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] App ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenApp(d, app))
}

func resourceRancher2AppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	id := d.Id()

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	app, err := client.App.ByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("description") || d.HasChange("annotations") || d.HasChange("labels") {
		log.Printf("[INFO] Updating App ID %s", id)

		app.Description = d.Get("description").(string)
		app.Annotations = toMapString(d.Get("annotations").(map[string]interface{}))
		app.Labels = toMapString(d.Get("labels").(map[string]interface{}))
		_, err := client.App.Replace(app)
		if err != nil {
			return diag.FromErr(err)
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
			return diag.FromErr(err)
		}
	} else if d.HasChange("answers") || d.HasChange("catalog_name") || d.HasChange("template_name") || d.HasChange("template_version") || d.HasChange("values_yaml") {
		log.Printf("[INFO] Upgrading App ID %s", id)

		values, err := Base64Decode(d.Get("values_yaml").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		upgrade := &projectClient.AppUpgradeConfig{
			Answers:      toMapString(d.Get("answers").(map[string]interface{})),
			ExternalID:   expandAppExternalID(d),
			ForceUpgrade: d.Get("force_upgrade").(bool),
			ValuesYaml:   values,
		}

		err = client.App.ActionUpgrade(app, upgrade)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("wait").(bool) {
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"yes"},
			Target:     []string{"no"},
			Refresh:    appTransitionRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr := stateConf.WaitForStateContext(ctx)
		if waitErr != nil {
			return diag.Errorf("[ERROR] waiting for app (%s) to finish transitioning: %s", id, waitErr)
		}
		stateConf = &retry.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    appStateRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr = stateConf.WaitForStateContext(ctx)
		if waitErr != nil {
			return diag.Errorf(
				"[ERROR] waiting for app (%s) to be updated: %s", id, waitErr)
		}
	}

	return resourceRancher2AppRead(ctx, d, meta)
}

func resourceRancher2AppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	id := d.Id()

	log.Printf("[INFO] Deleting App ID %s", id)

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	app, err := client.App.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] App ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.App.Delete(app)
	if err != nil {
		return diag.Errorf("[ERROR] removing App: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    appStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for App (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2AppGetVersion(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	template, err := client.Template.ByID(appID)
	if err != nil {
		return diag.FromErr(err)
	}

	appVersion, err = getLatestVersion(template.VersionLinks)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("template_version", appVersion)

	return nil
}

// appStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher App.
func appStateRefreshFunc(client *projectClient.Client, appID string) retry.StateRefreshFunc {
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

// appTransitionRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher App.
func appTransitionRefreshFunc(client *projectClient.Client, appID string) retry.StateRefreshFunc {
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
