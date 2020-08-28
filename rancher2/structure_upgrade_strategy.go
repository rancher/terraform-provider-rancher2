package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"k8s.io/apimachinery/pkg/util/intstr"
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

func flattenRollingUpdateDaemonSet(in *managementClient.RollingUpdateDaemonSet) []interface{} {
	obj := make(map[string]interface{})

	if in == nil {
		return []interface{}{}
	}

	if v := in.MaxUnavailable.IntValue(); v > 0 {
		obj["max_unavailable"] = v
	}

	return []interface{}{obj}
}

func flattenRollingUpdateDeployment(in *managementClient.RollingUpdateDeployment) []interface{} {
	obj := make(map[string]interface{})

	if in == nil {
		return []interface{}{}
	}

	if v := in.MaxSurge.IntValue(); v > 0 {
		obj["max_surge"] = in.MaxSurge.IntValue()
	}

	if v := in.MaxUnavailable.IntValue(); v > 0 {
		obj["max_unavailable"] = in.MaxUnavailable.IntValue()
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

func flattenDaemonSetStrategy(in *managementClient.DaemonSetUpdateStrategy) []interface{} {
	obj := make(map[string]interface{})

	if in == nil {
		return []interface{}{}
	}

	if in.RollingUpdate != nil {
		obj["rolling_update"] = flattenRollingUpdateDaemonSet(in.RollingUpdate)
	}

	if len(in.Strategy) > 0 {
		obj["strategy"] = in.Strategy
	}

	return []interface{}{obj}
}

func flattenDeploymentStrategy(in *managementClient.DeploymentStrategy) []interface{} {
	obj := make(map[string]interface{})

	if in == nil {
		return []interface{}{}
	}

	if in.RollingUpdate != nil {
		obj["rolling_update"] = flattenRollingUpdateDeployment(in.RollingUpdate)
	}

	if len(in.Strategy) > 0 {
		obj["strategy"] = in.Strategy
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

func expandRollingUpdateDaemonSet(p []interface{}) *managementClient.RollingUpdateDaemonSet {
	obj := &managementClient.RollingUpdateDaemonSet{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["max_unavailable"].(int); ok && v > 0 {
		obj.MaxUnavailable = intstr.FromInt(v)
	}

	return obj
}

func expandRollingUpdateDeployment(p []interface{}) *managementClient.RollingUpdateDeployment {
	obj := &managementClient.RollingUpdateDeployment{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["max_surge"].(int); ok && v > 0 {
		obj.MaxSurge = intstr.FromInt(v)
	}

	if v, ok := in["max_unavailable"].(int); ok && v > 0 {
		obj.MaxUnavailable = intstr.FromInt(v)
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

func expandDaemonSetStrategy(p []interface{}) *managementClient.DaemonSetUpdateStrategy {
	obj := &managementClient.DaemonSetUpdateStrategy{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["rolling_update"].([]interface{}); ok && len(v) > 0 {
		obj.RollingUpdate = expandRollingUpdateDaemonSet(v)
	}

	if v, ok := in["strategy"].(string); ok && len(v) > 0 {
		obj.Strategy = v
	}

	return obj
}

func expandDeploymentStrategy(p []interface{}) *managementClient.DeploymentStrategy {
	obj := &managementClient.DeploymentStrategy{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["rolling_update"].([]interface{}); ok && len(v) > 0 {
		obj.RollingUpdate = expandRollingUpdateDeployment(v)
	}

	if v, ok := in["strategy"].(string); ok && len(v) > 0 {
		obj.Strategy = v
	}

	return obj
}
