package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Schemas

func nodePoolFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"hostname_prefix": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"node_template_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"quantity": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"control_plane": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"etcd": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"worker": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
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
