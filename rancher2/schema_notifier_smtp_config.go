package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierSMTPConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default_recipient": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "SMTP default recipient address",
		},
		"host": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "SMTP host",
		},
		"port": &schema.Schema{
			Type:        schema.TypeInt,
			Required:    true,
			Description: "SMTP port",
		},
		"sender": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "SMTP sender",
		},
		"password": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SMTP password",
		},
		"tls": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "SMTP TLS",
		},
		"username": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SMTP username",
		},
	}

	return s
}
