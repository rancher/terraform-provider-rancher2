package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	globalDNSProviderAlidnsKind     = "alidns"
	globalDNSProviderCloudflareKind = "cloudflare"
	globalDNSProviderRoute53Kind    = "route53"
	globalDNSProviderRoute53Private = "private"
	globalDNSProviderRoute53Public  = "public"
)

var (
	globalDNSProviderKinds        = []string{globalDNSProviderAlidnsKind, globalDNSProviderCloudflareKind, globalDNSProviderRoute53Kind}
	globalDNSProviderRoute53Zones = []string{globalDNSProviderRoute53Private, globalDNSProviderRoute53Public}
)

// Schemas

func globalDNSProviderAliConfigSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"secret_key": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	return s
}

func globalDNSProviderCloudFareConfigSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"api_email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"api_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"proxy_setting": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	return s
}

func globalDNSProviderRoute53ConfigSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"secret_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"credentials_path": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/.aws",
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "us-west-2",
		},
		"role_arn": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"zone_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      globalDNSProviderRoute53Public,
			ValidateFunc: validation.StringInSlice(globalDNSProviderRoute53Zones, true),
		},
	}

	return s
}

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
		},
		"dns_provider": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"alidns_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"cloudflare_config", "route53_config"},
			Elem: &schema.Resource{
				Schema: globalDNSProviderAliConfigSchema(),
			},
		},
		"cloudflare_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"alidns_config", "route53_config"},
			Elem: &schema.Resource{
				Schema: globalDNSProviderCloudFareConfigSchema(),
			},
		},
		"route53_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"alidns_config", "cloudflare_config"},
			Elem: &schema.Resource{
				Schema: globalDNSProviderRoute53ConfigSchema(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
