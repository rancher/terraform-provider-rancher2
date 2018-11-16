package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	networkPluginFlannelName = "flannel"
)

//Schemas

func flannelNetworkProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"iface": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenFlannelNetworkProvider(in *managementClient.FlannelNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Iface) > 0 {
		obj["iface"] = in.Iface
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandFlannelNetworkProvider(p []interface{}) (*managementClient.FlannelNetworkProvider, error) {
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
