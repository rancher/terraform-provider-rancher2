package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	bootstrapDefaultTokenDesc   = "Terraform bootstrap admin token"
	bootstrapDefaultSessionDesc = "Terraform bootstrap admin session"
	bootstrapDefaultUser        = "admin"
	bootstrapDefaultPassword    = "admin"
	bootstrapDefaultTTL         = "60000"
	bootstrapSettingUILanding   = "ui-default-landing"
	bootstrapSettingURL         = "server-url"
	bootstrapSettingTelemetry   = "telemetry-opt"
	bootstrapUILandingExplorer  = "vue"
	bootstrapUILandingManager   = "ember"
)

var (
	bootstrapUILandingKinds = []string{bootstrapUILandingExplorer, bootstrapUILandingManager}
)

//Schemas

func bootstrapFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"current_password": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"initial_password": {
			Type:      schema.TypeString,
			Optional:  true,
			Default:   bootstrapDefaultPassword,
			Sensitive: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"token_ttl": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"token": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"token_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"token_update": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"telemetry": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"ui_default_landing": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      bootstrapUILandingManager,
			ValidateFunc: validation.StringInSlice(bootstrapUILandingKinds, true),
		},
		"temp_token": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"temp_token_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"user": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return s
}
