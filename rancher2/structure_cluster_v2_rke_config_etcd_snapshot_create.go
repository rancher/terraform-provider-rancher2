package rancher2

import rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"

func flattenClusterV2RKEConfigETCDSnapshotCreate(in *rkev1.ETCDSnapshotCreate) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if in.Generation > 0 {
		obj["generation"] = in.Generation
	}

	return []interface{}{obj}
}

func expandClusterV2RKEConfigETCDSnapshotCreate(p []interface{}) *rkev1.ETCDSnapshotCreate {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &rkev1.ETCDSnapshotCreate{}

	in := p[0].(map[string]interface{})

	if v, ok := in["generation"].(int); ok && v > 0 {
		obj.Generation = v
	}

	return obj
}
