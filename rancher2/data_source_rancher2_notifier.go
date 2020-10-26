package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2Notifier() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2NotifierRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notifier name",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notifier cluster ID",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Notifier description",
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
				Computed:      true,
				ConflictsWith: []string{"dingtalk_config", "msteams_config", "smtp_config", "slack_config", "webhook_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierPagerdutyConfigFields(),
				},
			},
			"slack_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Computed:      true,
				ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "webhook_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierSlackConfigFields(),
				},
			},
			"smtp_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Computed:      true,
				ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "slack_config", "webhook_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierSMTPConfigFields(),
				},
			},
			"webhook_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Computed:      true,
				ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "slack_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierWebhookConfigFields(),
				},
			},
			"wechat_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Computed:      true,
				ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "slack_config", "webhook_config"},
				Elem: &schema.Resource{
					Schema: notifierWechatConfigFields(),
				},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2NotifierRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	clusterID := d.Get("cluster_id").(string)

	filters := map[string]interface{}{
		"name":      name,
		"clusterId": clusterID,
	}
	listOpts := NewListOpts(filters)

	notifiers, err := client.Notifier.List(listOpts)
	if err != nil {
		return err
	}

	count := len(notifiers.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] notifier with name \"%s\" and cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d notifier with name \"%s\" and cluster ID \"%s\"", count, name, clusterID)
	}

	return flattenNotifier(d, &notifiers.Data[0])
}
