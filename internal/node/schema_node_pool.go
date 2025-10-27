package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

//Schemas

func nodePoolFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"hostname_prefix": {
			Type:     schema.TypeString,
			Required: true,
		},
		"node_template_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"node_taints": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: taintFields(),
			},
		},
		"delete_not_ready_after_secs": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"drain_before_delete": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"control_plane": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"etcd": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"quantity": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"worker": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
