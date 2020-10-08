package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigServices(in *managementClient.RKEConfigServices, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.Etcd != nil {
		v, ok := obj["etcd"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		etcd, err := flattenClusterRKEConfigServicesEtcd(in.Etcd, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["etcd"] = etcd
	}

	if in.KubeAPI != nil {
		kubeAPI, err := flattenClusterRKEConfigServicesKubeAPI(in.KubeAPI)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kube_api"] = kubeAPI
	}

	if in.KubeController != nil {
		kubeController, err := flattenClusterRKEConfigServicesKubeController(in.KubeController)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kube_controller"] = kubeController
	}

	if in.Kubelet != nil {
		kubelet, err := flattenClusterRKEConfigServicesKubelet(in.Kubelet)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kubelet"] = kubelet
	}

	if in.Kubeproxy != nil {
		kubeproxy, err := flattenClusterRKEConfigServicesKubeproxy(in.Kubeproxy)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kubeproxy"] = kubeproxy
	}

	if in.Scheduler != nil {
		scheduler, err := flattenClusterRKEConfigServicesScheduler(in.Scheduler)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["scheduler"] = scheduler
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigServices(p []interface{}) (*managementClient.RKEConfigServices, error) {
	obj := &managementClient.RKEConfigServices{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["etcd"].([]interface{}); ok && len(v) > 0 {
		etcd, err := expandClusterRKEConfigServicesEtcd(v)
		if err != nil {
			return obj, err
		}
		obj.Etcd = etcd
	}

	if v, ok := in["kube_api"].([]interface{}); ok && len(v) > 0 {
		kubeAPI, err := expandClusterRKEConfigServicesKubeAPI(v)
		if err != nil {
			return obj, err
		}
		obj.KubeAPI = kubeAPI
	}

	if v, ok := in["kube_controller"].([]interface{}); ok && len(v) > 0 {
		kubeController, err := expandClusterRKEConfigServicesKubeController(v)
		if err != nil {
			return obj, err
		}
		obj.KubeController = kubeController
	}

	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		kubelet, err := expandClusterRKEConfigServicesKubelet(v)
		if err != nil {
			return obj, err
		}
		obj.Kubelet = kubelet
	}

	if v, ok := in["kubeproxy"].([]interface{}); ok && len(v) > 0 {
		kubeproxy, err := expandClusterRKEConfigServicesKubeproxy(v)
		if err != nil {
			return obj, err
		}
		obj.Kubeproxy = kubeproxy
	}

	if v, ok := in["scheduler"].([]interface{}); ok && len(v) > 0 {
		scheduler, err := expandClusterRKEConfigServicesScheduler(v)
		if err != nil {
			return obj, err
		}
		obj.Scheduler = scheduler
	}

	return obj, nil
}
