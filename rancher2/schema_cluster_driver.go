package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterDriverFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"active": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
		"builtin": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"actual_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"checksum": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"ui_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"whitelist_domains": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
