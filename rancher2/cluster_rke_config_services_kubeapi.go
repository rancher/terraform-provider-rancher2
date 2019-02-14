package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schema

func kubeAPIFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"extra_binds": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"pod_security_policy": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"service_cluster_ip_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"service_node_port_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenKubeAPI(in *managementClient.KubeAPIService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["pod_security_policy"] = in.PodSecurityPolicy

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	if len(in.ServiceNodePortRange) > 0 {
		obj["service_node_port_range"] = in.ServiceNodePortRange
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandKubeAPI(p []interface{}) (*managementClient.KubeAPIService, error) {
	obj := &managementClient.KubeAPIService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["pod_security_policy"].(bool); ok {
		obj.PodSecurityPolicy = v
	}

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	if v, ok := in["service_node_port_range"].(string); ok && len(v) > 0 {
		obj.ServiceNodePortRange = v
	}

	return obj, nil
}
