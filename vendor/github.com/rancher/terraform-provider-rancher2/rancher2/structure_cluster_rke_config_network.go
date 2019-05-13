package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterRKEConfigNetworkCalico(in *managementClient.CalicoNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.CloudProvider) > 0 {
		obj["cloud_provider"] = in.CloudProvider
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetworkCanal(in *managementClient.CanalNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Iface) > 0 {
		obj["iface"] = in.Iface
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetworkFlannel(in *managementClient.FlannelNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Iface) > 0 {
		obj["iface"] = in.Iface
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetworkWeave(in *managementClient.WeaveNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetwork(in *managementClient.NetworkConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.CalicoNetworkProvider != nil {
		calicoNetwork, err := flattenClusterRKEConfigNetworkCalico(in.CalicoNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["calico_network_provider"] = calicoNetwork
	}

	if in.CanalNetworkProvider != nil {
		canalNetwork, err := flattenClusterRKEConfigNetworkCanal(in.CanalNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["canal_network_provider"] = canalNetwork
	}

	if in.FlannelNetworkProvider != nil {
		flannelNetwork, err := flattenClusterRKEConfigNetworkFlannel(in.FlannelNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["flannel_network_provider"] = flannelNetwork
	}

	if in.WeaveNetworkProvider != nil {
		weaveNetwork, err := flattenClusterRKEConfigNetworkWeave(in.WeaveNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["weave_network_provider"] = weaveNetwork
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Plugin) > 0 {
		obj["plugin"] = in.Plugin
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigNetworkCalico(p []interface{}) (*managementClient.CalicoNetworkProvider, error) {
	obj := &managementClient.CalicoNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cloud_provider"].(string); ok && len(v) > 0 {
		obj.CloudProvider = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetworkCanal(p []interface{}) (*managementClient.CanalNetworkProvider, error) {
	obj := &managementClient.CanalNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["iface"].(string); ok && len(v) > 0 {
		obj.Iface = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetworkFlannel(p []interface{}) (*managementClient.FlannelNetworkProvider, error) {
	obj := &managementClient.FlannelNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["iface"].(string); ok && len(v) > 0 {
		obj.Iface = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetworkWeave(p []interface{}) (*managementClient.WeaveNetworkProvider, error) {
	obj := &managementClient.WeaveNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetwork(p []interface{}) (*managementClient.NetworkConfig, error) {
	obj := &managementClient.NetworkConfig{}
	if len(p) == 0 || p[0] == nil {
		obj.Plugin = networkPluginDefault
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["calico_network_provider"].([]interface{}); ok && len(v) > 0 {
		calicoNetwork, err := expandClusterRKEConfigNetworkCalico(v)
		if err != nil {
			return obj, err
		}
		obj.CalicoNetworkProvider = calicoNetwork
	}

	if v, ok := in["canal_network_provider"].([]interface{}); ok && len(v) > 0 {
		canalNetwork, err := expandClusterRKEConfigNetworkCanal(v)
		if err != nil {
			return obj, err
		}
		obj.CanalNetworkProvider = canalNetwork
	}

	if v, ok := in["flannel_network_provider"].([]interface{}); ok && len(v) > 0 {
		flannelNetwork, err := expandClusterRKEConfigNetworkFlannel(v)
		if err != nil {
			return obj, err
		}
		obj.FlannelNetworkProvider = flannelNetwork
	}

	if v, ok := in["weave_network_provider"].([]interface{}); ok && len(v) > 0 {
		weaveNetwork, err := expandClusterRKEConfigNetworkWeave(v)
		if err != nil {
			return obj, err
		}
		obj.WeaveNetworkProvider = weaveNetwork
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["plugin"].(string); ok && len(v) > 0 {
		obj.Plugin = v
	}

	return obj, nil
}
