package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ProjectAlertGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ProjectAlertGroupCreate,
		Read:   resourceRancher2ProjectAlertGroupRead,
		Update: resourceRancher2ProjectAlertGroupUpdate,
		Delete: resourceRancher2ProjectAlertGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ProjectAlertGroupImport,
		},
		Schema: projectAlertGroupFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ProjectAlertGroupCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceRancher2ProjectAlertGroupRecients(d, meta)
	if err != nil {
		return err
	}
	projectAlertGroup := expandProjectAlertGroup(d)

	log.Printf("[INFO] Creating Project Alert Group %s", projectAlertGroup.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newProjectAlertGroup, err := client.ProjectAlertGroup.Create(projectAlertGroup)
	if err != nil {
		return err
	}

	d.SetId(newProjectAlertGroup.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    projectAlertGroupStateRefreshFunc(client, newProjectAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for project alert group (%s) to be created: %s", newProjectAlertGroup.ID, waitErr)
	}

	return resourceRancher2ProjectAlertGroupRead(d, meta)
}

func resourceRancher2ProjectAlertGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Project Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		projectAlertGroup, err := client.ProjectAlertGroup.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Project Alert Group ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenProjectAlertGroup(d, projectAlertGroup); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ProjectAlertGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Project Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectAlertGroup, err := client.ProjectAlertGroup.ByID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("recipients") {
		err = resourceRancher2ProjectAlertGroupRecients(d, meta)
		if err != nil {
			return err
		}
	}

	update := map[string]interface{}{
		"description":           d.Get("description").(string),
		"groupIntervalSeconds":  int64(d.Get("group_interval_seconds").(int)),
		"groupWaitSeconds":      int64(d.Get("group_wait_seconds").(int)),
		"name":                  d.Get("name").(string),
		"projectId":             d.Get("project_id").(string),
		"recipients":            expandRecipients(d.Get("recipients").([]interface{})),
		"repeatIntervalSeconds": int64(d.Get("repeat_interval_seconds").(int)),
		"annotations":           toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                toMapString(d.Get("labels").(map[string]interface{})),
	}

	newProjectAlertGroup, err := client.ProjectAlertGroup.Update(projectAlertGroup, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    projectAlertGroupStateRefreshFunc(client, newProjectAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project alert group (%s) to be updated: %s", newProjectAlertGroup.ID, waitErr)
	}

	return resourceRancher2ProjectAlertGroupRead(d, meta)
}

func resourceRancher2ProjectAlertGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Project Alert Group ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectAlertGroup, err := client.ProjectAlertGroup.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Project Alert Group ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ProjectAlertGroup.Delete(projectAlertGroup)
	if err != nil {
		return fmt.Errorf("Error removing Project Alert Group: %s", err)
	}

	log.Printf("[DEBUG] Waiting for project alert group (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    projectAlertGroupStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project alert group (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2ProjectAlertGroupRecients(d *schema.ResourceData, meta interface{}) error {
	recipients, ok := d.Get("recipients").([]interface{})
	if !ok {
		return fmt.Errorf("[ERROR] Getting Project Alert Group Recipients")
	}

	if len(recipients) > 0 {
		log.Printf("[INFO] Getting Project Alert Group Recipients")

		for i := range recipients {
			in := recipients[i].(map[string]interface{})

			recipient, err := meta.(*Config).GetRecipientByNotifier(in["notifier_id"].(string))
			if err != nil {
				return err
			}

			in["notifier_type"] = recipient.NotifierType
			if v, ok := in["default_recipient"].(bool); ok && v {
				in["recipient"] = recipient.Recipient
			}

			recipients[i] = in
		}
		d.Set("recipients", recipients)
	}

	return nil
}

// projectAlertGroupStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher ProjectAlertGroup.
func projectAlertGroupStateRefreshFunc(client *managementClient.Client, projectAlertGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ProjectAlertGroup.ByID(projectAlertGroupID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
