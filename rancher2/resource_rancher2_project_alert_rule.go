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

func resourceRancher2ProjectAlertRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2ProjectAlertRuleCreate,
		ReadContext:   resourceRancher2ProjectAlertRuleRead,
		UpdateContext: resourceRancher2ProjectAlertRuleUpdate,
		DeleteContext: resourceRancher2ProjectAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2ProjectAlertRuleImport,
		},
		Schema: projectAlertRuleFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ProjectAlertRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectAlertRule := expandProjectAlertRule(d)

	log.Printf("[INFO] Creating Project Alert Rule %s", projectAlertRule.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	newProjectAlertRule, err := client.ProjectAlertRule.Create(projectAlertRule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newProjectAlertRule.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    projectAlertRuleStateRefreshFunc(client, newProjectAlertRule.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for project alert rule (%s) to be created: %s", newProjectAlertRule.ID, waitErr)
	}

	return resourceRancher2ProjectAlertRuleRead(ctx, d, meta)
}

func resourceRancher2ProjectAlertRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Project Alert Rule ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		projectAlertRule, err := client.ProjectAlertRule.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Project Alert Rule ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenProjectAlertRule(d, projectAlertRule); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2ProjectAlertRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Project Alert Rule ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	projectAlertRule, err := client.ProjectAlertRule.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	inherited := d.Get("inherited").(bool)
	update := map[string]interface{}{
		"projectId":             d.Get("project_id").(string),
		"groupId":               d.Get("group_id").(string),
		"groupIntervalSeconds":  int64(d.Get("group_interval_seconds").(int)),
		"groupWaitSeconds":      int64(d.Get("group_wait_seconds").(int)),
		"inherited":             &inherited,
		"name":                  d.Get("name").(string),
		"repeatIntervalSeconds": int64(d.Get("repeat_interval_seconds").(int)),
		"severity":              d.Get("severity").(string),
		"annotations":           toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                toMapString(d.Get("labels").(map[string]interface{})),
	}

	if v, ok := d.Get("metric_rule").([]interface{}); ok && len(v) > 0 {
		update["metricRule"] = expandMetricRule(v)
	}

	if v, ok := d.Get("pod_rule").([]interface{}); ok && len(v) > 0 {
		update["podRule"] = expandPodRule(v)
	}

	if v, ok := d.Get("workload_rule").([]interface{}); ok && len(v) > 0 {
		update["workloadRule"] = expandWorkloadRule(v)
	}

	newProjectAlertRule, err := client.ProjectAlertRule.Update(projectAlertRule, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    projectAlertRuleStateRefreshFunc(client, newProjectAlertRule.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for project alert rule (%s) to be updated: %s", newProjectAlertRule.ID, waitErr)
	}

	return resourceRancher2ProjectAlertRuleRead(ctx, d, meta)
}

func resourceRancher2ProjectAlertRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Project Alert Rule ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	projectAlertRule, err := client.ProjectAlertRule.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Project Alert Rule ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.ProjectAlertRule.Delete(projectAlertRule)
	if err != nil {
		return diag.Errorf("Error removing Project Alert Rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for project alert rule (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    projectAlertRuleStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for project alert rule (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// projectAlertRuleStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher ProjectAlertRule.
func projectAlertRuleStateRefreshFunc(client *managementClient.Client, projectAlertRuleID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ProjectAlertRule.ByID(projectAlertRuleID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
