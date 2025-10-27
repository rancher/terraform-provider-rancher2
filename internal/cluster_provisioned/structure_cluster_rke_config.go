package rancher2

import (
	"fmt"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfig(in *managementClient.RancherKubernetesEngineConfig, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.AddonJobTimeout > 0 {
		obj["addon_job_timeout"] = int(in.AddonJobTimeout)
	}

	if len(in.Addons) > 0 {
		obj["addons"] = in.Addons
	}

	if len(in.AddonsInclude) > 0 {
		obj["addons_include"] = toArrayInterface(in.AddonsInclude)
	}

	if in.Authentication != nil {
		authn, err := flattenClusterRKEConfigAuthentication(in.Authentication)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["authentication"] = authn
	}

	if in.Authorization != nil {
		authz, err := flattenClusterRKEConfigAuthorization(in.Authorization)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["authorization"] = authz
	}

	if in.BastionHost != nil {
		v, ok := obj["bastion_host"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		bastion, err := flattenClusterRKEConfigBastionHost(in.BastionHost, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["bastion_host"] = bastion
	}

	if in.CloudProvider != nil {
		v, ok := obj["cloud_provider"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		cloudProvider, err := flattenClusterRKEConfigCloudProvider(in.CloudProvider, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["cloud_provider"] = cloudProvider
	}

	if in.DNS != nil {
		dns, err := flattenClusterRKEConfigDNS(in.DNS)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["dns"] = dns
	}

	if in.EnableCRIDockerd != nil {
		obj["enable_cri_dockerd"] = *in.EnableCRIDockerd
	}
	if in.IgnoreDockerVersion != nil {
		obj["ignore_docker_version"] = *in.IgnoreDockerVersion
	}

	if in.Ingress != nil {
		ingress, err := flattenClusterRKEConfigIngress(in.Ingress)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["ingress"] = ingress
	}

	if len(in.Version) > 0 {
		obj["kubernetes_version"] = in.Version
	}

	if in.Monitoring != nil {
		monitoring, err := flattenClusterRKEConfigMonitoring(in.Monitoring)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["monitoring"] = monitoring
	}

	if in.Network != nil {
		network, err := flattenClusterRKEConfigNetwork(in.Network)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["network"] = network
	}

	if in.Nodes != nil {
		nodes, err := flattenClusterRKEConfigNodes(in.Nodes)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["nodes"] = nodes
	}

	if len(in.PrefixPath) > 0 {
		obj["prefix_path"] = in.PrefixPath
	}

	if len(in.WindowsPrefixPath) > 0 {
		obj["win_prefix_path"] = in.WindowsPrefixPath
	}

	if in.PrivateRegistries != nil {
		v, ok := obj["private_registries"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		privReg, err := flattenClusterRKEConfigPrivateRegistries(in.PrivateRegistries, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["private_registries"] = privReg
	}

	if in.Services != nil {
		v, ok := obj["services"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		services, err := flattenClusterRKEConfigServices(in.Services, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["services"] = services
	}

	obj["ssh_agent_auth"] = in.SSHAgentAuth

	if len(in.SSHCertPath) > 0 {
		obj["ssh_cert_path"] = in.SSHCertPath
	}

	if len(in.SSHKeyPath) > 0 {
		obj["ssh_key_path"] = in.SSHKeyPath
	}

	if in.UpgradeStrategy != nil {
		obj["upgrade_strategy"] = flattenClusterRKEConfigNodeUpgradeStrategy(in.UpgradeStrategy)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfig(p []interface{}, name string) (*managementClient.RancherKubernetesEngineConfig, error) {
	obj := &managementClient.RancherKubernetesEngineConfig{}

	// Set default network
	network, err := expandClusterRKEConfigNetwork([]interface{}{})
	if err != nil {
		return obj, err
	}
	obj.Network = network

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	obj.ClusterName = name

	if v, ok := in["addon_job_timeout"].(int); ok && v > 0 {
		obj.AddonJobTimeout = int64(v)
	}

	if v, ok := in["addons"].(string); ok && len(v) > 0 {
		obj.Addons = v
	}

	if v, ok := in["addons_include"].([]interface{}); ok && len(v) > 0 {
		obj.AddonsInclude = toArrayString(v)
	}

	if v, ok := in["authentication"].([]interface{}); ok && len(v) > 0 {
		authn, err := expandClusterRKEConfigAuthentication(v)
		if err != nil {
			return obj, err
		}
		obj.Authentication = authn
	}

	if v, ok := in["authorization"].([]interface{}); ok && len(v) > 0 {
		authz, err := expandClusterRKEConfigAuthorization(v)
		if err != nil {
			return obj, err
		}
		obj.Authorization = authz
	}

	if v, ok := in["bastion_host"].([]interface{}); ok && len(v) > 0 {
		bastion, err := expandClusterRKEConfigBastionHost(v)
		if err != nil {
			return obj, err
		}
		obj.BastionHost = bastion
	}

	if v, ok := in["cloud_provider"].([]interface{}); ok && len(v) > 0 {
		cloudProvider, err := expandClusterRKEConfigCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.CloudProvider = cloudProvider
	}

	if v, ok := in["dns"].([]interface{}); ok && len(v) > 0 {
		dns, err := expandClusterRKEConfigDNS(v)
		if err != nil {
			return obj, err
		}
		obj.DNS = dns
	}

	if v, ok := in["enable_cri_dockerd"].(bool); ok {
		obj.EnableCRIDockerd = &v
	}
	if v, ok := in["ignore_docker_version"].(bool); ok {
		obj.IgnoreDockerVersion = &v
	}

	if v, ok := in["ingress"].([]interface{}); ok && len(v) > 0 {
		ingress, err := expandClusterRKEConfigIngress(v)
		if err != nil {
			return obj, err
		}
		obj.Ingress = ingress
	}

	if len(rancher2ClusterRKEK8SDefaultVersion) > 0 {
		obj.Version = rancher2ClusterRKEK8SDefaultVersion
	}
	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		obj.Version = v
		found := false
		if len(rancher2ClusterRKEK8SVersions) > 0 {
			for _, v := range rancher2ClusterRKEK8SVersions {
				if obj.Version == v {
					found = true
				}
			}
			if !found {
				return obj, fmt.Errorf("RKE version is not supported %s got %s", rancher2ClusterRKEK8SVersions, obj.Version)
			}
		}
	}

	if v, ok := in["monitoring"].([]interface{}); ok && len(v) > 0 {
		monitoring, err := expandClusterRKEConfigMonitoring(v)
		if err != nil {
			return obj, err
		}
		obj.Monitoring = monitoring
	}

	if v, ok := in["network"].([]interface{}); ok && len(v) > 0 {
		network, err := expandClusterRKEConfigNetwork(v)
		if err != nil {
			return obj, err
		}
		obj.Network = network
	}

	if v, ok := in["nodes"].([]interface{}); ok && len(v) > 0 {
		nodes, err := expandClusterRKEConfigNodes(v)
		if err != nil {
			return obj, err
		}
		obj.Nodes = nodes
	}

	if v, ok := in["prefix_path"].(string); ok && len(v) > 0 {
		obj.PrefixPath = v
	}
	if v, ok := in["win_prefix_path"].(string); ok && len(v) > 0 {
		obj.WindowsPrefixPath = v
	}

	if v, ok := in["private_registries"].([]interface{}); ok && len(v) > 0 {
		privReg, err := expandClusterRKEConfigPrivateRegistries(v)
		if err != nil {
			return obj, err
		}
		obj.PrivateRegistries = privReg
	}

	if v, ok := in["services"].([]interface{}); ok && len(v) > 0 {
		services, err := expandClusterRKEConfigServices(v)
		if err != nil {
			return obj, err
		}
		obj.Services = services
	}

	if v, ok := in["ssh_agent_auth"].(bool); ok {
		obj.SSHAgentAuth = v
	}

	if v, ok := in["ssh_cert_path"].(string); ok && len(v) > 0 {
		obj.SSHCertPath = v
	}

	if v, ok := in["ssh_key_path"].(string); ok && len(v) > 0 {
		obj.SSHKeyPath = v
	}

	if v, ok := in["upgrade_strategy"].([]interface{}); ok {
		obj.UpgradeStrategy = expandClusterRKEConfigNodeUpgradeStrategy(v)
	}

	return obj, nil
}
