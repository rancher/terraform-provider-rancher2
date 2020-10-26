---
page_title: "rancher2_notifier Resource"
---

# rancher2\_notifier Resource

Provides a Rancher v2 Notifier resource. This can be used to create notifiers for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Notifier
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "<cluster_id>"
  description = "Terraform notifier acceptance test"
  send_resolved = "true"
  pagerduty_config {
    service_key = "XXXXXXXX"
    proxy_url = "http://proxy.test.io"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the notifier (string)
* `cluster_id` - (Required/ForceNew) The cluster id where create notifier (string)
* `description` - (Optional) The notifier description (string)
* `send_resolved` = (Optional) Enable the notifier to send resolved notifications. Default `false` (bool)
* `dingtalk_config` - (Optional) Dingtalk config for notifier (list maxitems:1)
* `msteams_config` - (Optional) MSTeams config for notifier (list maxitems:1)
* `pagerduty_config` - (Optional) Pagerduty config for notifier (list maxitems:1)
* `slack_config` - (Optional) Slack config for notifier (list maxitems:1)
* `smtp_config` - (Optional) SMTP config for notifier (list maxitems:1)
* `webhook_config` - (Optional) Webhook config for notifier (list maxitems:1)
* `wechat_config` - (Optional) Wechat config for notifier (list maxitems:1)
* `annotations` - (Optional/Computed) Annotations for notifier object (map)
* `labels` - (Optional/Computed) Labels for notifier object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `dingtalk_config`

#### Arguments

* `url` - (Required) Dingtalk url (string)
* `proxy_url` - (Optional) Dingtalk proxy url (string)
* `secret` - (Optional) Secret for url sign enable (string)

### `msteams_config`

#### Arguments

* `url` - (Required) MSTeams url (string)
* `proxy_url` - (Optional) MSTeams proxy url (string)

### `pagerduty_config`

#### Arguments

* `service_key` - (Required) Pagerduty service key (string)
* `proxy_url` - (Optional) Pagerduty proxy url (string)

### `slack_config`

#### Arguments

* `default_recipient` - (Required) Slack default recipient (string)
* `url` - (Required) Slack url (string)
* `proxy_url` - (Optional) Slack proxy url (string)

### `smtp_config`

* `default_recipient` - (Required) SMTP default recipient (string)
* `host` - (Required) SMTP host (string)
* `port` - (Required) SMTP port (int)
* `sender` - (Required) SMTP sender (string)
* `password` - (Optional/Sensitive) SMTP password (string)
* `tls` - (Optional/Sensitive) SMTP tls. Default `true` (bool)
* `username` - (Optional/Sensitive) SMTP username (string)

#### Arguments

### `webhook_config`

* `url` - (Required) Webhook url (string)
* `proxy_url` - (Optional) Webhook proxy url (string)

#### Arguments

### `wechat_config`

#### Arguments

* `agent` - (Required) Wechat agent ID (string)
* `corp` - (Required) Wechat corporation ID (string)
* `default_recipient` - (Required) Wechat default recipient (string)
* `secret` - (Required/Sensitive) Wechat agent ID (string)
* `proxy_url` - (Optional) Wechat proxy url (string)
* `recipient_type` - (Optional) Wechat recipient type. Allowed values: `party` | `tag` | `user` (string)

## Timeouts

`rancher2_notifier` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating notifiers.
- `update` - (Default `10 minutes`) Used for notifier modifications.
- `delete` - (Default `10 minutes`) Used for deleting notifiers.

## Import

Notifiers can be imported using the Rancher nNtifier ID

```
$ terraform import rancher2_notifier.foo &lt;notifier_id&gt;
```

