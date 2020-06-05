package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierSMTPConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default_recipient": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SMTP default recipient address",
		},
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SMTP host",
		},
		"port": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "SMTP port",
		},
		"sender": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SMTP sender",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SMTP password",
		},
		"tls": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "SMTP TLS",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SMTP username",
		},
	}

	return s
}
