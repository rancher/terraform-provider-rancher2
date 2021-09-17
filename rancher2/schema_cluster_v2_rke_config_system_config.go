package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigSystemConfigLabelSelectorExpressionFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Label selector requirement key",
		},
		"operator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Label selector operator",
		},
		"values": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Label selector requirement values",
		},
	}

	return s
}

func clusterV2RKEConfigSystemConfigLabelSelectorFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"match_labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Label selector match labels",
		},
		"match_expressions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Label selector match expressions",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigSystemConfigLabelSelectorExpressionFields(),
			},
		},
	}

	return s
}

func clusterV2RKEConfigSystemConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"machine_label_selector": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Machine label selector",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigSystemConfigLabelSelectorFields(),
			},
		},
		"config": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Machine selector config",
		},
	}

	return s
}
