package rancher2

import (
	"reflect"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

// Flatteners

func flattenClusterV2RKEConfigUpgradeStrategyDrainOptions(in rkev1.DrainOptions) []interface{} {
	if reflect.ValueOf(in).IsZero() {
		return nil
	}

	obj := make(map[string]interface{})

	obj["enabled"] = in.Enabled
	obj["force"] = in.Force
	if in.IgnoreDaemonSets != nil {
		obj["ignore_daemon_sets"] = *in.IgnoreDaemonSets
	}
	obj["ignore_errors"] = in.IgnoreErrors
	obj["delete_empty_dir_data"] = in.DeleteEmptyDirData
	obj["disable_eviction"] = in.DisableEviction
	if in.GracePeriod > 0 {
		obj["grace_period"] = in.GracePeriod
	}
	if in.Timeout > 0 {
		obj["timeout"] = in.Timeout
	}
	if in.SkipWaitForDeleteTimeoutSeconds > 0 {
		obj["skip_wait_for_delete_timeout_seconds"] = in.SkipWaitForDeleteTimeoutSeconds
	}

	return []interface{}{obj}
}

func flattenClusterV2RKEConfigUpgradeStrategy(in rkev1.ClusterUpgradeStrategy) []interface{} {
	if reflect.ValueOf(in).IsZero() {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.ControlPlaneConcurrency) > 0 {
		obj["control_plane_concurrency"] = in.ControlPlaneConcurrency
	}
	obj["control_plane_drain_options"] = flattenClusterV2RKEConfigUpgradeStrategyDrainOptions(in.ControlPlaneDrainOptions)
	if len(in.WorkerConcurrency) > 0 {
		obj["worker_concurrency"] = in.WorkerConcurrency
	}
	obj["worker_drain_options"] = flattenClusterV2RKEConfigUpgradeStrategyDrainOptions(in.WorkerDrainOptions)

	return []interface{}{obj}
}

// Expanders

func expandClusterV2RKEConfigUpgradeStrategyDrainOptions(p []interface{}) rkev1.DrainOptions {
	if p == nil || len(p) == 0 || p[0] == nil {
		return rkev1.DrainOptions{}
	}

	obj := rkev1.DrainOptions{}

	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}
	if v, ok := in["force"].(bool); ok {
		obj.Force = v
	}
	if v, ok := in["ignore_daemon_sets"].(bool); ok {
		obj.IgnoreDaemonSets = &v
	}
	if v, ok := in["ignore_errors"].(bool); ok {
		obj.IgnoreErrors = v
	}
	if v, ok := in["delete_empty_dir_data"].(bool); ok {
		obj.DeleteEmptyDirData = v
	}
	if v, ok := in["disable_eviction"].(bool); ok {
		obj.DisableEviction = v
	}
	if v, ok := in["grace_period"].(int); ok && v > 0 {
		obj.GracePeriod = v
	}
	if v, ok := in["timeout"].(int); ok && v > 0 {
		obj.Timeout = v
	}
	if v, ok := in["skip_wait_for_delete_timeout_seconds"].(int); ok && v > 0 {
		obj.SkipWaitForDeleteTimeoutSeconds = v
	}

	return obj
}

func expandClusterV2RKEConfigUpgradeStrategy(p []interface{}) rkev1.ClusterUpgradeStrategy {
	if p == nil || len(p) == 0 || p[0] == nil {
		return rkev1.ClusterUpgradeStrategy{}
	}

	obj := rkev1.ClusterUpgradeStrategy{}

	in := p[0].(map[string]interface{})

	if v, ok := in["control_plane_concurrency"].(string); ok && len(v) > 0 {
		obj.ControlPlaneConcurrency = v
	}
	if v, ok := in["control_plane_drain_options"].([]interface{}); ok && len(v) > 0 {
		obj.ControlPlaneDrainOptions = expandClusterV2RKEConfigUpgradeStrategyDrainOptions(v)
	}
	if v, ok := in["worker_concurrency"].(string); ok && len(v) > 0 {
		obj.WorkerConcurrency = v
	}
	if v, ok := in["worker_drain_options"].([]interface{}); ok && len(v) > 0 {
		obj.WorkerDrainOptions = expandClusterV2RKEConfigUpgradeStrategyDrainOptions(v)
	}

	return obj
}
