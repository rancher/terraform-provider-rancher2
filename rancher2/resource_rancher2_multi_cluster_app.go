package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2MultiClusterApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2MultiClusterAppCreate,
		Read:   resourceRancher2MultiClusterAppRead,
		Update: resourceRancher2MultiClusterAppUpdate,
		Delete: resourceRancher2MultiClusterAppDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2MultiClusterAppImport,
		},

		Schema: multiClusterAppFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2MultiClusterAppCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	err := resourceRancher2AppGetVersion(d, meta)
	if err != nil {
		return err
	}

	multiClusterApp, err := expandMultiClusterApp(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating multi cluster app %s", name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newMultiClusterApp, err := client.MultiClusterApp.Create(multiClusterApp)
	if err != nil {
		return err
	}

	d.SetId(newMultiClusterApp.ID)

	if d.Get("wait").(bool) {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    multiClusterAppStateRefreshFunc(client, newMultiClusterApp.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr := stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf("[ERROR] waiting for multi cluster app (%s) to be created: %s", newMultiClusterApp.ID, waitErr)
		}
	}

	return resourceRancher2MultiClusterAppRead(d, meta)
}

func resourceRancher2MultiClusterAppRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	log.Printf("[INFO] Refreshing multi cluster app ID %s", id)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	multiClusterApp, err := client.MultiClusterApp.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] multi cluster app ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	templateVersion, err := client.TemplateVersion.ByID(multiClusterApp.TemplateVersionID)
	if err != nil {
		return err
	}

	return flattenMultiClusterApp(d, multiClusterApp, templateVersion.ExternalID)
}

func resourceRancher2MultiClusterAppUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	multiClusterApp, err := client.MultiClusterApp.ByID(id)
	if err != nil {
		return err
	}

	updateApp := true

	// Rollback or modify targets
	if d.HasChange("revision_id") {
		updateApp = false
		revID := d.Get("revision_id").(string)
		log.Printf("[INFO] Rollbacking multi cluster app ID %s to %s", id, revID)

		rollback := &managementClient.MultiClusterAppRollbackInput{
			RevisionID: revID,
		}
		err = client.MultiClusterApp.ActionRollback(multiClusterApp, rollback)
		if err != nil {
			return err
		}
	} else if d.HasChange("targets") {
		updateApp = false

		removeTarget := multiClusterAppTargetToRemove(d, multiClusterApp)
		addTarget := multiClusterAppTargetToAdd(d, multiClusterApp)

		if len(removeTarget.Projects) > 0 {
			log.Printf("[INFO] Removing targets on multi cluster app ID %s", id)
			err = client.MultiClusterApp.ActionRemoveProjects(multiClusterApp, removeTarget)
			if err != nil {
				return err
			}
			if d.HasChange("answers") {
				// answer for removed target has to be deleted manually
				updateApp = true
			}
		}

		if len(addTarget.Projects) > 0 {
			log.Printf("[INFO] Adding targets on multi cluster app ID %s", id)
			err = client.MultiClusterApp.ActionAddProjects(multiClusterApp, addTarget)
			if err != nil {
				return err
			}
		}
	}

	// Update app if needed
	if updateApp {
		log.Printf("[INFO] Updating multi cluster app ID %s", id)

		update := map[string]interface{}{
			"answers":              expandAnswers(d.Get("answers").([]interface{})),
			"members":              expandMembers(d.Get("members").([]interface{})),
			"revisionHistoryLimit": d.Get("revision_history_limit").(int),
			"roles":                toArrayString(d.Get("roles").([]interface{})),
			"templateVersionId":    expandMultiClusterAppTemplateVersionID(d),
			"upgradeStrategy":      expandUpgradeStrategy(d.Get("upgrade_strategy").([]interface{})),
			"annotations":          toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":               toMapString(d.Get("labels").(map[string]interface{})),
		}
		_, err := client.MultiClusterApp.Update(multiClusterApp, update)
		if err != nil {
			return err
		}
	}

	if d.Get("wait").(bool) {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    multiClusterAppStateRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr := stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf("[ERROR] waiting for multi cluster app (%s) to be created: %s", id, waitErr)
		}
	}

	return resourceRancher2MultiClusterAppRead(d, meta)
}

