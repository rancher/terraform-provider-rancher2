package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierDingtalkConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Webhook URL",
		},
		"proxy_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Dingtalk proxy URL",
		},
		"secret": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Required for webhook with sign enabled",
		},
	}

	return s
}
