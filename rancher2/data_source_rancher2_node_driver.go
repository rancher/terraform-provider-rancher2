package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2NodeDriver() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2NodeDriverRead,

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
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
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

func dataSourceRancher2NodeDriverRead(d *schema.ResourceData, meta interface{}) error {
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

	nodeDrivers, err := client.NodeDriver.List(listOpts)
	if err != nil {
		return err
	}

	count := len(nodeDrivers.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] node driver with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d node driver with name \"%s\"", count, name)
	}

	return flattenNodeDriver(d, &nodeDrivers.Data[0])
}
