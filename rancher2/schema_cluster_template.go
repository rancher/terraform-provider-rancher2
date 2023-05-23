package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	questionTypeInt      = "int"
	questionTypeBool     = "boolean"
	questionTypeString   = "string"
	questionTypePassword = "password"
)

var (
	questionTypeKinds = []string{questionTypeInt, questionTypeBool, questionTypeString, questionTypePassword}
)

//Schemas

func questionFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Default variable value",
		},
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Required variable",
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      questionTypeString,
			Description:  "Variable type",
			ValidateFunc: validation.StringInSlice(questionTypeKinds, true),
		},
		"variable": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Variable name",
		},
	}

	return s
}

func clusterSpecBaseFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_auth_endpoint": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Local cluster auth endpoint",
			Elem: &schema.Resource{
				Schema: clusterAuthEndpoint(),
			},
		},
		"default_cluster_role_for_project_members": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default cluster role for project members",
		},
		"default_pod_security_policy_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default pod security policy template ID",
		},
		"desired_agent_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Desired agent image",
		},
		"desired_auth_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Desired auth image",
		},
		"docker_root_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Docker Root Dir",
		},
		"enable_cluster_alerting": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable built-in cluster alerting",
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
		"rke_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Rancher Kubernetes Engine Config",
			Elem: &schema.Resource{
				Schema: clusterRKEConfigFieldsV0(),
			},
		},
		"windows_prefered_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Windows prefered cluster",
		},
	}

	return s
}

func clusterSpecBaseFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_auth_endpoint": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Local cluster auth endpoint",
			Elem: &schema.Resource{
				Schema: clusterAuthEndpoint(),
			},
		},
		"default_cluster_role_for_project_members": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default cluster role for project members",
		},
		"default_pod_security_policy_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default pod security policy template ID",
		},
		"desired_agent_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Desired agent image",
		},
		"desired_auth_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Desired auth image",
		},
		"docker_root_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Docker Root Dir",
		},
		"enable_cluster_alerting": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable built-in cluster alerting",
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
		"rke_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Rancher Kubernetes Engine Config",
			Elem: &schema.Resource{
				Schema: clusterRKEConfigFields(),
			},
		},
		"windows_prefered_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Windows prefered cluster",
		},
	}

	return s
}

// Used by datasource
func clusterSpecBaseFieldsData() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_auth_endpoint": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Local cluster auth endpoint",
			Elem: &schema.Resource{
				Schema: clusterAuthEndpoint(),
			},
		},
		"default_cluster_role_for_project_members": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default cluster role for project members",
		},
		"default_pod_security_policy_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Default pod security policy template ID",
		},
		"desired_agent_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Desired agent image",
		},
		"desired_auth_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Desired auth image",
		},
		"docker_root_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Docker Root Dir",
		},
		"enable_cluster_alerting": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable built-in cluster alerting",
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
		"rke_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Rancher Kubernetes Engine Config",
			Elem: &schema.Resource{
				Schema: clusterRKEConfigFieldsData(),
			},
		},
		"windows_prefered_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Windows prefered cluster",
		},
	}

	return s
}

func clusterTemplateRevisionFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Cluster configuration",
			Elem: &schema.Resource{
				Schema: clusterSpecBaseFieldsV0(),
			},
		},
		"cluster_template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template ID",
		},
		"default": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Default cluster template revision",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable cluster template revision",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template revision ID",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Cluster template revision name",
		},
		"questions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster template questions",
			Elem: &schema.Resource{
				Schema: questionFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func clusterTemplateRevisionFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Cluster configuration",
			Elem: &schema.Resource{
				Schema: clusterSpecBaseFields(),
			},
		},
		"cluster_template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template ID",
		},
		"default": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Default cluster template revision",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable cluster template revision",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template revision ID",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Cluster template revision name",
		},
		"questions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster template questions",
			Elem: &schema.Resource{
				Schema: questionFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

// Used by datasource
func clusterTemplateRevisionFieldsData() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Cluster configuration",
			Elem: &schema.Resource{
				Schema: clusterSpecBaseFieldsData(),
			},
		},
		"cluster_template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template ID",
		},
		"default": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Default cluster template revision",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable cluster template revision",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster template revision ID",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Cluster template revision name",
		},
		"questions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster template questions",
			Elem: &schema.Resource{
				Schema: questionFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func clusterTemplateFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default_revision_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Default cluster template revision ID",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster template description",
		},
		"members": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster template members",
			Elem: &schema.Resource{
				Schema: memberFields(),
			},
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Cluster template name",
		},
		"template_revisions": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster template revisions",
			Elem: &schema.Resource{
				Schema: clusterTemplateRevisionFieldsV0(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func clusterTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default_revision_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Default cluster template revision ID",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster template description",
		},
		"members": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster template members",
			Elem: &schema.Resource{
				Schema: memberFields(),
			},
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Cluster template name",
		},
		"template_revisions": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster template revisions",
			Elem: &schema.Resource{
				Schema: clusterTemplateRevisionFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
