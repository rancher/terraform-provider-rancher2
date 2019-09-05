package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenRollingUpdate(in *managementClient.RollingUpdate) []interface{} {
	obj := make(map[string]interface{})

	if in == nil {
		return []interface{}{}
	}

	if in.BatchSize > 0 {
		obj["batch_size"] = int(in.BatchSize)
	}

	if in.Interval > 0 {
		obj["interval"] = int(in.Interval)
	}

	return []interface{}{obj}
}

func flattenUpgradeStrategy(in *managementClient.UpgradeStrategy) []interface{} {
	obj := make(map[string]interface{})

	if in == nil {
		return []interface{}{}
	}

	if in.RollingUpdate != nil {
		obj["rolling_update"] = flattenRollingUpdate(in.RollingUpdate)
	}

	return []interface{}{obj}
}

// Expanders

func expandRollingUpdate(p []interface{}) *managementClient.RollingUpdate {
	obj := &managementClient.RollingUpdate{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["batch_size"].(int); ok && v > 0 {
		obj.BatchSize = int64(v)
	}

	if v, ok := in["interval"].(int); ok && v > 0 {
		obj.Interval = int64(v)
	}

	return obj
}

func expandUpgradeStrategy(p []interface{}) *managementClient.UpgradeStrategy {
	obj := &managementClient.UpgradeStrategy{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["rolling_update"].([]interface{}); ok && len(v) > 0 {
		obj.RollingUpdate = expandRollingUpdate(v)
	}

	return obj
}
