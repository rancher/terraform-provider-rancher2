package rancher2

import (
    v1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenPodSecurityPolicyAllowedProcMountTypes(p []v1.ProcMountType) []string {
	if len(p) == 0 {
		return []string{}
	}

	out := make([]string, len(p))
	for i, in := range p {
		out[i] = string(in)
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyAllowedProcMountTypes(p []string) []v1.ProcMountType {

	if len(p) == 0 {
		return []v1.ProcMountType{}
	}

	obj := make([]v1.ProcMountType, len(p))

    for i, in := range p {
        obj[i] = v1.ProcMountType(in)
	}

	return obj
}