package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	loggingSplunkKind = "splunk"
)

//Schemas

func loggingSplunkConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"token": {
			Type:      schema.TypeString,
			Required:  true,
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
		"index": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"source": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ssl_verify": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
