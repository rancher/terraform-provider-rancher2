package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
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

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_principal_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_principal_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRancher2ProjectRoleTemplateBindingCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	roleTemplateID := d.Get("role_template_id").(string)

	log.Printf("[INFO] Creating Project Role Template Binding %s", name)

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		return err
	}

	err = meta.(*Config).RoleTemplateExist(roleTemplateID)
	if err != nil {
		return err
	}

	projectRole := &managementClient.ProjectRoleTemplateBinding{
		ProjectID:        projectID,
		RoleTemplateID:   roleTemplateID,
		Name:             d.Get("name").(string),
		GroupID:          d.Get("group_id").(string),
		GroupPrincipalID: d.Get("group_principal_id").(string),
		UserID:           d.Get("user_id").(string),
		UserPrincipalID:  d.Get("user_principal_id").(string),
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newProjectRole, err := client.ProjectRoleTemplateBinding.Create(projectRole)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    projectRoleTemplateBindingStateRefreshFunc(client, newProjectRole.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for project role template binding (%s) to be created: %s", newProjectRole.ID, waitErr)
	}

	d.SetId(newProjectRole.ID)

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
		if IsNotFound(err) {
			log.Printf("[INFO] Project Role Template Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_id", projectRole.ProjectID)
	d.Set("role_template_id", projectRole.RoleTemplateID)
	d.Set("name", projectRole.Name)
	d.Set("group_id", projectRole.GroupID)
	d.Set("group_principal_id", projectRole.GroupPrincipalID)
	d.Set("user_id", projectRole.UserID)
	d.Set("user_principal_id", projectRole.UserPrincipalID)

	return nil
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
		"projectId":        d.Get("project_id").(string),
		"roleTemplateId":   d.Get("role_template_id").(string),
		"groupId":          d.Get("group_id").(string),
		"groupPrincipalId": d.Get("group_principal_id").(string),
		"userId":           d.Get("user_id").(string),
		"userPrincipalId":  d.Get("user_principal_id").(string),
	}

	if projectRole.ProjectID != update["projectId"].(string) {
		err = meta.(*Config).ProjectExist(update["projectId"].(string))
		if err != nil {
			return err
		}
	}

	if projectRole.RoleTemplateID != update["roleTemplateId"].(string) {
		err = meta.(*Config).RoleTemplateExist(update["roleTemplateId"].(string))
		if err != nil {
			return err
		}
	}

	newProjectRole, err := client.ProjectRoleTemplateBinding.Update(projectRole, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    projectRoleTemplateBindingStateRefreshFunc(client, newProjectRole.ID),
		Timeout:    10 * time.Minute,
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
		if IsNotFound(err) {
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
		Timeout:    10 * time.Minute,
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

func resourceRancher2ProjectRoleTemplateBindingImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	projectRole, err := client.ProjectRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(projectRole.ID)
	d.Set("project_id", projectRole.ProjectID)
	d.Set("role_template_id", projectRole.RoleTemplateID)
	d.Set("name", projectRole.Name)
	d.Set("group_id", projectRole.GroupID)
	d.Set("group_principal_id", projectRole.GroupPrincipalID)
	d.Set("user_id", projectRole.UserID)
	d.Set("user_principal_id", projectRole.UserPrincipalID)

	return []*schema.ResourceData{d}, nil
}

// ProjectRoleTemplateBindingStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Project Role Template Binding.
func projectRoleTemplateBindingStateRefreshFunc(client *managementClient.Client, projectRoleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pro, err := client.ProjectRoleTemplateBinding.ByID(projectRoleID)
		if err != nil {
			if IsNotFound(err) {
				return pro, "removed", nil
			}
			return nil, "", err
		}

		return pro, "active", nil
	}
}
