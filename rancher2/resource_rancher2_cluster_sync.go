package rancher2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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

	active, _, err := meta.(*Config).isClusterActive(clusterID)
	if err != nil {
		return err
	}
	if !active {
		client, err := meta.(*Config).ManagementClient()
		if err != nil {
			return err
		}

		stateCluster := &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    clusterStateRefreshFunc(client, clusterID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitClusterErr := stateCluster.WaitForState()
		if waitClusterErr != nil {
			return fmt.Errorf("[ERROR] waiting for cluster ID (%s) to be active: %s", clusterID, waitClusterErr)
		}
	}

	d.SetId(clusterID)
	d.Set("synced", active)

	return resourceRancher2ClusterSyncRead(d, meta)
}

func resourceRancher2ClusterSyncRead(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)

	active, clus, err := meta.(*Config).isClusterActive(clusterID)
	if err != nil {
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
