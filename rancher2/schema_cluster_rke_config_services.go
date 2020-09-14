package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigServicesFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"etcd": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesEtcdFields(),
			},
		},
		"kube_api": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIFieldsV0(),
			},
		},
		"kube_controller": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeControllerFields(),
			},
		},
		"kubelet": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeletFields(),
			},
		},
		"kubeproxy": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeproxyFields(),
			},
		},
		"scheduler": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesSchedulerFields(),
			},
		},
	}
	return s
}

func clusterRKEConfigServicesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"etcd": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesEtcdFields(),
			},
		},
		"kube_api": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIFields(),
			},
		},
		"kube_controller": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeControllerFields(),
			},
		},
		"kubelet": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeletFields(),
			},
		},
		"kubeproxy": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeproxyFields(),
			},
		},
		"scheduler": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesSchedulerFields(),
			},
		},
	}
	return s
}

// Used by datasource
func clusterRKEConfigServicesFieldsData() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"etcd": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesEtcdFields(),
			},
		},
		"kube_api": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIFieldsData(),
			},
		},
		"kube_controller": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeControllerFields(),
			},
		},
		"kubelet": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeletFields(),
			},
		},
		"kubeproxy": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeproxyFields(),
			},
		},
		"scheduler": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesSchedulerFields(),
			},
		},
	}
	return s
}
