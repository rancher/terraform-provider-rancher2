package cattle

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func dataSourceCattleCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCattleClusterRead,

		Schema: map[string]*schema.Schema{
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

func dataSourceCattleClusterRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher Environment: %s", name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing", "not found"},
		Target:     []string{"active"},
		Refresh:    findCluster(client, name),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	clus, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for cluster (%s) to be found: %s", name, waitErr)
	}

	cluster := clus.(managementClient.Cluster)
	d.SetId(cluster.ID)

	d.Set("description", cluster.Description)
	d.Set("name", cluster.Name)
	d.Set("api_endpoint", cluster.APIEndpoint)
	d.Set("ca_cert", cluster.CACert)

	return nil
}

func findCluster(client *managementClient.Client, clustername string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusters, err := client.Cluster.List(NewListOpts())
		if err != nil {
			return nil, "", err
		}

		for _, clus := range clusters.Data {
			if clus.Name == clustername {
				return clus, clus.State, nil
			}
		}

		return nil, "not found", nil
	}
}
