package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func globalDNSProviderAliConfigSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
		},
		"secret_key": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
		},
	}

	return s
}

func globalDNSProviderCloudFareConfigSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"api_email": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
		},
		"api_key": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
		},
		"proxy_setting": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: false,
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
			ForceNew: false,
		},
		"secret_key": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
		},
		"credentials_path": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: false,
			Default:  "/.aws",
		},
		"zone_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"public", "private"}, true),
		},
		"role_arn": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"region": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}

	return s
}
