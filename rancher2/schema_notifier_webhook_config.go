package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierWebhookConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"url": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Webhook URL",
		},
		"proxy_url": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Webhook proxy URL",
		},
	}

	return s
}
