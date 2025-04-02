package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func clusterV2FieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Cluster V2 name",
		},
		"fleet_namespace": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "fleet-default",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Cluster V2 kubernetes version",
		},
		"local_auth_endpoint": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster V2 local auth endpoint",
			Elem: &schema.Resource{
				Schema: clusterV2LocalAuthEndpointFields(),
			},
		},
		"rke_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 rke config",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigFieldsV0(),
			},
		},
		"agent_env_vars": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster V2 default agent env vars",
			Elem: &schema.Resource{
				Schema: envVarFields(),
			},
		},
		"cloud_credential_secret_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 cloud credential secret name",
		},
		"cluster_agent_deployment_customization": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional customization for cluster agent",
			Elem: &schema.Resource{
				Schema: agentDeploymentCustomizationFields(true),
			},
		},
		"fleet_agent_deployment_customization": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional customization for fleet agent",
			Elem: &schema.Resource{
				Schema: agentDeploymentCustomizationFields(false),
			},
		},
		"default_pod_security_admission_configuration_template_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 default pod security admission configuration template name",
		},
		"default_cluster_role_for_project_members": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 default cluster role for project members",
		},
		"enable_network_policy": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable k8s network policy",
		},
		// Computed attributes
		"cluster_registration_token": {
			Type:      schema.TypeList,
			MaxItems:  1,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Resource{
				Schema: clusterRegistrationTokenFields(),
			},
		},
		"kube_config": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"cluster_v1_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"resource_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func clusterV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Cluster V2 name",
		},
		"fleet_namespace": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "fleet-default",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Cluster V2 kubernetes version",
		},
		"local_auth_endpoint": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster V2 local auth endpoint",
			Elem: &schema.Resource{
				Schema: clusterV2LocalAuthEndpointFields(),
			},
		},
		"rke_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 rke config",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigFields(),
			},
		},
		"agent_env_vars": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster V2 default agent env vars",
			Elem: &schema.Resource{
				Schema: envVarFields(),
			},
		},
		"cloud_credential_secret_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 cloud credential secret name",
		},
		"cluster_agent_deployment_customization": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional customization for cluster agent",
			Elem: &schema.Resource{
				Schema: agentDeploymentCustomizationFields(true),
			},
		},
		"fleet_agent_deployment_customization": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional customization for fleet agent",
			Elem: &schema.Resource{
				Schema: agentDeploymentCustomizationFields(false),
			},
		},
		"default_pod_security_admission_configuration_template_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 default pod security admission configuration template name",
		},
		"default_cluster_role_for_project_members": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 default cluster role for project members",
		},
		"enable_network_policy": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable k8s network policy",
		},
		// Computed attributes
		"cluster_registration_token": {
			Type:      schema.TypeList,
			MaxItems:  1,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Resource{
				Schema: clusterRegistrationTokenFields(),
			},
		},
		"kube_config": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"cluster_v1_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"resource_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
