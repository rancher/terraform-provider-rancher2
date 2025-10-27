package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func namespaceResourceQuotaLimitFields() map[string]*schema.Schema {
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

func namespaceResourceQuotaFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: namespaceResourceQuotaLimitFields(),
			},
		},
	}

	return s
}

func namespaceFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Project ID where k8s namespace belongs",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if old == "" || new == "" {
					return false
				}
				oldClusterID, oldProjectID := splitProjectID(old)
				newClusterID, newProjectID := splitProjectID(new)
				// Just update project_id inside same cluster ID
				if oldClusterID == newClusterID && oldProjectID != newProjectID {
					return false
				}
				return true
			},
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the k8s namespace managed by rancher v2",
		},
		"container_resource_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: containerResourceLimitFields(),
			},
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the k8s namespace managed by rancher v2",
		},
		"resource_quota": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: namespaceResourceQuotaFields(),
			},
		},
		"wait_for_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Wait for cluster becomes active",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