func resourceRancher2MultiClusterAppDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	log.Printf("[INFO] Deleting multi cluster app ID %s", id)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	multiClusterApp, err := client.MultiClusterApp.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] multi cluster app ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.MultiClusterApp.Delete(multiClusterApp)
	if err != nil {
		return fmt.Errorf("[ERROR] removing multi cluster app: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    multiClusterAppStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for multi cluster app (%s) to be removed: %s", id, waitErr)
	}
	d.SetId("")

	for i := range multiClusterApp.Targets {
		client, err := meta.(*Config).ProjectClient(multiClusterApp.Targets[i].ProjectID)
		if err != nil {
			continue
		}
		mappID := splitProjectIDPart(multiClusterApp.Targets[i].ProjectID) + ":" + multiClusterApp.Targets[i].AppID
		stateConf = &resource.StateChangeConf{
			Pending:    []string{"removing"},
			Target:     []string{"removed"},
			Refresh:    appStateRefreshFunc(client, mappID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		stateConf.WaitForState()
	}
	time.Sleep(5 * time.Second)

	return nil
}

func resourceRancher2MultiClusterAppGetVersion(d *schema.ResourceData, meta interface{}) error {
	catalogName := d.Get("catalog_name").(string)
	appName := d.Get("template_name").(string)
	appVersion := d.Get("template_version").(string)

	if len(appVersion) > 0 {
		return nil
	}

	catalogName = MultiClusterAppTemplatePrefix + catalogName

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

// multiClusterAppStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher MultiClusterApp.
func multiClusterAppStateRefreshFunc(client *managementClient.Client, appID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.MultiClusterApp.ByID(appID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}

func multiClusterAppTargetToRemove(d *schema.ResourceData, mca *managementClient.MultiClusterApp) *managementClient.UpdateMultiClusterAppTargetsInput {
	newTargets := expandTargets(d.Get("targets").([]interface{}))

	removeTarget := &managementClient.UpdateMultiClusterAppTargetsInput{}
	for _, t := range mca.Targets {
		found := false
		for _, newT := range newTargets {
			if t == newT {
				found = true
				break
			}
		}
		if !found {
			var a *managementClient.Answer
			if d.HasChange("answers") {
				for _, answer := range mca.Answers {
					if t.ProjectID == answer.ProjectID {
						a = &answer
						break
					}
				}
			}
			removeTarget.Projects = append(removeTarget.Projects, t.ProjectID)
			if a != nil {
				removeTarget.Answers = append(removeTarget.Answers, *a)
			}
		}
	}

	return removeTarget
}

func multiClusterAppTargetToAdd(d *schema.ResourceData, mca *managementClient.MultiClusterApp) *managementClient.UpdateMultiClusterAppTargetsInput {
	newTargets := expandTargets(d.Get("targets").([]interface{}))

	addTarget := &managementClient.UpdateMultiClusterAppTargetsInput{}
	for _, newT := range newTargets {
		found := false
		for _, t := range mca.Targets {
			if t == newT {
				found = true
				break
			}
		}
		if !found {
			var a *managementClient.Answer
			if d.HasChange("answers") {
				newAnswers := expandAnswers(d.Get("answers").([]interface{}))
				for _, answer := range newAnswers {
					if newT.ProjectID == answer.ProjectID {
						a = &answer
						break
					}
				}
			}
			addTarget.Projects = append(addTarget.Projects, newT.ProjectID)
			if a != nil {
				addTarget.Answers = append(addTarget.Answers, *a)
			}
		}
	}

	return addTarget
}
