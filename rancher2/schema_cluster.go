package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	clusterDriverImported             = "imported"
	clusterRegistrationTokenName      = "default-token"
	clusterMonitoringV1Namespace      = "cattle-prometheus"
	clusterActiveCondition            = "Updated"
	clusterConnectedCondition         = "Connected"
	clusterMonitoringEnabledCondition = "MonitoringEnabled"
	clusterAlertingEnabledCondition   = "AlertingEnabled"
)

var (
	clusterDrivers                = []string{clusterDriverImported, clusterDriverAKS, clusterDriverEKS, clusterDriverGKE, clusterDriverGKEV2, clusterDriverK3S, clusterDriverOKE, clusterDriverRKE, clusterDriverRKE2}
	clusterRegistrationTokenNames = []string{clusterRegistrationTokenName, "system"}
)

//Types

type Cluster struct {
	managementClient.Cluster
	AmazonElasticContainerServiceConfig *AmazonElasticContainerServiceConfig `json:"amazonElasticContainerServiceConfig,omitempty" yaml:"amazonElasticContainerServiceConfig,omitempty"`
	AzureKubernetesServiceConfig        *AzureKubernetesServiceConfig        `json:"azureKubernetesServiceConfig,omitempty" yaml:"azureKubernetesServiceConfig,omitempty"`
	GoogleKubernetesEngineConfig        *GoogleKubernetesEngineConfig        `json:"googleKubernetesEngineConfig,omitempty" yaml:"googleKubernetesEngineConfig,omitempty"`
	OracleKubernetesEngineConfig        *OracleKubernetesEngineConfig        `json:"okeEngineConfig,omitempty" yaml:"okeEngineConfig,omitempty"`
}

// Schemas

func clusterRegistationTokenFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cluster_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"command": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"insecure_command": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"insecure_node_command": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"insecure_windows_node_command": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"manifest_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_command": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"token": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"windows_node_command": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func clusterAuthEndpoint() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ca_certs": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"fqdn": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}

func clusterDataFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"driver": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"kube_config": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rke_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigFieldsV0(),
			},
		},
		"k3s_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterK3SConfigFields(),
			},
		},
		"eks_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterEKSConfigFields(),
			},
		},
		"aks_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterAKSConfigFields(),
			},
		},
		"gke_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterGKEConfigFields(),
			},
		},
		"default_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cluster_auth_endpoint": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterAuthEndpoint(),
			},
		},
		"cluster_monitoring_input": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Computed:    true,
			Description: "Cluster monitoring configuration",
			Elem: &schema.Resource{
				Schema: monitoringInputFields(),
			},
		},
		"cluster_registration_token": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRegistationTokenFields(),
			},
		},
		"cluster_template_answers": {
			Type:        schema.TypeList,
			Computed:    true,
			MaxItems:    1,
			Description: "Cluster template answers",
			Elem: &schema.Resource{
				Schema: answerFields(),
			},
		},
		"cluster_template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template ID",
		},
		"cluster_template_questions": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Cluster template questions",
			Elem: &schema.Resource{
				Schema: questionFields(),
			},
		},
		"cluster_template_revision_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template revision ID",
		},
		"default_pod_security_policy_template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Default pod security policy template id",
		},
		"enable_cluster_alerting": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable built-in cluster alerting",
		},
		"enable_cluster_monitoring": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable built-in cluster monitoring",
		},
		"enable_network_policy": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable project network isolation",
		},
		"annotations": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"labels": {
			Type:     schema.TypeMap,
			Computed: true,
		},
	}
	return s
}

func clusterFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"driver": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(clusterDrivers, true),
		},
		"kube_config": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"rke_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "gke_config", "k3s_config"},
			Elem: &schema.Resource{
				Schema: clusterRKEConfigFieldsV0(),
			},
		},
		"k3s_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "gke_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterK3SConfigFields(),
			},
		},
		"eks_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "gke_config", "k3s_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterEKSConfigFields(),
			},
		},
		"aks_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"eks_config", "gke_config", "k3s_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterAKSConfigFields(),
			},
		},
		"gke_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "k3s_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterGKEConfigFields(),
			},
		},
		"default_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cluster_auth_endpoint": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterAuthEndpoint(),
			},
		},
		"cluster_monitoring_input": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster monitoring configuration",
			Elem: &schema.Resource{
				Schema: monitoringInputFields(),
			},
		},
		"cluster_registration_token": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRegistationTokenFields(),
			},
		},
		"cluster_template_answers": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Computed:    true,
			Description: "Cluster template answers",
			Elem: &schema.Resource{
				Schema: answerFields(),
			},
		},
		"cluster_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster template ID",
		},
		"cluster_template_questions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster template questions",
			Elem: &schema.Resource{
				Schema: questionFields(),
			},
		},
		"cluster_template_revision_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster template revision ID",
		},
		"default_pod_security_policy_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default pod security policy template id",
		},
		"desired_agent_image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"desired_auth_image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"docker_root_dir": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"enable_cluster_alerting": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable built-in cluster alerting",
		},
		"enable_cluster_monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable built-in cluster monitoring",
		},
		"enable_cluster_istio": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable built-in cluster istio",
			Deprecated:  "Deploy istio using rancher2_app resource instead",
		},
		"istio_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Is istio enabled at cluster?",
		},
		"enable_network_policy": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable project network isolation",
		},
		"windows_prefered_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Windows preferred cluster",
			ForceNew:    true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func clusterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"agent_env_vars": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional Agent Env Vars for Rancher agent",
			Elem: &schema.Resource{
				Schema: envVarFields(),
			},
		},
		"driver": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(clusterDrivers, true),
		},
		"kube_config": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"ca_cert": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"rke_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config", "eks_config_v2", "gke_config", "gke_config_v2", "k3s_config", "oke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterRKEConfigFields(),
			},
		},
		"rke2_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config", "eks_config_v2", "gke_config", "gke_config_v2", "k3s_config", "oke_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: clusterRKE2ConfigFields(),
			},
		},
		"k3s_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config", "eks_config_v2", "gke_config", "gke_config_v2", "oke_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterK3SConfigFields(),
			},
		},
		"eks_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config_v2", "gke_config", "gke_config_v2", "k3s_config", "oke_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterEKSConfigFields(),
			},
		},
		"eks_config_v2": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config", "gke_config", "gke_config_v2", "k3s_config", "oke_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterEKSConfigV2Fields(),
			},
		},
		"aks_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"eks_config", "aks_config_v2", "eks_config_v2", "gke_config", "gke_config_v2", "k3s_config", "oke_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterAKSConfigFields(),
			},
		},
		"aks_config_v2": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "eks_config_v2", "gke_config", "gke_config_v2", "k3s_config", "oke_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterAKSConfigV2Fields(),
			},
		},
		"gke_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config", "eks_config_v2", "gke_config_v2", "k3s_config", "oke_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterGKEConfigFields(),
			},
		},
		"gke_config_v2": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config", "eks_config_v2", "gke_config", "k3s_config", "oke_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterGKEConfigV2Fields(),
			},
		},
		"oke_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "aks_config_v2", "eks_config", "eks_config_v2", "gke_config", "gke_config_v2", "k3s_config", "rke_config", "rke2_config"},
			Elem: &schema.Resource{
				Schema: clusterOKEConfigFields(),
			},
		},
		"default_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cluster_auth_endpoint": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterAuthEndpoint(),
			},
		},
		"cluster_monitoring_input": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster monitoring configuration",
			Elem: &schema.Resource{
				Schema: monitoringInputFields(),
			},
		},
		"cluster_registration_token": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRegistationTokenFields(),
			},
		},
		"cluster_template_answers": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Computed:    true,
			Description: "Cluster template answers",
			Elem: &schema.Resource{
				Schema: answerFields(),
			},
		},
		"cluster_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster template ID",
		},
		"cluster_template_questions": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster template questions",
			Elem: &schema.Resource{
				Schema: questionFields(),
			},
		},
		"cluster_template_revision_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster template revision ID",
		},
		"default_pod_security_policy_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default pod security policy template id",
		},
		"desired_agent_image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"desired_auth_image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"docker_root_dir": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"enable_cluster_alerting": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable built-in cluster alerting",
		},
		"enable_cluster_monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable built-in cluster monitoring",
		},
		"enable_cluster_istio": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable built-in cluster istio",
			Deprecated:  "Deploy istio using rancher2_app resource instead",
		},
		"fleet_workspace_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"istio_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Is istio enabled at cluster?",
		},
		"enable_network_policy": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable project network isolation",
		},
		"windows_prefered_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Windows preferred cluster",
			ForceNew:    true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
