package rancher2

import (
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

// Flatteners

func flattenClusterV2Networking(in *rkev1.Networking) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.StackPreference) > 0 {
		obj["stack_preference"] = string(in.StackPreference)
	}
	return []interface{}{obj}
}

// Expanders

func expandClusterV2Networking(p []interface{}) *rkev1.Networking {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &rkev1.Networking{}

	in := p[0].(map[string]interface{})

	if v, ok := in["stack_preference"].(string); ok && len(v) > 0 {
		obj.StackPreference = rkev1.NetworkingStackPreference(v)
	}

	return obj
}
