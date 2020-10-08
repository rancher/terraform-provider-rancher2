package rancher2

import (
	"reflect"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	monitoringInputAnnotation = "field.cattle.io/overwriteAppAnswers"
	monitoringActionDisable   = "disableMonitoring"
	monitoringActionEdit      = "editMonitoring"
	monitoringActionEnable    = "enableMonitoring"
)

// Flatteners

func flattenMonitoringInput(in *managementClient.MonitoringInput) []interface{} {
	if in == nil || reflect.DeepEqual(in, &managementClient.MonitoringInput{}) {
		return []interface{}{}
	}
	obj := map[string]interface{}{}

	if len(in.Answers) > 0 {
		obj["answers"] = toMapInterface(in.Answers)
	}

	if len(in.Version) > 0 {
		obj["version"] = in.Version
	}

	return []interface{}{obj}
}

// Expanders

func expandMonitoringInput(p []interface{}) *managementClient.MonitoringInput {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &managementClient.MonitoringInput{}
	in := p[0].(map[string]interface{})

	if v, ok := in["answers"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Answers = toMapString(v)
	}

	if v, ok := in["version"].(string); ok && len(v) > 0 {
		obj.Version = v
	}

	return obj
}
