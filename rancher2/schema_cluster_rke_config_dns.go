package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	clusterRKEDNSProviderKube = "kube-dns"
	clusterRKEDNSProviderCore = "coredns"
	clusterRKEDNSProviderNone = "none"
)

var (
	clusterRKEDNSProviderList = []string{clusterRKEDNSProviderKube, clusterRKEDNSProviderCore, clusterRKEDNSProviderNone}
)

//Schemas

func clusterRKEConfigDNSFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"node_selector": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"provider": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      clusterRKEDNSProviderCore,
			ValidateFunc: validation.StringInSlice(clusterRKEDNSProviderList, true),
		},
		"reverse_cidrs": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"upstream_nameservers": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	return s
}
