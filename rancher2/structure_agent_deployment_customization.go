package rancher2

import (
	"encoding/json"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAgentDeploymentCustomization(in *managementClient.AgentDeploymentCustomization, includeScheduling bool) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.AppendTolerations) > 0 {
		obj["append_tolerations"] = flattenTolerations(in.AppendTolerations)
	}

	if in.OverrideAffinity != nil {
		overrideAffinity, _ := json.Marshal(in.OverrideAffinity)
		obj["override_affinity"] = string(overrideAffinity)
	}

	if in.OverrideResourceRequirements != nil {
		obj["override_resource_requirements"] = flattenResourceRequirements(in.OverrideResourceRequirements)
	}

	if includeScheduling {
		obj["scheduling_customization"] = flattenAgentSchedulingCustomization(in.SchedulingCustomization)
	}

	return []interface{}{obj}
}

// Expanders

func expandAgentDeploymentCustomization(p []interface{}, includeScheduling bool) (*managementClient.AgentDeploymentCustomization, error) {
	if len(p) == 0 || p[0] == nil {
		return nil, nil
	}

	obj := &managementClient.AgentDeploymentCustomization{}

	in := p[0].(map[string]interface{})

	if v, ok := in["append_tolerations"].([]interface{}); ok && len(v) > 0 {
		obj.AppendTolerations = expandTolerations(v)
	}

	if v, ok := in["override_affinity"].(string); ok && len(v) > 0 {
		var overrideAffinity *managementClient.Affinity
		if err := json.Unmarshal([]byte(v), &overrideAffinity); err != nil {
			return nil, err
		}
		obj.OverrideAffinity = overrideAffinity
	}

	if v, ok := in["override_resource_requirements"].([]interface{}); ok && len(v) > 0 {
		obj.OverrideResourceRequirements = expandResourceRequirements(v)
	}

	if includeScheduling {
		if v, ok := in["scheduling_customization"].([]interface{}); ok {
			obj.SchedulingCustomization = expandAgentSchedulingCustomization(v)
		}
	}

	return obj, nil
}
