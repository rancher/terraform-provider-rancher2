package rancher2

import (
	"encoding/json"

	"github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	corev1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenAgentDeploymentCustomizationV2(in *v1.AgentDeploymentCustomization) ([]interface{}, error) {
	if in == nil {
		return []interface{}{}, nil
	}

	obj := make(map[string]interface{})

	if len(in.AppendTolerations) > 0 {
		obj["append_tolerations"] = flattenTolerationsV2(in.AppendTolerations)
	}

	if in.OverrideAffinity != nil {
		b, err := interfaceToJSON(in.OverrideAffinity)
		if err != nil {
			return []interface{}{}, err
		}
		obj["override_affinity"] = b // TODO - ANDY, changed it to keep it simpler instead of creating a new schema.  The expand reads a json.
	}

	if in.OverrideResourceRequirements != nil {
		obj["override_resource_requirements"] = flattenResourceRequirementsV2(in.OverrideResourceRequirements) // TODO - ANDY added a flatten
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandAgentDeploymentCustomizationV2(p []interface{}) (*v1.AgentDeploymentCustomization, error) {
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

	return obj, nil
}
