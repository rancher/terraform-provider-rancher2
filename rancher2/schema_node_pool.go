package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
		"node_taints": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: taintFields(),
			},
		},
		"delete_not_ready_after_secs": &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"control_plane": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"etcd": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"quantity": &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(1),
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
