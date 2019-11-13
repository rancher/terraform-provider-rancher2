package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterLogging() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterLoggingRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"elasticsearch_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingElasticsearchConfigFields(),
				},
			},
			"fluentd_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingFluentdConfigFields(),
				},
			},
			"kafka_config": &schema.Schema{
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
			"output_tags": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"splunk_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingSplunkConfigFields(),
				},
			},
			"syslog_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingSyslogConfigFields(),
				},
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterLoggingRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster_id").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
	}
	listOpts := NewListOpts(filters)

	clusterLoggings, err := client.ClusterLogging.List(listOpts)
	if err != nil {
		return err
	}

	count := len(clusterLoggings.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster logging on cluster ID \"%s\" not found", clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster logging on cluster ID \"%s\"", count, clusterID)
	}

	return flattenClusterLogging(d, &clusterLoggings.Data[0])
}
