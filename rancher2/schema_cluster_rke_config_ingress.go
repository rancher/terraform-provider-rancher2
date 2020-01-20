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
		"extra_args": {
			Type:     schema.TypeMap,
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
	}
	return s
}
