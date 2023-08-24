package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2ClusterSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2ClusterSyncCreate,
		ReadContext:   resourceRancher2ClusterSyncRead,
		UpdateContext: resourceRancher2ClusterSyncUpdate,
		DeleteContext: resourceRancher2ClusterSyncDelete,

		Schema: clusterSyncFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRancher2ClusterSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID := d.Get("cluster_id").(string)

	cluster, err := meta.(*Config).WaitForClusterState(clusterID, clusterActiveCondition, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	if retries, ok := d.Get("state_confirm").(int); ok && retries > 1 {
		for i := 1; i < retries; i++ {
			time.Sleep(rancher2RetriesWait * time.Second)
			cluster, err = meta.(*Config).WaitForClusterState(clusterID, clusterActiveCondition, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Avoid race condition to generate kube_config for Rancher 2.6 clusters where cluster becomes active before connected
	isRancher26, err := meta.(*Config).IsRancherVersionGreaterThanOrEqual("2.6.0")
	if err != nil {
		return diag.FromErr(err)
	}
	if isRancher26 && cluster.LocalClusterAuthEndpoint != nil && cluster.LocalClusterAuthEndpoint.Enabled {
		// Retrying until resource create timeout
		for connected, _, _ := meta.(*Config).isClusterConnected(clusterID); !connected; connected, _, _ = meta.(*Config).isClusterConnected(clusterID) {
			time.Sleep(rancher2RetriesWait * time.Second)
		}
	}

	if cluster.EnableClusterMonitoring && d.Get("wait_monitoring").(bool) {
		_, err := meta.(*Config).WaitForClusterState(clusterID, clusterMonitoringEnabledCondition, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("[ERROR] waiting for cluster ID (%s) monitoring to be running: %v", clusterID, err)
		}
	}

	if cluster.EnableClusterAlerting && d.Get("wait_alerting").(bool) {
		_, err := meta.(*Config).WaitForClusterState(clusterID, clusterAlertingEnabledCondition, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("[ERROR] waiting for cluster ID (%s) alerting to be running: %v", clusterID, err)
		}
	}

	if d.Get("wait_catalogs").(bool) {
		_, err := waitAllCatalogV2Downloaded(meta.(*Config), clusterID)
		if err != nil {
			return diag.Errorf("[ERROR] waiting for cluster ID (%s) downloading catalogs: %v", clusterID, err)
		}
	}

	d.SetId(clusterID)

	return resourceRancher2ClusterSyncRead(ctx, d, meta)
}

func resourceRancher2ClusterSyncRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID := d.Get("cluster_id").(string)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		active, clus, err := meta.(*Config).isClusterActive(clusterID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster ID %s not found.", clusterID)
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if active {
			defaultProjectID, systemProjectID, err := meta.(*Config).GetClusterSpecialProjectsID(clusterID)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			d.Set("default_project_id", defaultProjectID)
			d.Set("system_project_id", systemProjectID)

			isRancher26, err := meta.(*Config).IsRancherVersionGreaterThanOrEqual("2.6.0")
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if isRancher26 && clus.LocalClusterAuthEndpoint != nil && clus.LocalClusterAuthEndpoint.Enabled {
				connected, _, err := meta.(*Config).isClusterConnected(clusterID)
				if err != nil {
					return retry.NonRetryableError(err)
				}
				if !connected {
					d.Set("synced", false)
					return nil
				}
			}
			kubeConfig, err := getClusterKubeconfig(meta.(*Config), clusterID, d.Get("kube_config").(string))
			if err != nil {
				return retry.NonRetryableError(err)
			}
			d.Set("kube_config", kubeConfig.Config)
			nodes, err := meta.(*Config).GetClusterNodes(clusterID)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			d.Set("nodes", flattenClusterNodes(nodes))

			if clus.EnableClusterMonitoring && d.Get("wait_monitoring").(bool) {
				monitor, _, err := meta.(*Config).isClusterMonitoringEnabledCondition(clusterID)
				if err != nil {
					return retry.NonRetryableError(err)
				}
				if !monitor {
					d.Set("synced", false)
					return nil
				}
			}

			if clus.EnableClusterAlerting && d.Get("wait_alerting").(bool) {
				alert, _, err := meta.(*Config).isClusterAlertingEnabledCondition(clusterID)
				if err != nil {
					return retry.NonRetryableError(err)
				}
				if !alert {
					d.Set("synced", false)
					return nil
				}
			}

			if d.Get("wait_catalogs").(bool) {
				_, err := waitAllCatalogV2Downloaded(meta.(*Config), clusterID)
				if err != nil {
					return retry.NonRetryableError(err)
				}
			}
		}

		d.Set("synced", active)
		return nil
	}))
}

func resourceRancher2ClusterSyncUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRancher2ClusterSyncCreate(ctx, d, meta)
}

func resourceRancher2ClusterSyncDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
