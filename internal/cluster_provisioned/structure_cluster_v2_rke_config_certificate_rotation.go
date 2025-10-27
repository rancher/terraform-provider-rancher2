package rancher2

import rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"

func flattenClusterV2RKEConfigRotateCertificates(in *rkev1.RotateCertificates) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if in.Generation > 0 {
		obj["generation"] = int(in.Generation)
	}

	if len(in.Services) > 0 {
		obj["services"] = toArrayInterfaceSorted(in.Services)
	}

	return []interface{}{obj}
}

func expandClusterV2RKEConfigRotateCertificates(p []interface{}) *rkev1.RotateCertificates {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &rkev1.RotateCertificates{}

	in := p[0].(map[string]interface{})

	if v, ok := in["generation"].(int); ok && v > 0 {
		obj.Generation = int64(v)
	}

	if v, ok := in["services"].([]interface{}); ok && len(v) > 0 {
		obj.Services = toArrayStringSorted(v)
	}

	return obj
}
