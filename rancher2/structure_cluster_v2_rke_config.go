package rancher2

import (
	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
)

// Flatteners

func flattenClusterV2RKEConfig(in *provisionv1.RKEConfig) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.AdditionalManifest) > 0 {
		obj["additional_manifest"] = in.AdditionalManifest
	}

	obj["local_auth_endpoint"] = flattenClusterV2RKEConfigLocalAuthEndpoint(in.LocalClusterAuthEndpoint)
	obj["upgrade_strategy"] = flattenClusterV2RKEConfigUpgradeStrategy(in.UpgradeStrategy)

	if in.ChartValues.Data != nil && len(in.ChartValues.Data) > 0 {
		yamlData, _ := interfaceToGhodssyaml(in.ChartValues.Data)
		obj["chart_values"] = yamlData
	}
	if in.MachineGlobalConfig.Data != nil && len(in.MachineGlobalConfig.Data) > 0 {
		yamlData, _ := interfaceToGhodssyaml(in.MachineGlobalConfig.Data)
		obj["machine_global_config"] = yamlData
	}
	if in.MachinePools != nil && len(in.MachinePools) > 0 {
		obj["machine_pools"] = flattenClusterV2RKEConfigMachinePools(in.MachinePools)
	}
	if in.MachineSelectorConfig != nil && len(in.MachineSelectorConfig) > 0 {
		obj["machine_selector_config"] = flattenClusterV2RKEConfigSystemConfig(in.MachineSelectorConfig)
	}
	if in.Registries != nil {
		obj["registries"] = flattenClusterV2RKEConfigRegistry(in.Registries)
	}
	if in.ETCD != nil {
		obj["etcd"] = flattenClusterV2RKEConfigETCD(in.ETCD)
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterV2RKEConfig(p []interface{}) *provisionv1.RKEConfig {
	obj := &provisionv1.RKEConfig{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["additional_manifest"].(string); ok && len(v) > 0 {
		obj.AdditionalManifest = v
	}

	if v, ok := in["local_auth_endpoint"].([]interface{}); ok && len(v) > 0 {
		obj.LocalClusterAuthEndpoint = expandClusterV2RKEConfigLocalAuthEndpoint(v)
	}
	if v, ok := in["upgrade_strategy"].([]interface{}); ok && len(v) > 0 {
		obj.UpgradeStrategy = expandClusterV2RKEConfigUpgradeStrategy(v)
	}

	if v, ok := in["chart_values"].(string); ok && len(v) > 0 {
		values, _ := ghodssyamlToMapInterface(v)
		obj.ChartValues.Data = values
	}
	if v, ok := in["machine_global_config"].(string); ok && len(v) > 0 {
		values, _ := ghodssyamlToMapInterface(v)
		obj.MachineGlobalConfig.Data = values
	}
	if v, ok := in["machine_pools"].([]interface{}); ok && len(v) > 0 {
		obj.MachinePools = expandClusterV2RKEConfigMachinePools(v)
	}
	if v, ok := in["machine_selector_config"].([]interface{}); ok && len(v) > 0 {
		obj.MachineSelectorConfig = expandClusterV2RKEConfigSystemConfig(v)
	}
	if v, ok := in["registries"].([]interface{}); ok && len(v) > 0 {
		obj.Registries = expandClusterV2RKEConfigRegistry(v)
	}
	if v, ok := in["etcd"].([]interface{}); ok && len(v) > 0 {
		obj.ETCD = expandClusterV2RKEConfigETCD(v)
	}

	return obj
}
