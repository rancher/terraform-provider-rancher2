package rancher2

import rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"

func flattenClusterV2RKEConfigRotateEncryptionKeys(in *rkev1.RotateEncryptionKeys) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if in.Generation > 0 {
		obj["generation"] = int(in.Generation)
	}

	return []interface{}{obj}
}

func expandClusterV2RKEConfigRotateEncryptionKeys(p []interface{}) *rkev1.RotateEncryptionKeys {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &rkev1.RotateEncryptionKeys{}

	in := p[0].(map[string]interface{})

	if v, ok := in["generation"].(int); ok && v > 0 {
		obj.Generation = int64(v)
	}

	return obj
}
