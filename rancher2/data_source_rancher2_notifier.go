package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2Notifier() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2NotifierRead,

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
				Type:     schema.TypeList,
				Computed: true,
				//ConflictsWith: []string{"dingtalk_config", "msteams_config", "smtp_config", "slack_config", "webhook_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierPagerdutyConfigFields(),
				},
			},
			"slack_config": {
				Type:     schema.TypeList,
				Computed: true,
				//ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "webhook_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierSlackConfigFields(),
				},
			},
			"smtp_config": {
				Type:     schema.TypeList,
				Computed: true,
				//ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "slack_config", "webhook_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierSMTPConfigFields(),
				},
			},
			"webhook_config": {
				Type:     schema.TypeList,
				Computed: true,
				//ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "slack_config", "wechat_config"},
				Elem: &schema.Resource{
					Schema: notifierWebhookConfigFields(),
				},
			},
			"wechat_config": {
				Type:     schema.TypeList,
				Computed: true,
				//ConflictsWith: []string{"dingtalk_config", "msteams_config", "pagerduty_config", "smtp_config", "slack_config", "webhook_config"},
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
			// TODO - ANDY I added the following fields here cause v2 was panicking. AS they didn't exists I used computed
			"send_resolved": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2NotifierRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	count := len(notifiers.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] notifier with name \"%s\" and cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d notifier with name \"%s\" and cluster ID \"%s\"", count, name, clusterID)
	}

	return diag.FromErr(flattenNotifier(d, &notifiers.Data[0]))
}
