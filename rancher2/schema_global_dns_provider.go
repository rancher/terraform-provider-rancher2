package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// Schemas
func globalDNSProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"root_domain": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"dns_provider": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.StringInSlice([]string{"route53", "cloudflare", "alidns"}, true),
		},
		"annotations": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"route53_config": {
			Type:          schema.TypeSet,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"alidns_config", "cloudflare_config"},
			Elem: &schema.Resource{
				Schema: globalDNSProviderRoute53ConfigSchema(),
			},
		},
		"alidns_config": {
			Type:          schema.TypeSet,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"alidns_config", "cloudflare_config"},
			Elem: &schema.Resource{
				Schema: globalDNSProviderAliConfigSchema(),
			},
		},
		"cloudflare_config": {
			Type:          schema.TypeSet,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"alidns_config", "cloudflare_config"},
			Elem: &schema.Resource{
				Schema: globalDNSProviderCloudFareConfigSchema(),
			},
		},
	}

	return s
}
