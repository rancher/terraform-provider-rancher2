package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterScan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterScanRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster ID to scan",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The cluster scan name",
			},
			"run_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster scan run type",
			},
			"scan_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterScanConfigFields(),
				},
				Description: "The cluster scan config",
			},
			"scan_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster scan type",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster scan status",
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

func dataSourceRancher2ClusterScanRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
	}
	if len(name) > 0 {
		filters["name"] = name
	}
	listOpts := NewListOpts(filters)

	clusterScans, err := client.ClusterScan.List(listOpts)
	if err != nil {
		return err
	}

	count := len(clusterScans.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster scan with cluster ID \"%s\" not found", clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster scan with cluster ID \"%s\"", count, clusterID)
	}

	d.SetId(clusterScans.Data[0].ID)

	return flattenClusterScan(d, &clusterScans.Data[0])
}
