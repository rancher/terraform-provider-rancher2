package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterDriver() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterDriverRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"builtin": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"actual_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ui_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"whitelist_domains": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterDriverRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
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
		return err
	}

	count := len(clusterDrivers.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster driver with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster driver with name \"%s\"", count, name)
	}

	return flattenClusterDriver(d, &clusterDrivers.Data[0])
}
