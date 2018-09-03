package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func rkeConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"addon_job_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Optional duration in seconds of addon job.",
		},
		"authentication": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster authentication",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: authenticationFields(),
			},
		},
		"cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: cloudProviderFields(),
			},
		},
		"ignore_docker_version": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Optional ignore docker version on nodes",
		},
		"ingress": {
			Type:        schema.TypeList,
			Description: "Kubernetes ingress configuration",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: ingressFields(),
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
				Schema: monitoringFields(),
			},
		},
		"network": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster networking",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: networkFields(),
			},
		},
		"services": {
			Type:        schema.TypeList,
			Description: "Kubernetes cluster services",
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: servicesFields(),
			},
		},
		"ssh_agent_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Optional use ssh agent auth",
		},
	}

	return s
}

// Flatteners

func flattenRkeConfig(in *managementClient.RancherKubernetesEngineConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.AddonJobTimeout > 0 {
		obj["addon_job_timeout"] = int(in.AddonJobTimeout)
	}

	if in.Authentication != nil {
		authn, err := flattenAuthentication(in.Authentication)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["authentication"] = authn
	}

	if in.CloudProvider != nil {
		cloudProvider, err := flattenCloudProvider(in.CloudProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["cloud_provider"] = cloudProvider
	}

	obj["ignore_docker_version"] = in.IgnoreDockerVersion

	if in.Ingress != nil {
		ingress, err := flattenIngress(in.Ingress)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["ingress"] = ingress
	}

	if len(in.Version) > 0 {
		obj["kubernetes_version"] = in.Version
	}

	if in.Monitoring != nil {
		monitoring, err := flattenMonitoring(in.Monitoring)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["monitoring"] = monitoring
	}

	if in.Network != nil {
		network, err := flattenNetwork(in.Network)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["network"] = network
	}

	if in.Services != nil {
		services, err := flattenServices(in.Services)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["services"] = services
	}

	obj["ssh_agent_auth"] = in.SSHAgentAuth

	return []interface{}{obj}, nil
}

// Expanders

func expandRkeConfig(p []interface{}) (*managementClient.RancherKubernetesEngineConfig, error) {
	obj := &managementClient.RancherKubernetesEngineConfig{}

	// Set default network
	network, err := expandNetwork([]interface{}{})
	if err != nil {
		return obj, err
	}
	obj.Network = network

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["addon_job_timeout"].(int); ok && v > 0 {
		obj.AddonJobTimeout = int64(v)
	}

	if v, ok := in["authentication"].([]interface{}); ok && len(v) > 0 {
		authn, err := expandAuthentication(v)
		if err != nil {
			return obj, err
		}
		obj.Authentication = authn
	}

	if v, ok := in["cloud_provider"].([]interface{}); ok && len(v) > 0 {
		cloudProvider, err := expandCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.CloudProvider = cloudProvider
	}

	if v, ok := in["ignore_docker_version"].(bool); ok {
		obj.IgnoreDockerVersion = v
	}

	if v, ok := in["ingress"].([]interface{}); ok && len(v) > 0 {
		ingress, err := expandIngress(v)
		if err != nil {
			return obj, err
		}
		obj.Ingress = ingress
	}

	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		obj.Version = v
	}

	if v, ok := in["monitoring"].([]interface{}); ok && len(v) > 0 {
		monitoring, err := expandMonitoring(v)
		if err != nil {
			return obj, err
		}
		obj.Monitoring = monitoring
	}

	if v, ok := in["network"].([]interface{}); ok && len(v) > 0 {
		network, err := expandNetwork(v)
		if err != nil {
			return obj, err
		}
		obj.Network = network
	}

	if v, ok := in["services"].([]interface{}); ok && len(v) > 0 {
		services, err := expandServices(v)
		if err != nil {
			return obj, err
		}
		obj.Services = services
	}

	if v, ok := in["ssh_agent_auth"].(bool); ok {
		obj.SSHAgentAuth = v
	}

	return obj, nil
}
