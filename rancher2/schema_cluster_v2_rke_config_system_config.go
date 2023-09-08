package rancher2

import (
	"fmt"
	"reflect"

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

func clusterV2RKEConfigSystemConfigFieldsV0() map[string]*schema.Schema {
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
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Machine selector config",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				_, err := ghodssyamlToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in yaml format, error: %v", key, err))
					return
				}
				return
			},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if old == "" || new == "" {
					return false
				}
				oldMap, _ := ghodssyamlToMapInterface(old)
				newMap, _ := ghodssyamlToMapInterface(new)
				return reflect.DeepEqual(oldMap, newMap)
			},
		},
	}

	return s
}
