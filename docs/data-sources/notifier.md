---
page_title: "rancher2_notifier Data Source"
---

# rancher2\_notifier Data Source

Use this data source to retrieve information about a Rancher v2 notifier.

## Example Usage

```
data "rancher2_notifier" "foo" {
    name = "foo"
    cluster_id = "<cluster_id>"
}
```

## Argument Reference

* `name` - (Required) The name of the notifier (string)
* `cluster_id` - (Required) The cluster id where create notifier (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `description` - (Computed) The notifier description (string)
* `send_resolved` - (Computed) If the notifier sends resolved notifications (bool)
* `dingtalk_config` - (Computed) Dingtalk config for notifier (list maxitems:1)
* `msteams_config` - (Computed) MSTeams config for notifier (list maxitems:1)
* `pagerduty_config` - (Computed) Pagerduty config for notifier (list maxitems:1)
* `slack_config` - (Computed) Slack config for notifier (list maxitems:1)
* `smtp_config` - (Computed) SMTP config for notifier (list maxitems:1)
* `webhook_config` - (Computed) Webhook config for notifier (list maxitems:1)
* `wechat_config` - (Computed) Wechat config for notifier (list maxitems:1)
* `annotations` - (Computed) Annotations for notifier object (map)
* `labels` - (Computed) Labels for notifier object (map)
