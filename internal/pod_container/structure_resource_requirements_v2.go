package rancher2

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Flatteners

func flattenResourceRequirementsV2(in *corev1.ResourceRequirements) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if in.Limits != nil {
		cpuLimitQuantity := in.Limits[corev1.ResourceCPU]
		cpuLimit := cpuLimitQuantity.String()
		if len(cpuLimit) > 0 {
			obj["cpu_limit"] = cpuLimit
		}

		memoryLimitQuantity := in.Limits[corev1.ResourceMemory]
		memoryLimit := memoryLimitQuantity.String()
		if len(memoryLimit) > 0 {
			obj["memory_limit"] = memoryLimit
		}
	}

	if in.Requests != nil {
		cpuRequestQuantity := in.Requests[corev1.ResourceCPU]
		cpuRequest := cpuRequestQuantity.String()
		if len(cpuRequest) > 0 {
			obj["cpu_request"] = cpuRequest
		}

		memoryRequestQuantity := in.Requests[corev1.ResourceMemory]
		memoryRequest := memoryRequestQuantity.String()
		if len(memoryRequest) > 0 {
			obj["memory_request"] = memoryRequest
		}
	}

	return []interface{}{obj}
}

// Expanders

func expandResourceRequirementsV2(p []interface{}) (*corev1.ResourceRequirements, error) {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil, nil
	}

	obj := &corev1.ResourceRequirements{}

	for i := range p {
		in := p[i].(map[string]interface{})

		obj.Limits = make(map[corev1.ResourceName]resource.Quantity)
		obj.Requests = make(map[corev1.ResourceName]resource.Quantity)

		if v, ok := in["cpu_limit"].(string); ok && len(v) > 0 {
			cpuLimit, err := resource.ParseQuantity(v)
			if err != nil {
				return nil, err
			}
			obj.Limits[corev1.ResourceCPU] = cpuLimit
		}

		if v, ok := in["cpu_request"].(string); ok && len(v) > 0 {
			cpuRequest, err := resource.ParseQuantity(v)
			if err != nil {
				return nil, err
			}
			obj.Requests[corev1.ResourceCPU] = cpuRequest
		}

		if v, ok := in["memory_limit"].(string); ok && len(v) > 0 {
			memoryLimit, err := resource.ParseQuantity(v)
			if err != nil {
				return nil, err
			}
			obj.Limits[corev1.ResourceMemory] = memoryLimit
		}

		if v, ok := in["memory_request"].(string); ok && len(v) > 0 {
			memoryRequest, err := resource.ParseQuantity(v)
			if err != nil {
				return nil, err
			}
			obj.Requests[corev1.ResourceMemory] = memoryRequest
		}
	}

	return obj, nil
}
