package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2ProjectLogging() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ProjectLoggingCreate,
		Read:   resourceRancher2ProjectLoggingRead,
		Update: resourceRancher2ProjectLoggingUpdate,
		Delete: resourceRancher2ProjectLoggingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ProjectLoggingImport,
		},

		Schema: projectLoggingFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ProjectLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	projectLogging, err := expandProjectLogging(d)
	if err != nil {
		return err
	}

	err = meta.(*Config).ProjectExist(projectLogging.ProjectID)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Project Logging %s", projectLogging.Name)

	newProjectLogging, err := client.ProjectLogging.Create(projectLogging)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"provisioning"},
		Target:     []string{"active"},
		Refresh:    projectLoggingStateRefreshFunc(client, newProjectLogging.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project logging (%s) to be created: %s", newProjectLogging.ID, waitErr)
	}

	err = flattenProjectLogging(d, newProjectLogging)
	if err != nil {
		return err
	}

	return resourceRancher2ProjectLoggingRead(d, meta)
}

func resourceRancher2ProjectLoggingRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Project Logging ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectLogging, err := client.ProjectLogging.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Project Logging ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenProjectLogging(d, projectLogging)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2ProjectLoggingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Project Logging ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectLogging, err := client.ProjectLogging.ByID(d.Id())
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

	newProjectLogging, err := client.ProjectLogging.Update(projectLogging, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"active"},
		Refresh:    projectLoggingStateRefreshFunc(client, newProjectLogging.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project Logging (%s) to be updated: %s", newProjectLogging.ID, waitErr)
	}

	return resourceRancher2ProjectLoggingRead(d, meta)
}

func resourceRancher2ProjectLoggingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Project Logging ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectLogging, err := client.ProjectLogging.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Project Logging ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ProjectLogging.Delete(projectLogging)
	if err != nil {
		return fmt.Errorf("Error removing Project Logging: %s", err)
	}

	log.Printf("[DEBUG] Waiting for project logging (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    projectLoggingStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project logging (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// projectLoggingStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster Role Template Binding.
func projectLoggingStateRefreshFunc(client *managementClient.Client, projectLoggingID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ProjectLogging.ByID(projectLoggingID)
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
