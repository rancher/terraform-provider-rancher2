---
layout: "rancher2"
page_title: "Rancher2: rancher2_project_logging"
sidebar_current: "docs-rancher2-resource-project_logging"
description: |-
  Provides a Rancher v2 Project Logging resource. This can be used to create Project Logging for Rancher v2 environments and retrieve their information.
---

# rancher2\_project\_logging

Provides a Rancher v2 Project Logging resource. This can be used to create Project Logging for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Project Logging
resource "rancher2_project_logging" "foo" {
  name = "foo"
  project_id = "<project_id>"
  kind = "syslog"
  syslog_config {
    endpoint = "<syslog_endpoint>"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project id to configure logging (string)
* `name` - (Required) The name of the Project Logging config (string)
* `kind` - (Required) The kind of the Project Logging. `elasticsearch`, `fluentd`, `kafka`, `splunk` and `syslog` are supported (string)
* `elasticsearch_config` - (Optional) The elasticsearch config for Project Logging. For `kind = elasticsearch`. Conflicts with `fluentd_config`, `kafka_config`, `splunk_config` and `syslog_config` (list maxitems:1)
* `fluentd_config` - (Optional) The fluentd config for Project Logging. For `kind = fluentd`. Conflicts with `elasticsearch_config`, `kafka_config`, `splunk_config` and `syslog_config` (list maxitems:1)
* `kafka_config` - (Optional) The kafka config for Project Logging. For `kind = kafka`. Conflicts with `elasticsearch_config`, `fluentd_config`, `splunk_config` and `syslog_config` (list maxitems:1)
* `namespace_id` - (Optional) The namespace id from Project logging (string)
* `output_flush_interval` - (Optional) How often buffered logs would be flushed. Default: `3` seconds (int)
* `output_tags` - (Optional/computed) The output tags for Project Logging (map)
* `splunk_config` - (Optional) The splunk config for Project Logging. For `kind = splunk`. Conflicts with `elasticsearch_config`, `fluentd_config`, `kafka_config`, and `syslog_config` (list maxitems:1)
* `syslog_config` - (Optional) The syslog config for Project Logging. For `kind = syslog`. Conflicts with `elasticsearch_config`, `fluentd_config`, `kafka_config`, and `splunk_config` (list maxitems:1)
* `annotations` - (Optional/Computed) Annotations for Project Logging object (map)
* `labels` - (Optional/Computed) Labels for Project Logging object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `elasticsearch_config`

#### Arguments

* `endpoint` - (Required) Endpoint of the elascticsearch service. Must include protocol, `http://` or `https://` (string)
* `auth_password` - (Optional/Sensitive) User password for the elascticsearch service (string)
* `auth_username` - (Optional/Sensitive) Username for the elascticsearch service (string)
* `certificate` - (Optional/Sensitive) SSL certificate for the elascticsearch service (string)
* `client_cert` - (Optional/Sensitive) SSL client certificate for the elascticsearch service (string)
* `client_key` - (Optional/Sensitive) SSL client key for the elascticsearch service (string)
* `client_key_pass` - (Optional/Sensitive) SSL client key password for the elascticsearch service (string)
* `date_format` - (Optional) Date format for the elascticsearch logs. Default: `YYYY-MM-DD` (string)
* `index_prefix` - (Optional) Index prefix for the elascticsearch logs. Default: `local` (string)
* `ssl_verify` - (Optional) SSL verify for the elascticsearch service (bool)
* `ssl_version` - (Optional) SSL version for the elascticsearch service (string)

### `fluentd_config`

#### Arguments

* `fluent_servers` - (Required) Servers for the fluentd service (list)
* `certificate` - (Optional/Sensitive) SSL certificate for the fluentd service (string)
* `compress` - (Optional) Compress data for the fluentd service (bool)
* `enable_tls` - (Optional) Enable TLS for the fluentd service (bool)

#### `fluent_servers`

##### Arguments

* `endpoint` - (Required) Endpoint of the fluentd service (string)
* `hostname` - (Optional) Hostname of the fluentd service (string)
* `password` - (Optional/Sensitive) User password of the fluentd service (string)
* `shared_key` - (Optional/Sensitive) Shared key of the fluentd service (string)
* `standby` - (Optional) Standby server of the fluentd service (bool)
* `username` - (Optional/Sensitive) Username of the fluentd service (string)
* `weight` - (Optional) Weight of the fluentd server (int)

### `kafka_config`

#### Arguments

* `topic` - (Required) Topic to publish on the kafka service (string)
* `broker_endpoints` - (Optional) Kafka endpoints for kafka service. Conflicts with `zookeeper_endpoint` (list)
* `certificate` - (Optional/Sensitive) SSL certificate for the kafka service (string)
* `client_cert` - (Optional/Sensitive) SSL client certificate for the kafka service (string)
* `client_key` - (Optional/Sensitive) SSL client key for the kafka service (string)
* `zookeeper_endpoint` - (Optional) Zookeeper endpoint for kafka service. Conflicts with `broker_endpoints` (string)

### `splunk_config`

#### Arguments

* `endpoint` - (Required) Endpoint of the splunk service. Must include protocol, `http://` or `https://` (string)
* `token` - (Required/Sensitive) Token for the splunk service (string)
* `certificate` - (Optional/Sensitive) SSL certificate for the splunk service (string)
* `client_cert` - (Optional/Sensitive) SSL client certificate for the splunk service (string)
* `client_key` - (Optional/Sensitive) SSL client key for the splunk service (string)
* `client_key_pass` - (Optional/Sensitive) SSL client key password for the splunk service (string)
* `index` - (Optional) Index prefix for the splunk logs (string)
* `source` - (Optional) Date format for the splunk logs (string)
* `ssl_verify` - (Optional) SSL verify for the splunk service (bool)

### `syslog_config`

#### Arguments

* `endpoint` - (Required) Endpoint of the syslog service (string)
* `certificate` - (Optional/Sensitive) SSL certificate for the syslog service (string)
* `client_cert` - (Optional/Sensitive) SSL client certificate for the syslog service (string)
* `client_key` - (Optional/Sensitive) SSL client key for the syslog service (string)
* `program` - (Optional) Program for the syslog service (string)
* `protocol` - (Optional) Protocol for the syslog service. `tcp` and `udp` are supported. Default: `udp` (string)
* `severity` - (Optional) Date format for the syslog logs. `emergency`, `alert`, `critical`, `error`, `warning`, `notice`, `info` and `debug` are supported. Default: `notice` (string)
* `ssl_verify` - (Optional) SSL verify for the syslog service (bool)
* `token` - (Optional/Sensitive) Token for the syslog service (string)

## Timeouts

`rancher2_project_logging` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating project logging configurations.
- `update` - (Default `10 minutes`) Used for project logging configuration modifications.
- `delete` - (Default `10 minutes`) Used for deleting project logging configurations.

## Import

Project Logging can be imported using the Rancher Project Logging ID

```
$ terraform import rancher2_project_logging.foo <project_logging_id>
```

