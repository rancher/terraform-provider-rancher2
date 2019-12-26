package rancher2

import (
    v1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenPodSecurityPolicyCapabilities(v []v1.Capability) []string {
	if len(v) == 0 {
		return []string{}
	}

	out := make([]string, len(v))
	for i, in := range v {
		out[i] = string(in)
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyCapabilities(v []string) []v1.Capability {

	if len(v) == 0 {
		return []v1.Capability{}
	}

	obj := make([]v1.Capability, len(v))

    for i, in := range v {
        obj[i] = v1.Capability(in)
	}

	return obj
}
