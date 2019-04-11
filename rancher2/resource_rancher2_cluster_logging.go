package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
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
				Schema: elasticsearchConfigFields(),
			},
		},
		"fluentd_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "kafka_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: fluentdConfigFields(),
			},
		},
		"kafka_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "fluentd_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: kafkaConfigFields(),
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
				Schema: splunkConfigFields(),
			},
		},
		"syslog_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "fluentd_config", "kafka_config", "splunk_config"},
			Elem: &schema.Resource{
				Schema: syslogConfigFields(),
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

// Flatteners

func flattenClusterLogging(d *schema.ResourceData, in *managementClient.ClusterLogging) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("cluster_id", in.ClusterID)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	kind := d.Get("kind").(string)
	if kind == "" {
		if in.ElasticsearchConfig != nil {
			kind = loggingElasticsearchKind
		}
		if in.FluentForwarderConfig != nil {
			kind = loggingFluentdKind
		}
		if in.KafkaConfig != nil {
			kind = loggingKafkaKind
		}
		if in.SplunkConfig != nil {
			kind = loggingSplunkKind
		}
		if in.SyslogConfig != nil {
			kind = loggingSyslogKind
		}

		err = d.Set("kind", kind)
		if err != nil {
			return err
		}
	}

	switch kind {
	case loggingElasticsearchKind:
		elkConfig, err := flattenElasticsearchConfig(in.ElasticsearchConfig)
		if err != nil {
			return err
		}
		err = d.Set("elasticsearch_config", elkConfig)
		if err != nil {
			return err
		}
	case loggingFluentdKind:
		fluentdConfig, err := flattenFluentdConfig(in.FluentForwarderConfig)
		if err != nil {
			return err
		}
		err = d.Set("fluentd_config", fluentdConfig)
		if err != nil {
			return err
		}
	case loggingKafkaKind:
		kafkaConfig, err := flattenKafkaConfig(in.KafkaConfig)
		if err != nil {
			return err
		}
		err = d.Set("kafka_config", kafkaConfig)
		if err != nil {
			return err
		}
	case loggingSplunkKind:
		splunkConfig, err := flattenSplunkConfig(in.SplunkConfig)
		if err != nil {
			return err
		}
		err = d.Set("splunk_config", splunkConfig)
		if err != nil {
			return err
		}
	case loggingSyslogKind:
		syslogConfig, err := flattenSyslogConfig(in.SyslogConfig)
		if err != nil {
			return err
		}
		err = d.Set("syslog_config", syslogConfig)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("[ERROR] Flattening Cluster Logging resource data: Kind %s not supported", kind)
	}

	err = d.Set("namespace_id", in.NamespaceId)
	if err != nil {
		return err
	}

	err = d.Set("output_flush_interval", int(in.OutputFlushInterval))
	if err != nil {
		return err
	}

	err = d.Set("output_tags", toMapInterface(in.OutputTags))
	if err != nil {
		return err
	}

	err = d.Set("annotations", toMapInterface(in.Annotations))
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

func expandClusterLogging(in *schema.ResourceData) (*managementClient.ClusterLogging, error) {
	obj := &managementClient.ClusterLogging{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] Expanding cluster config: Schema Resource data is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ClusterID = in.Get("cluster_id").(string)
	obj.Name = in.Get("name").(string)

	switch kind := in.Get("kind").(string); kind {
	case loggingElasticsearchKind:
		elkConfig, err := expandElasticsearchConfig(in.Get("elasticsearch_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.ElasticsearchConfig = elkConfig
	case loggingFluentdKind:
		fluentdConfig, err := expandFluentdConfig(in.Get("fluentd_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.FluentForwarderConfig = fluentdConfig
	case loggingKafkaKind:
		kafkaConfig, err := expandKafkaConfig(in.Get("kafka_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.KafkaConfig = kafkaConfig
	case loggingSplunkKind:
		splunkConfig, err := expandSplunkConfig(in.Get("splunk_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.SplunkConfig = splunkConfig
	case loggingSyslogKind:
		syslogConfig, err := expandSyslogConfig(in.Get("syslog_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.SyslogConfig = syslogConfig
	}

	if v, ok := in.Get("namespace_id").(string); ok && len(v) > 0 {
		obj.NamespaceId = v
	}

	if v, ok := in.Get("output_flush_interval").(int); ok && v > 0 {
		obj.OutputFlushInterval = int64(v)
	}

	if v, ok := in.Get("output_tags").(map[string]interface{}); ok && len(v) > 0 {
		obj.OutputTags = toMapString(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func resourceRancher2ClusterLogging() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterLoggingCreate,
		Read:   resourceRancher2ClusterLoggingRead,
		Update: resourceRancher2ClusterLoggingUpdate,
		Delete: resourceRancher2ClusterLoggingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterLoggingImport,
		},

		Schema: clusterLoggingFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ClusterLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	clusterLogging, err := expandClusterLogging(d)
	if err != nil {
		return err
	}

	err = meta.(*Config).ClusterExist(clusterLogging.ClusterID)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Cluster Logging %s", clusterLogging.Name)

	newClusterLogging, err := client.ClusterLogging.Create(clusterLogging)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"provisioning"},
		Target:     []string{"active"},
		Refresh:    clusterLoggingStateRefreshFunc(client, newClusterLogging.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster logging (%s) to be created: %s", newClusterLogging.ID, waitErr)
	}

	err = flattenClusterLogging(d, newClusterLogging)
	if err != nil {
		return err
	}

	return resourceRancher2ClusterLoggingRead(d, meta)
}

func resourceRancher2ClusterLoggingRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster Logging ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterLogging, err := client.ClusterLogging.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster Logging ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenClusterLogging(d, clusterLogging)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2ClusterLoggingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cluster Logging ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterLogging, err := client.ClusterLogging.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":                d.Get("name").(string),
		"namespaceId":         d.Get("namespace_id").(string),
		"outputFlushInterval": int64(d.Get("output_flush_interval").(int)),
		"outputTags":          toMapString(d.Get("output_tags").(map[string]interface{})),
		"annotations":         toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":              toMapString(d.Get("labels").(map[string]interface{})),
	}

	switch kind := d.Get("kind").(string); kind {
	case loggingElasticsearchKind:
		elkConfig, err := expandElasticsearchConfig(d.Get("elasticsearch_config").([]interface{}))
		if err != nil {
			return err
		}
		update["elasticsearchConfig"] = elkConfig
	case loggingFluentdKind:
		fluentdConfig, err := expandFluentdConfig(d.Get("fluentd_config").([]interface{}))
		if err != nil {
			return err
		}
		update["fluentForwarderConfig"] = fluentdConfig
	case loggingKafkaKind:
		kafkaConfig, err := expandKafkaConfig(d.Get("kafka_config").([]interface{}))
		if err != nil {
			return err
		}
		update["kafkaConfig"] = kafkaConfig
	case loggingSplunkKind:
		splunkConfig, err := expandSplunkConfig(d.Get("splunk_config").([]interface{}))
		if err != nil {
			return err
		}
		update["splunkConfig"] = splunkConfig
	case loggingSyslogKind:
		syslogConfig, err := expandSyslogConfig(d.Get("syslog_config").([]interface{}))
		if err != nil {
			return err
		}
		update["syslogConfig"] = syslogConfig
	}

	newClusterLogging, err := client.ClusterLogging.Update(clusterLogging, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"active"},
		Refresh:    clusterLoggingStateRefreshFunc(client, newClusterLogging.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster Logging (%s) to be updated: %s", newClusterLogging.ID, waitErr)
	}

	return resourceRancher2ClusterLoggingRead(d, meta)
}

func resourceRancher2ClusterLoggingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster Logging ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterLogging, err := client.ClusterLogging.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster Logging ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ClusterLogging.Delete(clusterLogging)
	if err != nil {
		return fmt.Errorf("Error removing Cluster Logging: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cluster logging (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    clusterLoggingStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster logging (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// clusterLoggingStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster Role Template Binding.
func clusterLoggingStateRefreshFunc(client *managementClient.Client, clusterLoggingID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ClusterLogging.ByID(clusterLoggingID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, "active", nil
	}
}
