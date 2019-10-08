package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenMonitoringInput(in *managementClient.MonitoringInput) []interface{} {
	obj := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	if len(in.Answers) > 0 {
		obj["answers"] = toMapInterface(in.Answers)
	}

	return []interface{}{obj}
}

// Expanders

func expandMonitoringInput(p []interface{}) *managementClient.MonitoringInput {
	obj := &managementClient.MonitoringInput{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["answers"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Answers = toMapString(v)
	}

	return obj
}
