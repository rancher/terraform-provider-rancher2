package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	recipientTypePagerduty = "pagerduty"
	recipientTypeSlack     = "slack"
	recipientTypeSMTP      = "email"
	recipientTypeWebhook   = "webhook"
	recipientTypeWechat    = "wechat"
)

//Schemas

func recipientFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"notifier_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Recipient notifier ID",
		},
		"notifier_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Recipient notifier type",
		},
		"recipient": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Recipient",
		},
	}
	return s
}
