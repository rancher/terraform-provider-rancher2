package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func etcdFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"retention": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"snapshot": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func kubeApiFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"extra_args": {
			Type:     schema.TypeMap,
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

func kubeControllerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_cidr": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"service_cluster_ip_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

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
	}
	return s
}

func kubeproxyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

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
				Schema: kubeApiFields(),
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

func flattenEtcd(in *managementClient.ETCDService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.Retention) > 0 {
		obj["retention"] = in.Retention
	}

	obj["snapshot"] = in.Snapshot

	return []interface{}{obj}, nil
}

func flattenKubeApi(in *managementClient.KubeAPIService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
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

func flattenKubeController(in *managementClient.KubeControllerService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ClusterCIDR) > 0 {
		obj["cluster_cidr"] = in.ClusterCIDR
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	return []interface{}{obj}, nil
}

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

	return []interface{}{obj}, nil
}

func flattenKubeproxy(in *managementClient.KubeproxyService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	return []interface{}{obj}, nil
}

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
		kubeApi, err := flattenKubeApi(in.KubeAPI)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kube_api"] = kubeApi
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

func expandEtcd(p []interface{}) (*managementClient.ETCDService, error) {
	obj := &managementClient.ETCDService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["retention"].(string); ok && len(v) > 0 {
		obj.Retention = v
	}

	if v, ok := in["snapshot"].(bool); ok {
		obj.Snapshot = v
	}

	return obj, nil
}

func expandKubeApi(p []interface{}) (*managementClient.KubeAPIService, error) {
	obj := &managementClient.KubeAPIService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
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

func expandKubeController(p []interface{}) (*managementClient.KubeControllerService, error) {
	obj := &managementClient.KubeControllerService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_cidr"].(string); ok && len(v) > 0 {
		obj.ClusterCIDR = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	return obj, nil
}

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

	return obj, nil
}

func expandKubeproxy(p []interface{}) (*managementClient.KubeproxyService, error) {
	obj := &managementClient.KubeproxyService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	return obj, nil
}

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
		kubeApi, err := expandKubeApi(v)
		if err != nil {
			return obj, err
		}
		obj.KubeAPI = kubeApi
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
