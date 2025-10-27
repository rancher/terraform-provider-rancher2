package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKE2Config(in *managementClient.Rke2Config) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if in.ClusterUpgradeStrategy != nil {
		obj["upgrade_strategy"] = flattenClusterK3SUpgradeStrategyConfig(in.ClusterUpgradeStrategy)
	}

	if len(in.Version) > 0 {
		obj["version"] = in.Version
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterRKE2Config(p []interface{}) *managementClient.Rke2Config {
	obj := &managementClient.Rke2Config{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["upgrade_strategy"].([]interface{}); ok && len(v) > 0 {
		obj.ClusterUpgradeStrategy = expandClusterK3SUpgradeStrategyConfig(v)
	}

	if v, ok := in["version"].(string); ok && len(v) > 0 {
		obj.Version = v
	}

	return obj
}
