package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schema

func kubeletFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_dns_server": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"cluster_domain": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
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
		"fail_swap_on": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"infra_container_image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenKubelet(in *managementClient.KubeletService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ClusterDNSServer) > 0 {
		obj["cluster_dns_server"] = in.ClusterDNSServer
	}

	if len(in.ClusterDomain) > 0 {
		obj["cluster_domain"] = in.ClusterDomain
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

	obj["fail_swap_on"] = in.FailSwapOn

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.InfraContainerImage) > 0 {
		obj["infra_container_image"] = in.InfraContainerImage
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandKubelet(p []interface{}) (*managementClient.KubeletService, error) {
	obj := &managementClient.KubeletService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_dns_server"].(string); ok && len(v) > 0 {
		obj.ClusterDNSServer = v
	}

	if v, ok := in["cluster_domain"].(string); ok && len(v) > 0 {
		obj.ClusterDomain = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["fail_swap_on"].(bool); ok {
		obj.FailSwapOn = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["infra_container_image"].(string); ok && len(v) > 0 {
		obj.InfraContainerImage = v
	}

	return obj, nil
}
