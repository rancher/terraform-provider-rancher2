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

func resourceRancher2NodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2NodePoolCreate,
		ReadContext:   resourceRancher2NodePoolRead,
		UpdateContext: resourceRancher2NodePoolUpdate,
		DeleteContext: resourceRancher2NodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2NodePoolImport,
		},

		Schema: nodePoolFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	nodePool := expandNodePool(d)

	log.Printf("[INFO] Creating Node Pool %s", nodePool.Name)

	err := meta.(*Config).ClusterExist(nodePool.ClusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	newNodePool, err := client.NodePool.Create(nodePool)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newNodePool.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    nodePoolStateRefreshFunc(client, newNodePool.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for node pool (%s) to be created: %s", newNodePool.ID, waitErr)
	}

	return resourceRancher2NodePoolRead(ctx, d, meta)
}

func resourceRancher2NodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Node Pool ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		nodePool, err := client.NodePool.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Node Pool ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenNodePool(d, nodePool); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2NodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Node Pool ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodePool, err := client.NodePool.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	update := map[string]interface{}{
		"hostnamePrefix":          d.Get("hostname_prefix").(string),
		"deleteNotReadyAfterSecs": int64(d.Get("delete_not_ready_after_secs").(int)),
		"drainBeforeDelete":       d.Get("drain_before_delete").(bool),
		"nodeTemplateId":          d.Get("node_template_id").(string),
		"nodeTaints":              expandTaints(d.Get("node_taints").([]interface{})),
		"quantity":                int64(d.Get("quantity").(int)),
		"controlPlane":            d.Get("control_plane").(bool),
		"etcd":                    d.Get("etcd").(bool),
		"worker":                  d.Get("worker").(bool),
		"annotations":             toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                  toMapString(d.Get("labels").(map[string]interface{})),
	}

	newNodePool, err := client.NodePool.Update(nodePool, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    nodePoolStateRefreshFunc(client, newNodePool.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for node pool (%s) to be updated: %s", newNodePool.ID, waitErr)
	}

	return resourceRancher2NodePoolRead(ctx, d, meta)
}

func resourceRancher2NodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Node Pool ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodePool, err := client.NodePool.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Node Pool ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.NodePool.Delete(nodePool)
	if err != nil {
		return diag.Errorf("Error removing Node Pool: %s", err)
	}

	log.Printf("[DEBUG] Waiting for node pool (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    nodePoolStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for node pool (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// nodePoolStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher NodePool.
func nodePoolStateRefreshFunc(client *managementClient.Client, nodePoolID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.NodePool.ByID(nodePoolID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
