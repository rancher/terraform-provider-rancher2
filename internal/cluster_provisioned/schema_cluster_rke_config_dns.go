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

func clusterRKEConfigDNSNodelocalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ip_address": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"node_selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Node selector key pair",
		},
	}
	return s
}

func clusterRKEConfigDNSLinearAutoscalerParamsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cores_per_replica": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		"nodes_per_replica": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		"prevent_single_point_failure": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"min": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"max": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	return s
}

func clusterRKEConfigDNSFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"node_selector": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"nodelocal": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Nodelocal dns",
			Elem: &schema.Resource{
				Schema: clusterRKEConfigDNSNodelocalFields(),
			},
		},
		"linear_autoscaler_params": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Linear Autoscaler Params",
			Elem: &schema.Resource{
				Schema: clusterRKEConfigDNSLinearAutoscalerParamsFields(),
			},
		},
		"options": {
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
		"tolerations": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "DNS service tolerations",
			Elem: &schema.Resource{
				Schema: tolerationFields(),
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
		"update_strategy": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Update deployment strategy",
			Elem: &schema.Resource{
				Schema: deploymentStrategyFields(),
			},
		},
	}

	return s
}
