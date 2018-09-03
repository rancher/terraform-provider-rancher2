package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func dataSourceRancher2Cluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"api_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_cert": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher Environment: %s", name)

	cluster, err := meta.(*Config).GetClusterByName(name)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing"},
		Target:     []string{"active"},
		Refresh:    findCluster(client, cluster.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for cluster (%s) to be found: %s", name, waitErr)
	}

	d.SetId(cluster.ID)

	d.Set("description", cluster.Description)
	d.Set("name", cluster.Name)
	d.Set("api_endpoint", cluster.APIEndpoint)
	d.Set("ca_cert", cluster.CACert)

	return nil
}

func findCluster(client *managementClient.Client, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clus, err := client.Cluster.ByID(clusterID)
		if err != nil {
			if IsNotFound(err) {
				return clus, "removed", nil
			}
			return nil, "", err
		}

		return clus, clus.State, nil
	}
}
