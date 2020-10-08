package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterK3SUpgradeStrategyConfig(in *managementClient.ClusterUpgradeStrategy) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["drain_server_nodes"] = in.DrainServerNodes
	obj["drain_worker_nodes"] = in.DrainWorkerNodes
	obj["server_concurrency"] = int(in.ServerConcurrency)
	obj["worker_concurrency"] = int(in.WorkerConcurrency)

	return []interface{}{obj}
}

func flattenClusterK3SConfig(in *managementClient.K3sConfig) []interface{} {
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

func expandClusterK3SUpgradeStrategyConfig(p []interface{}) *managementClient.ClusterUpgradeStrategy {
	obj := &managementClient.ClusterUpgradeStrategy{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["drain_server_nodes"].(bool); ok {
		obj.DrainServerNodes = v
	}

	if v, ok := in["drain_worker_nodes"].(bool); ok {
		obj.DrainWorkerNodes = v
	}

	if v, ok := in["server_concurrency"].(int); ok && v > 0 {
		obj.ServerConcurrency = int64(v)
	}

	if v, ok := in["worker_concurrency"].(int); ok && v > 0 {
		obj.WorkerConcurrency = int64(v)
	}

	return obj
}

func expandClusterK3SConfig(p []interface{}) *managementClient.K3sConfig {
	obj := &managementClient.K3sConfig{}
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
