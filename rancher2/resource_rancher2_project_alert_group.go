package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ProjectAlertGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2ProjectAlertGroupCreate,
		ReadContext:   resourceRancher2ProjectAlertGroupRead,
		UpdateContext: resourceRancher2ProjectAlertGroupUpdate,
		DeleteContext: resourceRancher2ProjectAlertGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2ProjectAlertGroupImport,
		},
		Schema: projectAlertGroupFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ProjectAlertGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diag2 := resourceRancher2ProjectAlertGroupRecients(ctx, d, meta)
	if diag2.HasError() {
		return diag2
	}
	projectAlertGroup := expandProjectAlertGroup(d)

	log.Printf("[INFO] Creating Project Alert Group %s", projectAlertGroup.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	newProjectAlertGroup, err := client.ProjectAlertGroup.Create(projectAlertGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newProjectAlertGroup.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    projectAlertGroupStateRefreshFunc(client, newProjectAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for project alert group (%s) to be created: %s", newProjectAlertGroup.ID, waitErr)
	}

	return resourceRancher2ProjectAlertGroupRead(ctx, d, meta)
}

func resourceRancher2ProjectAlertGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Project Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		projectAlertGroup, err := client.ProjectAlertGroup.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Project Alert Group ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenProjectAlertGroup(d, projectAlertGroup); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2ProjectAlertGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Project Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	projectAlertGroup, err := client.ProjectAlertGroup.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("recipients") {
		diag2 := resourceRancher2ProjectAlertGroupRecients(ctx, d, meta)
		if diag2.HasError() {
			return diag2
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
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    projectAlertGroupStateRefreshFunc(client, newProjectAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for project alert group (%s) to be updated: %s", newProjectAlertGroup.ID, waitErr)
	}

	return resourceRancher2ProjectAlertGroupRead(ctx, d, meta)
}

func resourceRancher2ProjectAlertGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Project Alert Group ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	projectAlertGroup, err := client.ProjectAlertGroup.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Project Alert Group ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.ProjectAlertGroup.Delete(projectAlertGroup)
	if err != nil {
		return diag.Errorf("Error removing Project Alert Group: %s", err)
	}

	log.Printf("[DEBUG] Waiting for project alert group (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    projectAlertGroupStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for project alert group (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2ProjectAlertGroupRecients(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	recipients, ok := d.Get("recipients").([]interface{})
	if !ok {
		return diag.Errorf("[ERROR] Getting Project Alert Group Recipients")
	}

	if len(recipients) > 0 {
		log.Printf("[INFO] Getting Project Alert Group Recipients")

		for i := range recipients {
			in := recipients[i].(map[string]interface{})

			recipient, err := meta.(*Config).GetRecipientByNotifier(in["notifier_id"].(string))
			if err != nil {
				return diag.FromErr(err)
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

// projectAlertGroupStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher ProjectAlertGroup.
func projectAlertGroupStateRefreshFunc(client *managementClient.Client, projectAlertGroupID string) retry.StateRefreshFunc {
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
