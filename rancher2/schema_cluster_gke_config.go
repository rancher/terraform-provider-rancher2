package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	clusterGKEKind = "gke"
)

//Schemas

func clusterGKEConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_ipv4_cidr": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Cluster ipv4 CIDR for GKE",
		},
		"credential": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Required Credential for GKE",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Required Description for GKE cluster",
		},
		"disk_size_gb": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Required Disk size for agents for GKE cluster",
		},
		"enable_alpha_feature": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable alpha features on GKE cluster",
		},
		"enable_http_load_balancing": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable HTTP load balancing on GKE cluster",
		},
		"enable_horizontal_pod_autoscaling": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable Horitzontal Pod Autoscaling on GKE cluster",
		},
		"enable_kubernetes_dashboard": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable kubernetes dashboard on GKE cluster",
		},
		"enable_legacy_abac": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable legacy abac on GKE cluster",
		},
		"enable_network_policy_config": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable network policy config on GKE cluster",
		},
		"enable_stackdriver_logging": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable stackdriver logging on GKE cluster",
		},
		"enable_stackdriver_monitoring": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required Enable stackdriver monitoring on GKE cluster",
		},
		"image_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Image type for GKE cluster",
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Optional Labels for GKE",
		},
		"locations": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Required Locations for GKE cluster",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"machine_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Machine type for GKE cluster",
		},
		"maintenance_window": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Maintenance window for GKE cluster",
		},
		"master_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Master version for GKE cluster",
		},
		"network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Network for GKE cluster",
		},
		"node_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Required Node count for GKE cluster",
		},
		"node_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Node version for GKE cluster",
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Project ID for GKE cluster",
		},
		"sub_network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Subnetwork for GKE cluster",
		},
		"zone": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Zone GKE cluster",
		},
	}

	return s
}
