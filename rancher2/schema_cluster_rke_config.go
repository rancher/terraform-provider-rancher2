package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterRKEKind   = "rke"
	clusterDriverRKE = "rancherKubernetesEngine"
)

//Schemas

func clusterRKEConfigFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"addon_job_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Optional duration in seconds of addon job.",
		},
		"addons": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional addons descripton to deploy on rke cluster.",
			StateFunc:   TrimSpace,
		},
		"addons_include": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional addons yaml manisfest to deploy on rke cluster.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"authentication": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster authentication",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigAuthenticationFields(),
			},
		},
		"authorization": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster authorization",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigAuthorizationFields(),
			},
		},
		"bastion_host": {
			Type:        schema.TypeList,
			Description: "RKE bastion host",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigBastionHostFields(),
			},
		},
		"cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderFields(),
			},
		},
		"dns": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigDNSFields(),
			},
		},
		"ignore_docker_version": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Optional ignore docker version on nodes",
		},
		"ingress": {
			Type:        schema.TypeList,
			Description: "Kubernetes ingress configuration",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigIngressFields(),
			},
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional kubernetes version to deploy",
		},
		"monitoring": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster monitoring",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigMonitoringFields(),
			},
		},
		"network": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster networking",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNetworkFields(),
			},
		},
		"nodes": {
			Type:        schema.TypeList,
			Description: "Optional RKE cluster nodes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNodesFields(),
			},
		},
		"prefix_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional prefix to customize kubernetes path",
		},
		"private_registries": {
			Type:        schema.TypeList,
			Description: "Optional private registries for docker images",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigPrivateRegistriesFields(),
			},
		},
		"services": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster services",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesFieldsV0(),
			},
		},
		"ssh_agent_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Optional use ssh agent auth",
		},
		"ssh_cert_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional cluster level SSH certificate path",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional cluster level SSH private key path",
		},
		"upgrade_strategy": {
			Type:        schema.TypeList,
			Description: "RKE upgrade strategy",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNodeUpgradeStrategyFields(),
			},
		},
	}

	return s
}

func clusterRKEConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"addon_job_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Optional duration in seconds of addon job.",
		},
		"addons": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional addons descripton to deploy on rke cluster.",
			StateFunc:   TrimSpace,
		},
		"addons_include": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional addons yaml manisfest to deploy on rke cluster.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"authentication": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster authentication",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigAuthenticationFields(),
			},
		},
		"authorization": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster authorization",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigAuthorizationFields(),
			},
		},
		"bastion_host": {
			Type:        schema.TypeList,
			Description: "RKE bastion host",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigBastionHostFields(),
			},
		},
		"cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderFields(),
			},
		},
		"dns": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigDNSFields(),
			},
		},
		"enable_cri_dockerd": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable/disable using cri-dockerd",
		},
		"ignore_docker_version": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Optional ignore docker version on nodes",
		},
		"ingress": {
			Type:        schema.TypeList,
			Description: "Kubernetes ingress configuration",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigIngressFields(),
			},
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional kubernetes version to deploy",
		},
		"monitoring": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster monitoring",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigMonitoringFields(),
			},
		},
		"network": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster networking",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNetworkFields(),
			},
		},
		"nodes": {
			Type:        schema.TypeList,
			Description: "Optional RKE cluster nodes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNodesFields(),
			},
		},
		"prefix_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional prefix to customize kubernetes path",
		},
		"win_prefix_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional prefix to customize kubernetes path for windows",
		},
		"private_registries": {
			Type:        schema.TypeList,
			Description: "Optional private registries for docker images",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigPrivateRegistriesFields(),
			},
		},
		"services": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster services",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesFields(),
			},
		},
		"ssh_agent_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Optional use ssh agent auth",
		},
		"ssh_cert_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional cluster level SSH certificate path",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional cluster level SSH private key path",
		},
		"upgrade_strategy": {
			Type:        schema.TypeList,
			Description: "RKE upgrade strategy",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNodeUpgradeStrategyFields(),
			},
		},
	}

	return s
}

// Used by datasource
func clusterRKEConfigFieldsData() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"addon_job_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Optional duration in seconds of addon job.",
		},
		"addons": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional addons descripton to deploy on rke cluster.",
			StateFunc:   TrimSpace,
		},
		"addons_include": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Optional addons yaml manisfest to deploy on rke cluster.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"authentication": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster authentication",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigAuthenticationFields(),
			},
		},
		"authorization": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster authorization",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigAuthorizationFields(),
			},
		},
		"bastion_host": {
			Type:        schema.TypeList,
			Description: "RKE bastion host",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigBastionHostFields(),
			},
		},
		"cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderFields(),
			},
		},
		"dns": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigDNSFields(),
			},
		},
		"enable_cri_dockerd": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable/disable using cri-dockerd",
		},
		"ignore_docker_version": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Optional ignore docker version on nodes",
		},
		"ingress": {
			Type:        schema.TypeList,
			Description: "Kubernetes ingress configuration",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigIngressFields(),
			},
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional kubernetes version to deploy",
		},
		"monitoring": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster monitoring",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigMonitoringFields(),
			},
		},
		"network": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster networking",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNetworkFields(),
			},
		},
		"nodes": {
			Type:        schema.TypeList,
			Description: "Optional RKE cluster nodes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNodesFields(),
			},
		},
		"prefix_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional prefix to customize kubernetes path",
		},
		"win_prefix_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional prefix to customize kubernetes path for windows nodes",
		},
		"private_registries": {
			Type:        schema.TypeList,
			Description: "Optional private registries for docker images",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigPrivateRegistriesFields(),
			},
		},
		"services": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster services",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesFieldsData(),
			},
		},
		"ssh_agent_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Optional use ssh agent auth",
		},
		"ssh_cert_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional cluster level SSH certificate path",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional cluster level SSH private key path",
		},
		"upgrade_strategy": {
			Type:        schema.TypeList,
			Description: "RKE upgrade strategy",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNodeUpgradeStrategyFields(),
			},
		},
	}

	return s
}
