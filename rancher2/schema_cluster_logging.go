package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	clusterLoggingKinds = []string{loggingElasticsearchKind, loggingFluentdKind, loggingKafkaKind, loggingSplunkKind, loggingSyslogKind}
)

// Shemas

func clusterLoggingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"kind": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(clusterLoggingKinds, true),
		},
		"elasticsearch_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"fluentd_config", "kafka_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingElasticsearchConfigFields(),
			},
		},
		"fluentd_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "kafka_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingFluentdConfigFields(),
			},
		},
		"kafka_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "fluentd_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingKafkaConfigFields(),
			},
		},
		"namespace_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"output_flush_interval": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  3,
		},
		"output_tags": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"splunk_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "fluentd_config", "kafka_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingSplunkConfigFields(),
			},
		},
		"syslog_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "fluentd_config", "kafka_config", "splunk_config"},
			Elem: &schema.Resource{
				Schema: loggingSyslogConfigFields(),
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
