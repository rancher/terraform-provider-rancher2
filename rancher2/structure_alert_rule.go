package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenEventRule(in *managementClient.EventRule) []interface{} {
	obj := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	if len(in.EventType) > 0 {
		obj["event_type"] = in.EventType
	}

	if len(in.ResourceKind) > 0 {
		obj["resource_kind"] = in.ResourceKind
	}

	return []interface{}{obj}
}

func flattenMetricRule(in *managementClient.MetricRule) []interface{} {
	obj := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	if len(in.Comparison) > 0 {
		obj["comparison"] = in.Comparison
	}

	if len(in.Duration) > 0 {
		obj["duration"] = in.Duration
	}

	if len(in.Expression) > 0 {
		obj["expression"] = in.Expression
	}

	if in.ThresholdValue > 0 {
		obj["threshold_value"] = in.ThresholdValue
	}

	if len(in.Description) > 0 {
		obj["description"] = in.Description
	}

	return []interface{}{obj}
}

func flattenNodeRule(in *managementClient.NodeRule) []interface{} {
	obj := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	if in.CPUThreshold > 0 {
		obj["cpu_threshold"] = int(in.CPUThreshold)
	}

	if len(in.Condition) > 0 {
		obj["condition"] = in.Condition
	}

	if in.MemThreshold > 0 {
		obj["mem_threshold"] = int(in.MemThreshold)
	}

	if len(in.NodeID) > 0 {
		obj["node_id"] = in.NodeID
	}

	if len(in.Selector) > 0 {
		obj["selector"] = toMapInterface(in.Selector)
	}

	return []interface{}{obj}
}

func flattenPodRule(in *managementClient.PodRule) []interface{} {
	obj := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	if len(in.Condition) > 0 {
		obj["condition"] = in.Condition
	}

	if len(in.PodID) > 0 {
		obj["pod_id"] = in.PodID
	}

	if in.RestartIntervalSeconds > 0 {
		obj["restart_interval_seconds"] = int(in.RestartIntervalSeconds)
	}

	if in.RestartTimes > 0 {
		obj["restart_times"] = int(in.RestartTimes)
	}

	return []interface{}{obj}
}

func flattenSystemServiceRule(in *managementClient.SystemServiceRule) []interface{} {
	obj := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	if len(in.Condition) > 0 {
		obj["condition"] = in.Condition
	}

	return []interface{}{obj}
}

func flattenWorkloadRule(in *managementClient.WorkloadRule) []interface{} {
	obj := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	if in.AvailablePercentage > 0 {
		obj["available_percentage"] = int(in.AvailablePercentage)
	}

	if len(in.Selector) > 0 {
		obj["selector"] = toMapInterface(in.Selector)
	}

	if len(in.WorkloadID) > 0 {
		obj["workload_id"] = in.WorkloadID
	}

	return []interface{}{obj}
}

// Expanders

func expandEventRule(p []interface{}) *managementClient.EventRule {
	obj := &managementClient.EventRule{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["event_type"].(string); ok && len(v) > 0 {
		obj.EventType = v
	}

	if v, ok := in["resource_kind"].(string); ok && len(v) > 0 {
		obj.ResourceKind = v
	}

	return obj
}

func expandMetricRule(p []interface{}) *managementClient.MetricRule {
	obj := &managementClient.MetricRule{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["comparison"].(string); ok && len(v) > 0 {
		obj.Comparison = v
	}

	if v, ok := in["duration"].(string); ok && len(v) > 0 {
		obj.Duration = v
	}

	if v, ok := in["expression"].(string); ok && len(v) > 0 {
		obj.Expression = v
	}

	if v, ok := in["threshold_value"].(float64); ok && v > 0 {
		obj.ThresholdValue = v
	}

	if v, ok := in["description"].(string); ok && len(v) > 0 {
		obj.Description = v
	}

	return obj
}

func expandNodeRule(p []interface{}) *managementClient.NodeRule {
	obj := &managementClient.NodeRule{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cpu_threshold"].(int); ok && v > 0 {
		obj.CPUThreshold = int64(v)
	}

	if v, ok := in["condition"].(string); ok && len(v) > 0 {
		obj.Condition = v
	}

	if v, ok := in["mem_threshold"].(int); ok && v > 0 {
		obj.MemThreshold = int64(v)
	}

	if v, ok := in["node_id"].(string); ok && len(v) > 0 {
		obj.NodeID = v
	}

	if v, ok := in["selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Selector = toMapString(v)
	}

	return obj
}

func expandPodRule(p []interface{}) *managementClient.PodRule {
	obj := &managementClient.PodRule{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["condition"].(string); ok && len(v) > 0 {
		obj.Condition = v
	}

	if v, ok := in["pod_id"].(string); ok && len(v) > 0 {
		obj.PodID = v
	}

	if v, ok := in["restart_interval_seconds"].(int); ok && v > 0 {
		obj.RestartIntervalSeconds = int64(v)
	}

	if v, ok := in["restart_times"].(int); ok && v > 0 {
		obj.RestartTimes = int64(v)
	}

	return obj
}

func expandSystemServiceRule(p []interface{}) *managementClient.SystemServiceRule {
	obj := &managementClient.SystemServiceRule{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["condition"].(string); ok && len(v) > 0 {
		obj.Condition = v
	}

	return obj
}

func expandWorkloadRule(p []interface{}) *managementClient.WorkloadRule {
	obj := &managementClient.WorkloadRule{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["available_percentage"].(int); ok && v > 0 {
		obj.AvailablePercentage = int64(v)
	}

	if v, ok := in["selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Selector = toMapString(v)
	}

	if v, ok := in["workload_id"].(string); ok && len(v) > 0 {
		obj.WorkloadID = v
	}

	return obj
}
