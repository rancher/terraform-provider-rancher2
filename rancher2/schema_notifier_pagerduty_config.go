package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierPagerdutyConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"service_key": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Pagerduty service key",
		},
		"proxy_url": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Pagerduty proxy URL",
		},
	}

	return s
}
