package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ProjectLogging() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ProjectLoggingRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_target_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingCustomTargetConfigFields(),
				},
			},
			"enable_json_parsing": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Optional enable json log parsing",
			},
			"elasticsearch_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingElasticsearchConfigFields(),
				},
			},
			"fluentd_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingFluentdConfigFields(),
				},
			},
			"kafka_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingKafkaConfigFields(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_flush_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"output_tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"splunk_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingSplunkConfigFields(),
				},
			},
			"syslog_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingSyslogConfigFields(),
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

func dataSourceRancher2ProjectLoggingRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)

	filters := map[string]interface{}{
		"projectId": projectID,
	}
	listOpts := NewListOpts(filters)

	projectLoggings, err := client.ProjectLogging.List(listOpts)
	if err != nil {
		return err
	}

	count := len(projectLoggings.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] project logging on project ID \"%s\" not found", projectID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d project logging on project ID \"%s\"", count, projectID)
	}

	return flattenProjectLogging(d, &projectLoggings.Data[0])
}
