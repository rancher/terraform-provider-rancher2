package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ClusterAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterAlertRuleCreate,
		Read:   resourceRancher2ClusterAlertRuleRead,
		Update: resourceRancher2ClusterAlertRuleUpdate,
		Delete: resourceRancher2ClusterAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterAlertRuleImport,
		},
		Schema: clusterAlertRuleFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ClusterAlertRuleCreate(d *schema.ResourceData, meta interface{}) error {
	clusterAlertRule := expandClusterAlertRule(d)

	log.Printf("[INFO] Creating Cluster Alert Rule %s", clusterAlertRule.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newClusterAlertRule, err := client.ClusterAlertRule.Create(clusterAlertRule)
	if err != nil {
		return err
	}

	d.SetId(newClusterAlertRule.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    clusterAlertRuleStateRefreshFunc(client, newClusterAlertRule.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster alert rule (%s) to be created: %s", newClusterAlertRule.ID, waitErr)
	}

	return resourceRancher2ClusterAlertRuleRead(d, meta)
}

func resourceRancher2ClusterAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster Alert Rule ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		clusterAlertRule, err := client.ClusterAlertRule.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster Alert Rule ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}

			return resource.NonRetryableError(err)
		}

		if err = flattenClusterAlertRule(d, clusterAlertRule); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ClusterAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cluster Alert Rule ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterAlertRule, err := client.ClusterAlertRule.ByID(d.Id())
	if err != nil {
		return err
	}

	inherited := d.Get("inherited").(bool)
	update := map[string]interface{}{
		"clusterId":             d.Get("cluster_id").(string),
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

	if v, ok := d.Get("event_rule").([]interface{}); ok && len(v) > 0 {
		update["eventRule"] = expandEventRule(v)
	}

	if v, ok := d.Get("metric_rule").([]interface{}); ok && len(v) > 0 {
		update["metricRule"] = expandMetricRule(v)
	}

	if v, ok := d.Get("node_rule").([]interface{}); ok && len(v) > 0 {
		update["nodeRule"] = expandNodeRule(v)
	}

	if v, ok := d.Get("system_service_rule").([]interface{}); ok && len(v) > 0 {
		update["systemServiceRule"] = expandSystemServiceRule(v)
	}

	newClusterAlertRule, err := client.ClusterAlertRule.Update(clusterAlertRule, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    clusterAlertRuleStateRefreshFunc(client, newClusterAlertRule.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster alert rule (%s) to be updated: %s", newClusterAlertRule.ID, waitErr)
	}

	return resourceRancher2ClusterAlertRuleRead(d, meta)
}

func resourceRancher2ClusterAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster Alert Rule ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterAlertRule, err := client.ClusterAlertRule.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cluster Alert Rule ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ClusterAlertRule.Delete(clusterAlertRule)
	if err != nil {
		return fmt.Errorf("Error removing Cluster Alert Rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cluster alert rule (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    clusterAlertRuleStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster alert rule (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// clusterAlertRuleStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher ClusterAlertRule.
func clusterAlertRuleStateRefreshFunc(client *managementClient.Client, clusterAlertRuleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ClusterAlertRule.ByID(clusterAlertRuleID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
