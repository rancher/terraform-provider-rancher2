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

func resourceRancher2ClusterAlertGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2ClusterAlertGroupCreate,
		ReadContext:   resourceRancher2ClusterAlertGroupRead,
		UpdateContext: resourceRancher2ClusterAlertGroupUpdate,
		DeleteContext: resourceRancher2ClusterAlertGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2ClusterAlertGroupImport,
		},
		Schema: clusterAlertGroupFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ClusterAlertGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diagnostics := resourceRancher2ClusterAlertGroupRecients(ctx, d, meta)
	if diagnostics.HasError() {
		return diagnostics
	}
	clusterAlertGroup := expandClusterAlertGroup(d)

	log.Printf("[INFO] Creating Cluster Alert Group %s", clusterAlertGroup.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	newClusterAlertGroup, err := client.ClusterAlertGroup.Create(clusterAlertGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newClusterAlertGroup.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    clusterAlertGroupStateRefreshFunc(client, newClusterAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for cluster alert group (%s) to be created: %s", newClusterAlertGroup.ID, waitErr)
	}

	return resourceRancher2ClusterAlertGroupRead(ctx, d, meta)
}

func resourceRancher2ClusterAlertGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Cluster Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		clusterAlertGroup, err := client.ClusterAlertGroup.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster Alert Group ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
		}

		if err = flattenClusterAlertGroup(d, clusterAlertGroup); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2ClusterAlertGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Cluster Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	clusterAlertGroup, err := client.ClusterAlertGroup.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("recipients") {
		diagnostics := resourceRancher2ClusterAlertGroupRecients(ctx, d, meta)
		if diagnostics.HasError() {
			return diagnostics
		}
	}

	update := map[string]interface{}{
		"clusterId":             d.Get("cluster_id").(string),
		"description":           d.Get("description").(string),
		"groupIntervalSeconds":  int64(d.Get("group_interval_seconds").(int)),
		"groupWaitSeconds":      int64(d.Get("group_wait_seconds").(int)),
		"name":                  d.Get("name").(string),
		"recipients":            expandRecipients(d.Get("recipients").([]interface{})),
		"repeatIntervalSeconds": int64(d.Get("repeat_interval_seconds").(int)),
		"annotations":           toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                toMapString(d.Get("labels").(map[string]interface{})),
	}

	newClusterAlertGroup, err := client.ClusterAlertGroup.Update(clusterAlertGroup, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    clusterAlertGroupStateRefreshFunc(client, newClusterAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for cluster alert group (%s) to be updated: %s", newClusterAlertGroup.ID, waitErr)
	}

	return resourceRancher2ClusterAlertGroupRead(ctx, d, meta)
}

func resourceRancher2ClusterAlertGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Cluster Alert Group ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	clusterAlertGroup, err := client.ClusterAlertGroup.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cluster Alert Group ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.ClusterAlertGroup.Delete(clusterAlertGroup)
	if err != nil {
		return diag.Errorf("Error removing Cluster Alert Group: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cluster alert group (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    clusterAlertGroupStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for cluster alert group (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2ClusterAlertGroupRecients(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	recipients, ok := d.Get("recipients").([]interface{})
	if !ok {
		return diag.Errorf("[ERROR] Getting Cluster Alert Group Recipients")
	}

	if len(recipients) > 0 {
		log.Printf("[INFO] Getting Cluster Alert Group Recipients")

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

// clusterAlertGroupStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher ClusterAlertGroup.
func clusterAlertGroupStateRefreshFunc(client *managementClient.Client, clusterAlertGroupID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ClusterAlertGroup.ByID(clusterAlertGroupID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
