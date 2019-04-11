package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func projectFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"resource_quota": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: projectResourceQuotaFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["annotations"],
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["labels"],
		},
	}

	return s
}

// Flatteners

func flattenProject(d *schema.ResourceData, in *managementClient.Project) error {
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

	err = d.Set("description", in.Description)
	if err != nil {
		return err
	}

	if in.ResourceQuota != nil && in.NamespaceDefaultResourceQuota != nil {
		resourceQuota, err := flattenProjectResourceQuota(in.ResourceQuota, in.NamespaceDefaultResourceQuota)
		if err != nil {
			return err
		}
		err = d.Set("resource_quota", resourceQuota)
		if err != nil {
			return err
		}
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

func expandProject(in *schema.ResourceData) (*managementClient.Project, error) {
	obj := &managementClient.Project{}
	if in == nil {
		return nil, nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ClusterID = in.Get("cluster_id").(string)
	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	if v, ok := in.Get("resource_quota").([]interface{}); ok && len(v) > 0 {
		resourceQuota, nsResourceQuota, err := expandProjectResourceQuota(v)
		if err != nil {
			return obj, err
		}
		obj.ResourceQuota = resourceQuota
		obj.NamespaceDefaultResourceQuota = nsResourceQuota
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

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

	project, err := expandProject(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Project %s", project.Name)

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

	err = flattenProject(d, newProject)
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

	err = flattenProject(d, project)
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

	resourceQuota, nsResourceQuota, err := expandProjectResourceQuota(d.Get("resource_quota").([]interface{}))
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":                          d.Get("name").(string),
		"description":                   d.Get("description").(string),
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
