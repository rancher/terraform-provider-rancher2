package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyVolumes(v []policyv1.FSType) []string {
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

func expandPodSecurityPolicyVolumes(v []string) []policyv1.FSType {

	if len(v) == 0 {
		return []policyv1.FSType{}
	}

	obj := make([]policyv1.FSType, len(v))

    for i, in := range v {
        obj[i] = policyv1.FSType(in)
	}

	return obj
}
