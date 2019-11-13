package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func projectAlertRuleFields() map[string]*schema.Schema {
	r := alertRuleFields()
	s := map[string]*schema.Schema{
		"project_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert rule Project ID",
		},
		"metric_rule": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"pod_rule", "workload_rule"},
			Elem: &schema.Resource{
				Schema: metricRuleFields(),
			},
			Description: "Alert metric rule",
		},
		"pod_rule": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"metric_rule", "workload_rule"},
			Elem: &schema.Resource{
				Schema: podRuleFields(),
			},
			Description: "Alert pod rule",
		},
		"workload_rule": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"metric_rule", "pod_rule"},
			Elem: &schema.Resource{
				Schema: workloadRuleFields(),
			},
			Description: "Alert workload rule",
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}
