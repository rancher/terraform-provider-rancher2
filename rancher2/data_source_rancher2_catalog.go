package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRancher2Catalog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2CatalogRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"branch": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceRancher2CatalogRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	filters := map[string]interface{}{"name": name}
	listOpts := NewListOpts(filters)

	catalogs, err := client.Catalog.List(listOpts)
	if err != nil {
		return err
	}

	count := len(catalogs.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] catalog with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d catalogs with name \"%s\"", count, name)
	}

	return flattenCatalog(d, &catalogs.Data[0])
}
