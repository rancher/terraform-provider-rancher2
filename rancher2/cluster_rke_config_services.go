package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func servicesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"etcd": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: etcdFields(),
			},
		},
		"kube_api": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: kubeAPIFields(),
			},
		},
		"kube_controller": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: kubeControllerFields(),
			},
		},
		"kubelet": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: kubeletFields(),
			},
		},
		"kubeproxy": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: kubeproxyFields(),
			},
		},
	}
	return s
}

// Flatteners

func flattenServices(in *managementClient.RKEConfigServices) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.Etcd != nil {
		etcd, err := flattenEtcd(in.Etcd)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["etcd"] = etcd
	}

	if in.KubeAPI != nil {
		kubeAPI, err := flattenKubeAPI(in.KubeAPI)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kube_api"] = kubeAPI
	}

	if in.KubeController != nil {
		kubeController, err := flattenKubeController(in.KubeController)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kube_controller"] = kubeController
	}

	if in.Kubelet != nil {
		kubelet, err := flattenKubelet(in.Kubelet)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kubelet"] = kubelet
	}

	if in.Kubeproxy != nil {
		kubeproxy, err := flattenKubeproxy(in.Kubeproxy)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kubeproxy"] = kubeproxy
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandServices(p []interface{}) (*managementClient.RKEConfigServices, error) {
	obj := &managementClient.RKEConfigServices{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["etcd"].([]interface{}); ok && len(v) > 0 {
		etcd, err := expandEtcd(v)
		if err != nil {
			return obj, err
		}
		obj.Etcd = etcd
	}

	if v, ok := in["kube_api"].([]interface{}); ok && len(v) > 0 {
		kubeAPI, err := expandKubeAPI(v)
		if err != nil {
			return obj, err
		}
		obj.KubeAPI = kubeAPI
	}

	if v, ok := in["kube_controller"].([]interface{}); ok && len(v) > 0 {
		kubeController, err := expandKubeController(v)
		if err != nil {
			return obj, err
		}
		obj.KubeController = kubeController
	}

	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		kubelet, err := expandKubelet(v)
		if err != nil {
			return obj, err
		}
		obj.Kubelet = kubelet
	}

	if v, ok := in["kubeproxy"].([]interface{}); ok && len(v) > 0 {
		kubeproxy, err := expandKubeproxy(v)
		if err != nil {
			return obj, err
		}
		obj.Kubeproxy = kubeproxy
	}

	return obj, nil
}
