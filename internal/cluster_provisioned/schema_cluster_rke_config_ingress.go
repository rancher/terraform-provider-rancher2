package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	clusterRKEConfigIngressDNSPolicyClusterFirst            = "ClusterFirst"
	clusterRKEConfigIngressDNSPolicyClusterFirstWithHostNet = "ClusterFirstWithHostNet"
	clusterRKEConfigIngressDNSPolicyDefault                 = "Default"
	clusterRKEConfigIngressDNSPolicyNone                    = "None"
)

var (
	clusterRKEConfigIngressDNSPolicyList = []string{
		clusterRKEConfigIngressDNSPolicyClusterFirst,
		clusterRKEConfigIngressDNSPolicyClusterFirstWithHostNet,
		clusterRKEConfigIngressDNSPolicyDefault,
		clusterRKEConfigIngressDNSPolicyNone,
	}
)

//Schemas

func clusterRKEConfigIngressFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default_backend": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"http_port": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"https_port": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"network_mode": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"node_selector": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"dns_policy": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(clusterRKEConfigIngressDNSPolicyList, true),
		},
		"tolerations": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ingress add-on tolerations",
			Elem: &schema.Resource{
				Schema: tolerationFields(),
			},
		},
		"update_strategy": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Update daemon set strategy",
			Elem: &schema.Resource{
				Schema: daemonSetStrategyFields(),
			},
		},
	}
	return s
}
