package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	networkPluginCalicoName = "calico"
)

//Schemas

func calicoNetworkProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cloud_provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenCalicoNetworkProvider(in *managementClient.CalicoNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.CloudProvider) > 0 {
		obj["cloud_provider"] = in.CloudProvider
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandCalicoNetworkProvider(p []interface{}) (*managementClient.CalicoNetworkProvider, error) {
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
