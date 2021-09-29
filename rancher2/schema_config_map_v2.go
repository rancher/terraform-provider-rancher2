package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func configMapV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "K8s cluster ID",
		},
		"data": {
			Type:        schema.TypeMap,
			Required:    true,
			Description: "ConfigMap V2 data map",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ConfigMap V2 name",
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "default",
			Description: "ConfigMap V2 namespace",
		},
		"immutable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If set to true, ensures that data stored in the ConfigMap cannot be updated",
		},
		"resource_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
