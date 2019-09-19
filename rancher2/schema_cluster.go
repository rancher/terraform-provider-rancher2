package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	clusterDriverImported        = "imported"
	clusterRegistrationTokenName = "system"
)

var (
	clusterDrivers           = []string{clusterDriverImported, clusterDriverAKS, clusterDriverEKS, clusterDriverGKE, clusterDriverRKE}
	clusterPodSecurityPolicy = []string{"restricted", "unrestricted"}
)

//Types

type Cluster struct {
	managementClient.Cluster
	AmazonElasticContainerServiceConfig *AmazonElasticContainerServiceConfig `json:"amazonElasticContainerServiceConfig,omitempty" yaml:"amazonElasticContainerServiceConfig,omitempty"`
	AzureKubernetesServiceConfig        *AzureKubernetesServiceConfig        `json:"azureKubernetesServiceConfig,omitempty" yaml:"azureKubernetesServiceConfig,omitempty"`
	GoogleKubernetesEngineConfig        *GoogleKubernetesEngineConfig        `json:"googleKubernetesEngineConfig,omitempty" yaml:"googleKubernetesEngineConfig,omitempty"`
}

// Schemas

func clusterRegistationTokenFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"insecure_command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"manifest_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"token": &schema.Schema{
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"windows_node_command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}

func clusterAuthEndpoint() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ca_certs": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"fqdn": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}

func clusterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"driver": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(clusterDrivers, true),
		},
		"kube_config": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"rke_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "gke_config"},
			Elem: &schema.Resource{
				Schema: clusterRKEConfigFields(),
			},
		},
		"eks_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "gke_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterEKSConfigFields(),
			},
		},
		"aks_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"eks_config", "gke_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterAKSConfigFields(),
			},
		},
		"gke_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterGKEConfigFields(),
			},
		},
		"default_project_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_project_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"cluster_auth_endpoint": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterAuthEndpoint(),
			},
		},
		"cluster_registration_token": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRegistationTokenFields(),
			},
		},
		"default_pod_security_policy_template_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(clusterPodSecurityPolicy, true),
			Description:  "Default pod security policy template id",
		},
		"enable_cluster_monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable built-in cluster monitoring",
		},
		"enable_network_policy": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable project network isolation",
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
