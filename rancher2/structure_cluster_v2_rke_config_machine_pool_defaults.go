package rancher2

import (
	"reflect"

	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
)

func flattenClusterV2RKEConfigMachinePoolDefaults(in provisionv1.RKEMachinePoolDefaults) []any {
	if reflect.ValueOf(in).IsZero() {
		return nil
	}

	obj := map[string]any{}

	if in.HostnameLengthLimit > 0 {
		obj["hostname_length_limit"] = in.HostnameLengthLimit
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterV2RKEConfigMachinePoolDefaults(d []any) provisionv1.RKEMachinePoolDefaults {
	if d == nil || len(d) == 0 || d[0] == nil {
		return provisionv1.RKEMachinePoolDefaults{}
	}

	obj := provisionv1.RKEMachinePoolDefaults{}

	in := d[0].(map[string]interface{})

	if v, ok := in["hostname_length_limit"].(int); ok {
		obj.HostnameLengthLimit = v
	}

	return obj
}
