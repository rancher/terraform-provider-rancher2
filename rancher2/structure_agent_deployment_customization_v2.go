package rancher2

import (
	"encoding/json"

	"github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	corev1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenAgentDeploymentCustomizationV2(in *v1.AgentDeploymentCustomization, includeScheduling bool) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.AppendTolerations) > 0 {
		obj["append_tolerations"] = flattenTolerationsV2(in.AppendTolerations)
	}

	if in.OverrideAffinity != nil {
		obj["override_affinity"] = in.OverrideAffinity
	}

	if in.OverrideResourceRequirements != nil {
		obj["override_resource_requirements"] = in.OverrideResourceRequirements
	}

	if includeScheduling {
		obj["scheduling_customization"] = flattenAgentSchedulingCustomizationV2(in.SchedulingCustomization)
	}

	return []interface{}{obj}
}

// Expanders

func expandAgentDeploymentCustomizationV2(p []interface{}, includeScheduling bool) (*v1.AgentDeploymentCustomization, error) {
	if len(p) == 0 || p[0] == nil {
		return nil, nil
	}

	obj := &v1.AgentDeploymentCustomization{}

	in := p[0].(map[string]interface{})

	if v, ok := in["append_tolerations"].([]interface{}); ok && len(v) > 0 {
		obj.AppendTolerations = expandTolerationsV2(v)
	}

	if v, ok := in["override_affinity"].(string); ok && len(v) > 0 {
		var overrideAffinity *corev1.Affinity
		if err := json.Unmarshal([]byte(v), &overrideAffinity); err != nil {
			return nil, err
		}
		obj.OverrideAffinity = overrideAffinity
	}

	if v, ok := in["override_resource_requirements"].([]interface{}); ok && len(v) > 0 {
		overrideResourceRequirements, err := expandResourceRequirementsV2(v)
		if err != nil {
			return nil, err
		}
		obj.OverrideResourceRequirements = overrideResourceRequirements
	}

	if includeScheduling {
		if v, ok := in["scheduling_customization"].([]interface{}); ok && len(v) > 0 {
			obj.SchedulingCustomization = expandAgentSchedulingCustomizationV2(v)
		}
	}

	return obj, nil
}
