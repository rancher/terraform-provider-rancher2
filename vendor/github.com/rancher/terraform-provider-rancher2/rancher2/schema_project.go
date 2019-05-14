package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Schemas

func projectResourceQuotaLimitFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"config_maps": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"limits_cpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"limits_memory": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"persistent_volume_claims": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pods": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"replication_controllers": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_cpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_memory": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_storage": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"secrets": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"services": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"services_load_balancers": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"services_node_ports": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}

func projectResourceQuotaFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: projectResourceQuotaLimitFields(),
			},
		},
		"namespace_default_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: projectResourceQuotaLimitFields(),
			},
		},
	}

	return s
}

func projectFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"resource_quota": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: projectResourceQuotaFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["annotations"],
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["labels"],
		},
	}

	return s
}
