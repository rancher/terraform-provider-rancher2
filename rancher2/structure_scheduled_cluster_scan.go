package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenScheduledClusterScanConfig(in *managementClient.ScheduledClusterScanConfig) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})
	obj["cron_schedule"] = in.CronSchedule
	obj["retention"] = int(in.Retention)

	return []interface{}{obj}
}

func flattenScheduledClusterScan(in *managementClient.ScheduledClusterScan) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})
	obj["enabled"] = in.Enabled
	obj["scan_config"] = flattenClusterScanConfig(in.ScanConfig)
	obj["schedule_config"] = flattenScheduledClusterScanConfig(in.ScheduleConfig)

	return []interface{}{obj}
}

// Expanders

func expandScheduledClusterScanConfig(p []interface{}) *managementClient.ScheduledClusterScanConfig {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})
	obj := &managementClient.ScheduledClusterScanConfig{}
	if v, ok := in["cron_schedule"].(string); ok && len(v) > 0 {
		obj.CronSchedule = v
	}
	if v, ok := in["retention"].(int); ok {
		obj.Retention = int64(v)
	}

	return obj
}

func expandScheduledClusterScan(p []interface{}) *managementClient.ScheduledClusterScan {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	obj := &managementClient.ScheduledClusterScan{}
	obj.Enabled = in["enabled"].(bool)

	if v, ok := in["scan_config"].([]interface{}); ok && len(v) > 0 {
		obj.ScanConfig = expandClusterScanConfig(v)
	}

	if v, ok := in["schedule_config"].([]interface{}); ok && len(v) > 0 {
		obj.ScheduleConfig = expandScheduledClusterScanConfig(v)
	}

	return obj
}
