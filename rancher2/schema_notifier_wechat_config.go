package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	notifierWechatRecipientTypeParty = "party"
	notifierWechatRecipientTypeTag   = "tag"
	notifierWechatRecipientTypeUser  = "user"
)

var (
	notifierWechatRecipientTypes = []string{notifierWechatRecipientTypeParty, notifierWechatRecipientTypeTag, notifierWechatRecipientTypeUser}
)

// Schemas

func notifierWechatConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"agent": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Wechat application agent ID",
		},
		"corp": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Wechat corporation ID",
		},
		"default_recipient": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Wechat default channel",
		},
		"secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Wechat application secret",
		},
		"proxy_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Wechat proxy URL",
		},
		"recipient_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      notifierWechatRecipientTypeParty,
			Description:  "Wechat recipient type",
			ValidateFunc: validation.StringInSlice(notifierWechatRecipientTypes, true),
		},
	}

	return s
}
