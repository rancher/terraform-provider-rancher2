package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const clusterSyncMonitoringEnabledCondition = "MonitoringEnabled"

func resourceRancher2ClusterSync() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterSyncCreate,
		Read:   resourceRancher2ClusterSyncRead,
		Update: resourceRancher2ClusterSyncUpdate,
		Delete: resourceRancher2ClusterSyncDelete,

		Schema: clusterSyncFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRancher2ClusterSyncCreate(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)

	start := time.Now()
	cluster, err := meta.(*Config).GetClusterByID(clusterID)
	if err != nil {
		return err
	}
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}
	stateCluster := &resource.StateChangeConf{
		Pending:                   []string{},
		Target:                    []string{"active"},
		Refresh:                   clusterStateRefreshFunc(client, clusterID),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     1 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: d.Get("state_confirm").(int),
	}
	_, waitClusterErr := stateCluster.WaitForState()
	if waitClusterErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster ID (%s) to be active: %s", clusterID, waitClusterErr)
	}

	if cluster.EnableClusterMonitoring && d.Get("wait_monitoring").(bool) {
		enabled := false
		for cluster, err := meta.(*Config).GetClusterByID(clusterID); ; cluster, err = meta.(*Config).GetClusterByID(clusterID) {
			if err != nil {
				return err
			}
			for _, v := range cluster.Conditions {
				if v.Type == clusterSyncMonitoringEnabledCondition {
					if v.Status == "True" {
						enabled = true
					}
					break
				}
			}
			if time.Since(start) >= d.Timeout(schema.TimeoutCreate) || enabled {
				break
			}
			time.Sleep(5 * time.Second)
		}
		if !enabled {
			return fmt.Errorf("[ERROR] waiting for cluster ID (%s) monitoring to be running: Timeout", clusterID)
		}
	}

	d.SetId(clusterID)

	return resourceRancher2ClusterSyncRead(d, meta)
}

func resourceRancher2ClusterSyncRead(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)

	active, clus, err := meta.(*Config).isClusterActive(clusterID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cluster ID %s not found.", clusterID)
			d.SetId("")
			return nil
		}
		return err
	}

	if active {
		defaultProjectID, systemProjectID, err := meta.(*Config).GetClusterSpecialProjectsID(clusterID)
		if err != nil {
			return err
		}
		d.Set("default_project_id", defaultProjectID)
		d.Set("system_project_id", systemProjectID)

		client, err := meta.(*Config).ManagementClient()
		if err != nil {
			return err
		}
		kubeConfig, err := client.Cluster.ActionGenerateKubeconfig(clus)
		if err != nil {
			return err
		}
		d.Set("kube_config", kubeConfig.Config)
	}

	d.Set("synced", active)

	return nil
}

func resourceRancher2ClusterSyncUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRancher2ClusterSyncCreate(d, meta)
}

func resourceRancher2ClusterSyncDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
