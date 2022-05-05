package rancher2

import rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"

func flattenClusterV2RKEConfigETCDSnapshotRestore(in *rkev1.ETCDSnapshotRestore) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.Name) > 0 {
		obj["name"] = in.Name
	}
	if in.Generation > 0 {
		obj["generation"] = in.Generation
	}
	if len(in.RestoreRKEConfig) > 0 {
		obj["restore_rke_config"] = in.RestoreRKEConfig
	}

	return []interface{}{obj}
}

func expandClusterV2RKEConfigETCDSnapshotRestore(p []interface{}) *rkev1.ETCDSnapshotRestore {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &rkev1.ETCDSnapshotRestore{}

	in := p[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	if v, ok := in["generation"].(int); ok && v > 0 {
		obj.Generation = v
	}
	if v, ok := in["restore_rke_config"].(string); ok && len(v) > 0 {
		obj.RestoreRKEConfig = v
	}

	return obj
}
