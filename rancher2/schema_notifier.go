package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func notifierFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Notifier name",
		},
		"cluster_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Notifier cluster ID",
		},
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Notifier description",
		},
		"send_resolved": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Notifier send resolved",
			Default:     false,
		},
		"pagerduty_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"smtp_config", "slack_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierPagerdutyConfigFields(),
			},
		},
		"slack_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"pagerduty_config", "smtp_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierSlackConfigFields(),
			},
		},
		"smtp_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"pagerduty_config", "slack_config", "webhook_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierSMTPConfigFields(),
			},
		},
		"webhook_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"pagerduty_config", "smtp_config", "slack_config", "wechat_config"},
			Elem: &schema.Resource{
				Schema: notifierWebhookConfigFields(),
			},
		},
		"wechat_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"pagerduty_config", "smtp_config", "slack_config", "webhook_config"},
			Elem: &schema.Resource{
				Schema: notifierWechatConfigFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
