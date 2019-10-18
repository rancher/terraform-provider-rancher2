package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierSlackConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default_recipient": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Slack default channel",
		},
		"url": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Slack URL",
		},
		"proxy_url": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Slack proxy URL",
		},
	}

	return s
}
