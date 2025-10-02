package rancher2

import rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"

func flattenClusterV2RKEConfigDataDirectories(in rkev1.DataDirectories) []any {
	obj := make(map[string]any)

	if in.SystemAgent != "" {
		obj["system_agent"] = in.SystemAgent
	}

	if in.Provisioning != "" {
		obj["provisioning"] = in.Provisioning
	}

	if in.K8sDistro != "" {
		obj["k8s_distro"] = in.K8sDistro
	}

	return []any{obj}
}

func expandClusterV2RKEConfigDataDirectories(p []any) rkev1.DataDirectories {
	obj := rkev1.DataDirectories{}

	if p == nil || len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["system_agent"].(string); ok && v != "" {
		obj.SystemAgent = v
	}

	if v, ok := in["provisioning"].(string); ok && v != "" {
		obj.Provisioning = v
	}

	if v, ok := in["k8s_distro"].(string); ok && v != "" {
		obj.K8sDistro = v
	}

	return obj
}
