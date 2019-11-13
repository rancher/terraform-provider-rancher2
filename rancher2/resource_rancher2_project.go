package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2Project() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ProjectCreate,
		Read:   resourceRancher2ProjectRead,
		Update: resourceRancher2ProjectUpdate,
		Delete: resourceRancher2ProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ProjectImport,
		},

		Schema: projectFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	project := expandProject(d)

	active, _, err := meta.(*Config).isClusterActive(project.ClusterID)
	if err != nil {
		return err
	}
	if !active {
		if v, ok := d.Get("wait_for_cluster").(bool); ok && !v {
			return fmt.Errorf("[ERROR] Creating Project: Cluster ID %s is not active", project.ClusterID)
		}

		mgmtClient, err := meta.(*Config).ManagementClient()
		if err != nil {
			return err
		}

		stateCluster := &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"active"},
			Refresh:    clusterStateRefreshFunc(mgmtClient, project.ClusterID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitClusterErr := stateCluster.WaitForState()
		if waitClusterErr != nil {
			return fmt.Errorf("[ERROR] waiting for cluster ID (%s) to be active: %s", project.ClusterID, waitClusterErr)
		}
	}

	log.Printf("[INFO] Creating Project %s on Cluster ID %s", project.Name, project.ClusterID)

	newProject, err := client.Project.Create(project)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"initializing", "configuring", "active"},
		Target:     []string{"active"},
		Refresh:    projectStateRefreshFunc(client, newProject.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project (%s) to be created: %s", newProject.ID, waitErr)
	}

	monitoringInput := expandMonitoringInput(d.Get("project_monitoring_input").([]interface{}))
	if newProject.EnableProjectMonitoring && monitoringInput != nil {
		err = client.Project.ActionEditMonitoring(newProject, monitoringInput)
		if err != nil {
			return err
		}
	}

	err = flattenProject(d, newProject, monitoringInput)
	if err != nil {
		return err
	}

	return resourceRancher2ProjectRead(d, meta)
}

func resourceRancher2ProjectRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Project ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	project, err := client.Project.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Project ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	monitoringInput := &managementClient.MonitoringInput{}
	if project.EnableProjectMonitoring {
		monitoringOutput, err := client.Project.ActionViewMonitoring(project)
		if err != nil {
			return err
		}

		if monitoringOutput != nil && len(monitoringOutput.Answers) > 0 {
			monitoringInput.Answers = monitoringOutput.Answers
		}
	}

	err = flattenProject(d, project, monitoringInput)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2ProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Project ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	project, err := client.Project.ByID(d.Id())
	if err != nil {
		return err
	}

	resourceQuota, nsResourceQuota := expandProjectResourceQuota(d.Get("resource_quota").([]interface{}))

	update := map[string]interface{}{
		"name":                          d.Get("name").(string),
		"description":                   d.Get("description").(string),
		"containerDefaultResourceLimit": expandProjectContainerResourceLimit(d.Get("container_resource_limit").([]interface{})),
		"enableProjectMonitoring":       d.Get("enable_project_monitoring").(bool),
		"namespaceDefaultResourceQuota": nsResourceQuota,
		"resourceQuota":                 resourceQuota,
		"annotations":                   toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                        toMapString(d.Get("labels").(map[string]interface{})),
	}

	newProject, err := client.Project.Update(project, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    projectStateRefreshFunc(client, newProject.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project (%s) to be updated: %s", newProject.ID, waitErr)
	}

	if d.HasChange("pod_security_policy_template_id") {
		pspInput := &managementClient.SetPodSecurityPolicyTemplateInput{
			PodSecurityPolicyTemplateName: d.Get("pod_security_policy_template_id").(string),
		}
		_, err = client.Project.ActionSetpodsecuritypolicytemplate(newProject, pspInput)
		if err != nil {
			// Checking error due to ActionSetpodsecuritypolicytemplate() issue
			if error.Error(err) != "unexpected end of JSON input" {
				return err
			}
		}
	}

	if d.HasChange("enable_project_monitoring") || d.HasChange("project_monitoring_input") {
		monitoringInput := expandMonitoringInput(d.Get("project_monitoring_input").([]interface{}))
		if newProject.EnableProjectMonitoring && monitoringInput != nil {
			err = client.Project.ActionEditMonitoring(newProject, monitoringInput)
			if err != nil {
				return err
			}
		}
	}

	return resourceRancher2ProjectRead(d, meta)
}

func resourceRancher2ProjectDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Project ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	project, err := client.Project.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Project ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Project.Delete(project)
	if err != nil {
		return fmt.Errorf("Error removing Project: %s", err)
	}

	log.Printf("[DEBUG] Waiting for project (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    projectStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// projectStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Project.
func projectStateRefreshFunc(client *managementClient.Client, projectID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Project.ByID(projectID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, obj.State, nil
	}
}
