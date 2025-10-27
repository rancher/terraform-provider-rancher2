package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterRKEConfigMonitoringProviderDisabled = "none"
)

//Schemas

func clusterRKEConfigMonitoringFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"node_selector": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"replicas": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"update_strategy": {
			Type:        schema.TypeList,
			Description: "Update deployment strategy",
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: deploymentStrategyFields(),
			},
		},
		"tolerations": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Monitoring add-on tolerations",
			Elem: &schema.Resource{
				Schema: tolerationFields(),
			},
		},
	}
	return s
}
