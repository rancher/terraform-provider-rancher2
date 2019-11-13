package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

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
		elkConfig, err := expandLoggingElasticsearchConfig(d.Get("elasticsearch_config").([]interface{}))
		if err != nil {
			return err
		}
		update["elasticsearchConfig"] = elkConfig
	case loggingFluentdKind:
		fluentdConfig, err := expandLoggingFluentdConfig(d.Get("fluentd_config").([]interface{}))
		if err != nil {
			return err
		}
		update["fluentForwarderConfig"] = fluentdConfig
	case loggingKafkaKind:
		kafkaConfig, err := expandLoggingKafkaConfig(d.Get("kafka_config").([]interface{}))
		if err != nil {
			return err
		}
		update["kafkaConfig"] = kafkaConfig
	case loggingSplunkKind:
		splunkConfig, err := expandLoggingSplunkConfig(d.Get("splunk_config").([]interface{}))
		if err != nil {
			return err
		}
		update["splunkConfig"] = splunkConfig
	case loggingSyslogKind:
		syslogConfig, err := expandLoggingSyslogConfig(d.Get("syslog_config").([]interface{}))
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
