package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	networkPluginDefault = networkPluginCanalName
	networkPluginList    = []string{networkPluginCanalName, networkPluginFlannelName, networkPluginCalicoName}
)

//Schemas

func networkFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"calico_network_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: calicoNetworkProviderFields(),
			},
		},
		"canal_network_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: canalNetworkProviderFields(),
			},
		},
		"flannel_network_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: flannelNetworkProviderFields(),
			},
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"plugin": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(networkPluginList, true),
		},
	}
	return s
}

// Flatteners

func flattenNetwork(in *managementClient.NetworkConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.CalicoNetworkProvider != nil {
		calicoNetwork, err := flattenCalicoNetworkProvider(in.CalicoNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["calico_network_provider"] = calicoNetwork
	}

	if in.CanalNetworkProvider != nil {
		canalNetwork, err := flattenCanalNetworkProvider(in.CanalNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["canal_network_provider"] = canalNetwork
	}

	if in.FlannelNetworkProvider != nil {
		flannelNetwork, err := flattenFlannelNetworkProvider(in.FlannelNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["flannel_network_provider"] = flannelNetwork
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

func expandNetwork(p []interface{}) (*managementClient.NetworkConfig, error) {
	obj := &managementClient.NetworkConfig{}
	if len(p) == 0 || p[0] == nil {
		obj.Plugin = networkPluginDefault
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["calico_network_provider"].([]interface{}); ok && len(v) > 0 {
		calicoNetwork, err := expandCalicoNetworkProvider(v)
		if err != nil {
			return obj, err
		}
		obj.CalicoNetworkProvider = calicoNetwork
	}

	if v, ok := in["canal_network_provider"].([]interface{}); ok && len(v) > 0 {
		canalNetwork, err := expandCanalNetworkProvider(v)
		if err != nil {
			return obj, err
		}
		obj.CanalNetworkProvider = canalNetwork
	}

	if v, ok := in["flannel_network_provider"].([]interface{}); ok && len(v) > 0 {
		flannelNetwork, err := expandFlannelNetworkProvider(v)
		if err != nil {
			return obj, err
		}
		obj.FlannelNetworkProvider = flannelNetwork
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["plugin"].(string); ok && len(v) > 0 {
		obj.Plugin = v
	}

	return obj, nil
}
