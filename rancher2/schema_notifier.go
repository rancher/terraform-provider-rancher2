package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Notifier name",
		},
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Notifier cluster ID",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Notifier description",
		},
		"send_resolved": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Notifier send resolved",
			Default:     false,
		},
		"dingtalk_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"msteams_config", "pagerduty_config", "smtp_config", "slack_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierDingtalkConfigFields(),
			},
		},
		"msteams_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"dingtalk_config", "pagerduty_config", "smtp_config", "slack_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierMSTeamsConfigFields(),
			},
		},
		"pagerduty_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"dingtalk_config", "msteams_config", "smtp_config", "slack_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierPagerdutyConfigFields(),
			},
		},
		"slack_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierSlackConfigFields(),
			},
		},
		"smtp_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "slack_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierSMTPConfigFields(),
			},
		},
		"webhook_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "slack_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierWebhookConfigFields(),
			},
		},
		"wechat_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "slack_config", "webhook_config"},
			Elem: &schema.Resource{
				Schema: notifierWechatConfigFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
