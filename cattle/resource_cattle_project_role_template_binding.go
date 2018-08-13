package cattle

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceCattleProjectRoleTemplateBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceCattleProjectRoleTemplateBindingCreate,
		Read:   resourceCattleProjectRoleTemplateBindingRead,
		Update: resourceCattleProjectRoleTemplateBindingUpdate,
		Delete: resourceCattleProjectRoleTemplateBindingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCattleProjectRoleTemplateBindingImport,
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

func resourceCattleProjectRoleTemplateBindingCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	//annotations := d.Get("annotations").(map[string]string)
	//labels := d.Get("labels").(map[string]string)

	log.Printf("[INFO] Creating Project Role Template Binding %s", name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectRole := &managementClient.ProjectRoleTemplateBinding{
		ProjectId:        d.Get("project_id").(string),
		RoleTemplateId:   d.Get("role_template_id").(string),
		Name:             d.Get("name").(string),
		GroupId:          d.Get("group_id").(string),
		GroupPrincipalId: d.Get("group_principal_id").(string),
		UserId:           d.Get("user_id").(string),
		UserPrincipalId:  d.Get("user_principal_id").(string),
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

	return resourceCattleProjectRoleTemplateBindingRead(d, meta)
}

func resourceCattleProjectRoleTemplateBindingRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("project_id", projectRole.ProjectId)
	d.Set("role_template_id", projectRole.RoleTemplateId)
	d.Set("name", projectRole.Name)
	d.Set("group_id", projectRole.GroupId)
	d.Set("group_principal_id", projectRole.GroupPrincipalId)
	d.Set("user_id", projectRole.UserId)
	d.Set("user_principal_id", projectRole.UserPrincipalId)
	//d.Set("annotations", projectRole.Annotations)
	//d.Set("labels", projectRole.Labels)

	return nil
}

func resourceCattleProjectRoleTemplateBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Project Role Template Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectRole, err := client.ProjectRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]string{
		"projectId":        d.Get("project_id").(string),
		"roleTemplateId":   d.Get("role_template_id").(string),
		"name":             d.Get("name").(string),
		"groupId":          d.Get("group_id").(string),
		"groupPrincipalId": d.Get("group_principal_id").(string),
		"userId":           d.Get("user_id").(string),
		"userPrincipalId":  d.Get("user_principal_id").(string),
		//"annotations": d.Get("annotations").(map[string]string),
		//"labels":      d.Get("labels").(map[string]string),
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

	return resourceCattleProjectRoleTemplateBindingRead(d, meta)
}

func resourceCattleProjectRoleTemplateBindingDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceCattleProjectRoleTemplateBindingImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	projectRole, err := client.ProjectRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(projectRole.ID)
	d.Set("project_id", projectRole.ProjectId)
	d.Set("role_template_id", projectRole.RoleTemplateId)
	d.Set("name", projectRole.Name)
	d.Set("group_id", projectRole.GroupId)
	d.Set("group_principal_id", projectRole.GroupPrincipalId)
	d.Set("user_id", projectRole.UserId)
	d.Set("user_principal_id", projectRole.UserPrincipalId)
	//d.Set("annotations", project.Annotations)
	//d.Set("labels", project.Labels)

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
