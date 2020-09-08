package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	cloudProviderAzureLoadBalancerSkuBasic    = "basic"
	cloudProviderAzureLoadBalancerSkuStandard = "standard"
)

var (
	cloudProviderAzureLoadBalancerSkuList = []string{
		cloudProviderAzureLoadBalancerSkuBasic,
		cloudProviderAzureLoadBalancerSkuStandard,
	}
)

//Schemas

func clusterRKEConfigCloudProviderAzureFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"aad_client_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"aad_client_secret": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"subscription_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"tenant_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"aad_client_cert_password": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"aad_client_cert_path": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"cloud": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_backoff": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_backoff_duration": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_backoff_exponent": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_backoff_jitter": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_backoff_retries": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_rate_limit": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_rate_limit_bucket": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_rate_limit_qps": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"load_balancer_sku": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      cloudProviderAzureLoadBalancerSkuBasic,
			Description:  "Load balancer type (basic | standard). Must be standard for auto-scaling",
			ValidateFunc: validation.StringInSlice(cloudProviderAzureLoadBalancerSkuList, true),
		},
		"location": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"maximum_load_balancer_rule_count": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"primary_availability_set_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"primary_scale_set_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"resource_group": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"route_table_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"security_group_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"use_instance_metadata": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"use_managed_identity_extension": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"vm_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vnet_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vnet_resource_group": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}
