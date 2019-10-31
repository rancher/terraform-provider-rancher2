package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	loggingElasticsearchKind = "elasticsearch"
)

func loggingElasticsearchConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"auth_password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"auth_username": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_cert": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_key_pass": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"date_format": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "YYYY-MM-DD",
		},
		"index_prefix": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "local",
		},
		"ssl_verify": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"ssl_version": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}
