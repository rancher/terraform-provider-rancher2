package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	bootstrapDefaultTokenDesc   = "Terraform bootstrap admin token"
	bootstrapDefaultSessionDesc = "Terraform bootstrap admin session"
	bootstrapDefaultUser        = "admin"
	bootstrapDefaultPassword    = "admin"
	bootstrapDefaultTTL         = "60000"
	bootstrapSettingURL         = "server-url"
	bootstrapSettingTelemetry   = "telemetry-opt"
)

//Schemas

func bootstrapFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"current_password": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"password": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"token_ttl": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"token": &schema.Schema{
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"token_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"token_update": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"telemetry": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"temp_token": &schema.Schema{
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"temp_token_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"user": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return s
}
