package rancher2

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Flatteners

func flattenAgentDeploymentCustomizationV2(in *v1.AgentDeploymentCustomization) []interface{} {
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

func expandAgentDeploymentCustomizationOverrideResourceRequirementsV2(p []interface{}) (*corev1.ResourceRequirements, error) {
	obj := &corev1.ResourceRequirements{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil, nil
	}

	for i := range p {
		in := p[i].(map[string]interface{})

		obj.Limits = make(map[corev1.ResourceName]resource.Quantity)
		obj.Requests = make(map[corev1.ResourceName]resource.Quantity)

		if v, ok := in["cpu_limit"].(string); ok && len(v) > 0 {
			value, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			obj.Limits[corev1.ResourceCPU] = *resource.NewQuantity(value, resource.DecimalSI)
		}

		if v, ok := in["cpu_request"].(string); ok && len(v) > 0 {
			value, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			obj.Requests[corev1.ResourceCPU] = *resource.NewQuantity(value, resource.DecimalSI)
		}

		if v, ok := in["memory_limit"].(string); ok && len(v) > 0 {
			value, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			obj.Limits[corev1.ResourceMemory] = *resource.NewQuantity(value, resource.DecimalSI)
		}

		if v, ok := in["memory_request"].(string); ok && len(v) > 0 {
			value, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			obj.Requests[corev1.ResourceMemory] = *resource.NewQuantity(value, resource.DecimalSI)
		}
	}

	return obj, nil
}

func expandAgentDeploymentCustomizationV2(p []interface{}) (*v1.AgentDeploymentCustomization, error) {
	obj := &v1.AgentDeploymentCustomization{}
	if len(p) == 0 || p[0] == nil {
		return obj, fmt.Errorf("agent deployment customization was nil")
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["append_tolerations"].(string); ok && len(v) > 0 {
		var appendTolerations []corev1.Toleration
		if err := json.Unmarshal([]byte(v), &appendTolerations); err != nil {
			return nil, err
		}
		obj.AppendTolerations = appendTolerations
	}

	if v, ok := in["override_affinity"].(string); ok && len(v) > 0 {
		var overrideAffinity *corev1.Affinity
		if err := json.Unmarshal([]byte(v), &overrideAffinity); err != nil {
			return nil, err
		}
		obj.OverrideAffinity = overrideAffinity
	}

	if v, ok := in["override_resource_requirements"].([]interface{}); ok && len(v) > 0 {
		overrideResourceRequirements, err := expandAgentDeploymentCustomizationOverrideResourceRequirementsV2(v)
		if err != nil {
			return nil, err
		}
		obj.OverrideResourceRequirements = overrideResourceRequirements
	}

	return obj, nil
}
