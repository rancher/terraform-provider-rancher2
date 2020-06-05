package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterAlertGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterAlertGroupRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert group cluster ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert group name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert group description",
			},
			"group_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert group interval seconds",
			},
			"group_wait_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert group wait seconds",
			},
			"recipients": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Alert group recipients",
				Elem: &schema.Resource{
					Schema: recipientFields(),
				},
			},
			"repeat_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert group repeat interval seconds",
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterAlertGroupRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
		"name":      name,
	}
	listOpts := NewListOpts(filters)

	alertGroups, err := client.ClusterAlertGroup.List(listOpts)
	if err != nil {
		return err
	}

	count := len(alertGroups.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster alert group with name \"%s\" on cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster alert group with name \"%s\" on cluster ID \"%s\"", count, name, clusterID)
	}

	return flattenClusterAlertGroup(d, &alertGroups.Data[0])
}
