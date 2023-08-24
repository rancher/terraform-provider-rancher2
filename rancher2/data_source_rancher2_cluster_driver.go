package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2ClusterDriver() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2ClusterDriverRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"builtin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"actual_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ui_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"whitelist_domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func dataSourceRancher2ClusterDriverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	url := d.Get("url").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	if len(url) > 0 {
		filters["url"] = url
	}
	listOpts := NewListOpts(filters)

	clusterDrivers, err := client.KontainerDriver.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(clusterDrivers.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] cluster driver with name \"%s\" not found", name)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d cluster driver with name \"%s\"", count, name)
	}

	return diag.FromErr(flattenClusterDriver(d, &clusterDrivers.Data[0]))
}
