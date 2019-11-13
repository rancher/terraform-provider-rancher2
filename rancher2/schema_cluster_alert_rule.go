package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterAlertRuleFields() map[string]*schema.Schema {
	r := alertRuleFields()
	s := map[string]*schema.Schema{
		"cluster_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert rule cluster ID",
		},
		"event_rule": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"metric_rule", "node_rule", "system_service_rule"},
			Elem: &schema.Resource{
				Schema: eventRuleFields(),
			},
			Description: "Alert event rule",
		},
		"metric_rule": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"event_rule", "node_rule", "system_service_rule"},
			Elem: &schema.Resource{
				Schema: metricRuleFields(),
			},
			Description: "Alert metric rule",
		},
		"node_rule": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"event_rule", "metric_rule", "system_service_rule"},
			Elem: &schema.Resource{
				Schema: nodeRuleFields(),
			},
			Description: "Alert node rule",
		},
		"system_service_rule": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"event_rule", "metric_rule", "node_rule"},
			Elem: &schema.Resource{
				Schema: systemServiceRuleFields(),
			},
			Description: "Alert system service rule",
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}
