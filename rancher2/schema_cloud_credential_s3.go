package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func cloudCredentialS3Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "AWS Access Key",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "AWS Secret Key",
		},
		"default_bucket": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS default bucket",
		},
		"default_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS default endpoint",
		},
		"default_endpoint_ca": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "AWS default endpoint CA",
		},
		"default_folder": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS default folder",
		},
		"default_region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS default region",
		},
		"default_skip_ssl_verify": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "AWS default skip ssl verify",
		},
	}

	return s
}
