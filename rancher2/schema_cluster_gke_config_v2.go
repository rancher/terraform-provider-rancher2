package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterGKEV2Kind                     = "gkeV2"
	clusterDriverGKEV2                   = "GKE"
	clusterGKEV2LoggingAudit             = "audit"
	clusterGKEV2LoggingAPI               = "api"
	clusterGKEV2LoggingScheduler         = "scheduler"
	clusterGKEV2LoggingcontrollerManager = "controllerManager"
	clusterGKEV2LoggingAuthenticator     = "authenticator"
)

//Schemas

func clusterGKEConfigV2NodeTaintFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"effect": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	return s
}

func clusterGKEConfigV2ClusterAddonsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"http_load_balancing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable GKE HTTP load balancing",
		},
		"horizontal_pod_autoscaling": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable GKE horizontal pod autoscaling",
		},
		"network_policy_config": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable GKE network policy config",
		},
	}

	return s
}

func clusterGKEConfigV2IPAllocationPolicyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster ip v4 allocation cidr block",
		},
		"cluster_secondary_range_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster ip v4 allocation secondary range name",
		},
		"create_subnetwork": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Create GKE subnetwork?",
		},
		"node_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE node ip v4 allocation cidr block",
		},
		"services_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE services ip v4 allocation cidr block",
		},
		"services_secondary_range_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE services ip v4 allocation secondary range name",
		},
		"subnetwork_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster subnetwork name",
		},
		"use_ip_aliases": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Use GKE ip aliases?",
		},
	}

	return s
}

func clusterGKEConfigV2CidrBlocksFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cidr_block": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The GKE master authorized network config cidr block",
		},
		"display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The GKE master authorized network config cidr block dispaly name",
		},
	}

	return s
}

func clusterGKEConfigV2MasterAuthorizedNetworksConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cidr_blocks": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The GKE master authorized network config cidr blocks",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2CidrBlocksFields(),
			},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable GKE master authorized network config",
		},
	}

	return s
}

func clusterGKEConfigV2NodePoolConfigAutoscalingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable GKE node pool config autoscaling",
		},
		"max_node_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The GKE node pool config max node count",
		},
		"min_node_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The GKE node pool config min node count",
		},
	}

	return s
}

func clusterGKEConfigV2NodeConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"disk_size_gb": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The GKE node config disk size (Gb)",
		},
		"disk_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The GKE node config disk type",
		},
		"image_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The GKE node config image type",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "The GKE node config labels",
		},
		"local_ssd_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The GKE node config local ssd count",
		},
		"machine_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The GKE node config machine type",
		},
		"oauth_scopes": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The GKE node config oauth scopes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"preemptible": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable GKE node config preemptible",
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The GKE node config tags",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"taints": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The GKE node config taints",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2NodeTaintFields(),
			},
		},
	}

	return s
}

func clusterGKEConfigV2NodePoolConfigManagementFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"auto_repair": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable GKE node pool config management auto repair",
		},
		"auto_upgrade": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable GKE node pool config management auto upgrade",
		},
	}

	return s
}

func clusterGKEConfigV2NodePoolConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The GKE node pool config name",
		},
		"initial_node_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The GKE node pool config initial node count",
		},
		"version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The GKE node pool config version",
		},
		"autoscaling": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "The GKE node pool config autoscaling",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2NodePoolConfigAutoscalingFields(),
			},
		},
		"config": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "The GKE node pool node config",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2NodeConfigFields(),
			},
		},
		"management": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "The GKE node pool config management",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2NodePoolConfigManagementFields(),
			},
		},
		"max_pods_constraint": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The GKE node pool config max pods constraint",
		},
	}

	return s
}

func clusterGKEConfigV2PrivateClusterConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"master_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The GKE cluster private master ip v4 cidr block",
		},
		"enable_private_endpoint": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable GKE cluster private endpoint",
		},
		"enable_private_nodes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable GKE cluster private nodes",
		},
	}

	return s
}

func clusterGKEConfigV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The GKE cluster name",
		},
		"google_credential_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Google credential secret",
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The GKE project id",
		},
		"cluster_addons": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "The GKE cluster addons",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2ClusterAddonsFields(),
			},
		},
		"cluster_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "The GKE ip v4 cidr block",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "The GKE cluster description",
		},
		"enable_kubernetes_alpha": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Enable Kubernetes alpha",
		},
		"ip_allocation_policy": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "The GKE ip allocation policy",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2IPAllocationPolicyFields(),
			},
		},
		"imported": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Is GKE cluster imported?",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The kubernetes master version",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster labels",
		},
		"locations": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster locations",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"logging_service": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster logging service",
		},
		"maintenance_window": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster maintenance window",
		},
		"master_authorized_networks_config": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "The GKE cluster master authorized networks config",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2MasterAuthorizedNetworksConfigFields(),
			},
		},
		"monitoring_service": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster monitoring service",
		},
		"network": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "The GKE cluster network",
		},
		"network_policy_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Is GKE cluster network policy enabled?",
		},
		"node_pools": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The GKE cluster node pools",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2NodePoolConfigFields(),
			},
		},
		"private_cluster_config": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "The GKE private cluster config",
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2PrivateClusterConfigFields(),
			},
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "The GKE cluster region. Required if `zone` is empty",
		},
		"subnetwork": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "The GKE cluster subnetwork",
		},
		"zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "The GKE cluster zone. Required if `region` is empty",
		},
	}

	return s
}
