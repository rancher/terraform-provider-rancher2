package rancher2

import (
	"encoding/json"
	"fmt"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAgentDeploymentCustomization(in *managementClient.AgentDeploymentCustomization) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AppendTolerations) > 0 {
		obj["append_tolerations"] = in.AppendTolerations
	}

	if in.OverrideAffinity != nil {
		obj["override_affinity"] = in.OverrideAffinity
	}

	if in.OverrideResourceRequirements != nil {
		obj["override_resource_requirements"] = in.OverrideResourceRequirements
	}

	return []interface{}{obj}
}

// Expanders

func expandAgentDeploymentCustomizationOverrideResourceRequirements(p []interface{}) *managementClient.ResourceRequirements {
	obj := &managementClient.ResourceRequirements{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	for i := range p {
		in := p[i].(map[string]interface{})

		obj.Limits = make(map[string]string)
		obj.Requests = make(map[string]string)

		if v, ok := in["cpu_limit"].(string); ok && len(v) > 0 {
			obj.Limits["cpu"] = v
		}

		if v, ok := in["cpu_request"].(string); ok && len(v) > 0 {
			obj.Requests["cpu"] = v
		}

		if v, ok := in["memory_limit"].(string); ok && len(v) > 0 {
			obj.Limits["memory"] = v
		}

		if v, ok := in["memory_request"].(string); ok && len(v) > 0 {
			obj.Requests["memory"] = v
		}
	}

	return obj
}

func expandAgentDeploymentCustomization(p []interface{}) (*managementClient.AgentDeploymentCustomization, error) {
	obj := &managementClient.AgentDeploymentCustomization{}
	if len(p) == 0 || p[0] == nil {
		return obj, fmt.Errorf("agent deployment customization was nil")
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["append_tolerations"].(string); ok && len(v) > 0 {
		var appendTolerations []managementClient.Toleration
		if err := json.Unmarshal([]byte(v), &appendTolerations); err != nil {
			return nil, err
		}
		obj.AppendTolerations = appendTolerations
	}

	if v, ok := in["override_affinity"].(string); ok && len(v) > 0 {
		var overrideAffinity *managementClient.Affinity
		if err := json.Unmarshal([]byte(v), &overrideAffinity); err != nil {
			return nil, err
		}
		obj.OverrideAffinity = overrideAffinity
	}

	if v, ok := in["override_resource_requirements"].([]interface{}); ok && len(v) > 0 {
		obj.OverrideResourceRequirements = expandAgentDeploymentCustomizationOverrideResourceRequirements(v)
	}

	return obj, nil
}
