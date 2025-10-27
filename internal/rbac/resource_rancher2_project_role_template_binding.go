package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ProjectRoleTemplateBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ProjectRoleTemplateBindingCreate,
		Read:   resourceRancher2ProjectRoleTemplateBindingRead,
		Update: resourceRancher2ProjectRoleTemplateBindingUpdate,
		Delete: resourceRancher2ProjectRoleTemplateBindingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ProjectRoleTemplateBindingImport,
		},

		Schema: projectRoleTemplateBindingFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ProjectRoleTemplateBindingCreate(d *schema.ResourceData, meta interface{}) error {
	projectRole := expandProjectRoleTemplateBinding(d)

	err := meta.(*Config).ProjectExist(projectRole.ProjectID)
	if err != nil {
		return err
	}

	err = meta.(*Config).RoleTemplateExist(projectRole.RoleTemplateID)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Project Role Template Binding %s", projectRole.Name)

	newProjectRole, err := client.ProjectRoleTemplateBinding.Create(projectRole)
	if err != nil {
		return err
	}

	d.SetId(newProjectRole.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    projectRoleTemplateBindingStateRefreshFunc(client, newProjectRole.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project role template binding (%s) to be created: %s", newProjectRole.ID, waitErr)
	}

	return resourceRancher2ProjectRoleTemplateBindingRead(d, meta)
}

func resourceRancher2ProjectRoleTemplateBindingRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Project Role Template Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectRole, err := client.ProjectRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Project Role Template Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenProjectRoleTemplateBinding(d, projectRole)
}

func resourceRancher2ProjectRoleTemplateBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Project Role Template Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectRole, err := client.ProjectRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"groupId":          d.Get("group_id").(string),
		"groupPrincipalId": d.Get("group_principal_id").(string),
		"roleTemplateId":   d.Get("role_template_id").(string),
		"userId":           d.Get("user_id").(string),
		"userPrincipalId":  d.Get("user_principal_id").(string),
		"annotations":      toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":           toMapString(d.Get("labels").(map[string]interface{})),
	}

	newProjectRole, err := client.ProjectRoleTemplateBinding.Update(projectRole, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    projectRoleTemplateBindingStateRefreshFunc(client, newProjectRole.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project role template binding (%s) to be updated: %s", newProjectRole.ID, waitErr)
	}

	return resourceRancher2ProjectRoleTemplateBindingRead(d, meta)
}

func resourceRancher2ProjectRoleTemplateBindingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Project Role Template Binding ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectRole, err := client.ProjectRoleTemplateBinding.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Project Role Template Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ProjectRoleTemplateBinding.Delete(projectRole)
	if err != nil {
		return fmt.Errorf("Error removing Project Role Template Binding: %s", err)
	}

	log.Printf("[DEBUG] Waiting for project role template binding (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    projectRoleTemplateBindingStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project role template binding (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// PpojectRoleTemplateBindingStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Project Role Template Binding.
func projectRoleTemplateBindingStateRefreshFunc(client *managementClient.Client, projectRoleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ProjectRoleTemplateBinding.ByID(projectRoleID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}
