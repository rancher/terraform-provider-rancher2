package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenNotifier(d *schema.ResourceData, in *managementClient.Notifier) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("cluster_id", in.ClusterID)
	d.Set("name", in.Name)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	d.Set("send_resolved", in.SendResolved)

	if in.DingtalkConfig != nil {
		v, ok := d.Get("dingtalk_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		d.Set("dingtalk_config", flattenNotifierDingtalkConfig(in.DingtalkConfig, v))
	}

	if in.MSTeamsConfig != nil {
		v, ok := d.Get("msteams_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		d.Set("msteams_config", flattenNotifierMSTeamsConfig(in.MSTeamsConfig, v))
	}

	if in.PagerdutyConfig != nil {
		v, ok := d.Get("pagerduty_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		d.Set("pagerduty_config", flattenNotifierPagerdutyConfig(in.PagerdutyConfig, v))
	}

	if in.SlackConfig != nil {
		v, ok := d.Get("slack_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		d.Set("slack_config", flattenNotifierSlackConfig(in.SlackConfig, v))
	}

	if in.SMTPConfig != nil {
		v, ok := d.Get("smtp_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		d.Set("smtp_config", flattenNotifierSMTPConfig(in.SMTPConfig, v))
	}

	if in.WebhookConfig != nil {
		v, ok := d.Get("webhook_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		d.Set("webhook_config", flattenNotifierWebhookConfig(in.WebhookConfig, v))
	}

	if in.WechatConfig != nil {
		v, ok := d.Get("wechat_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		d.Set("wechat_config", flattenNotifierWechatConfig(in.WechatConfig, v))
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

// Expanders

func expandNotifier(in *schema.ResourceData) (*managementClient.Notifier, error) {
	obj := &managementClient.Notifier{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] Expanding notifier: Schema Resource data is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ClusterID = in.Get("cluster_id").(string)
	obj.Name = in.Get("name").(string)

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("send_resolved").(bool); ok {
		obj.SendResolved = v
	}

	if v, ok := in.Get("dingtalk_config").([]interface{}); ok && len(v) > 0 {
		obj.DingtalkConfig = expandNotifierDingtalkConfig(v)
	}

	if v, ok := in.Get("msteams_config").([]interface{}); ok && len(v) > 0 {
		obj.MSTeamsConfig = expandNotifierMSTeamsConfig(v)
	}

	if v, ok := in.Get("pagerduty_config").([]interface{}); ok && len(v) > 0 {
		obj.PagerdutyConfig = expandNotifierPagerdutyConfig(v)
	}

	if v, ok := in.Get("slack_config").([]interface{}); ok && len(v) > 0 {
		obj.SlackConfig = expandNotifierSlackConfig(v)
	}

	if v, ok := in.Get("smtp_config").([]interface{}); ok && len(v) > 0 {
		obj.SMTPConfig = expandNotifierSMTPConfig(v)
	}

	if v, ok := in.Get("webhook_config").([]interface{}); ok && len(v) > 0 {
		obj.WebhookConfig = expandNotifierWebhookConfig(v)
	}

	if v, ok := in.Get("wechat_config").([]interface{}); ok && len(v) > 0 {
		obj.WechatConfig = expandNotifierWechatConfig(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
