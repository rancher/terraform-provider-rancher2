package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2ClusterRoleTemplateBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterRoleTemplateBindingCreate,
		Read:   resourceRancher2ClusterRoleTemplateBindingRead,
		Update: resourceRancher2ClusterRoleTemplateBindingUpdate,
		Delete: resourceRancher2ClusterRoleTemplateBindingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterRoleTemplateBindingImport,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
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

func resourceRancher2ClusterRoleTemplateBindingCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	clusterID := d.Get("cluster_id").(string)
	roleTemplateID := d.Get("role_template_id").(string)

	log.Printf("[INFO] Creating Cluster Role Template Binding %s", name)

	err := meta.(*Config).ClusterExist(clusterID)
	if err != nil {
		return err
	}

	err = meta.(*Config).RoleTemplateExist(roleTemplateID)
	if err != nil {
		return err
	}

	clusterRole := &managementClient.ClusterRoleTemplateBinding{
		ClusterID:        clusterID,
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

	newClusterRole, err := client.ClusterRoleTemplateBinding.Create(clusterRole)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    clusterRoleTemplateBindingStateRefreshFunc(client, newClusterRole.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster role template binding (%s) to be created: %s", newClusterRole.ID, waitErr)
	}

	d.SetId(newClusterRole.ID)

	return resourceRancher2ClusterRoleTemplateBindingRead(d, meta)
}

func resourceRancher2ClusterRoleTemplateBindingRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster Role Template Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterRole, err := client.ClusterRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster Role Template Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cluster_id", clusterRole.ClusterID)
	d.Set("role_template_id", clusterRole.RoleTemplateID)
	d.Set("name", clusterRole.Name)
	d.Set("group_id", clusterRole.GroupID)
	d.Set("group_principal_id", clusterRole.GroupPrincipalID)
	d.Set("user_id", clusterRole.UserID)
	d.Set("user_principal_id", clusterRole.UserPrincipalID)

	return nil
}

func resourceRancher2ClusterRoleTemplateBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cluster Role Template Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterRole, err := client.ClusterRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"clusterId":        d.Get("cluster_id").(string),
		"roleTemplateId":   d.Get("role_template_id").(string),
		"groupId":          d.Get("group_id").(string),
		"groupPrincipalId": d.Get("group_principal_id").(string),
		"userId":           d.Get("user_id").(string),
		"userPrincipalId":  d.Get("user_principal_id").(string),
	}

	if clusterRole.ClusterID != update["clusterId"].(string) {
		err = meta.(*Config).ClusterExist(update["clusterId"].(string))
		if err != nil {
			return err
		}
	}

	if clusterRole.RoleTemplateID != update["roleTemplateId"].(string) {
		err = meta.(*Config).RoleTemplateExist(update["roleTemplateId"].(string))
		if err != nil {
			return err
		}
	}

	newClusterRole, err := client.ClusterRoleTemplateBinding.Update(clusterRole, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    clusterRoleTemplateBindingStateRefreshFunc(client, newClusterRole.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster role template binding (%s) to be updated: %s", newClusterRole.ID, waitErr)
	}

	return resourceRancher2ClusterRoleTemplateBindingRead(d, meta)
}

func resourceRancher2ClusterRoleTemplateBindingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster Role Template Binding ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterRole, err := client.ClusterRoleTemplateBinding.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster Role Template Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ClusterRoleTemplateBinding.Delete(clusterRole)
	if err != nil {
		return fmt.Errorf("Error removing Cluster Role Template Binding: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cluster role template binding (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    clusterRoleTemplateBindingStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster role template binding (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2ClusterRoleTemplateBindingImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	clusterRole, err := client.ClusterRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(clusterRole.ID)
	d.Set("cluster_id", clusterRole.ClusterID)
	d.Set("role_template_id", clusterRole.RoleTemplateID)
	d.Set("name", clusterRole.Name)
	d.Set("group_id", clusterRole.GroupID)
	d.Set("group_principal_id", clusterRole.GroupPrincipalID)
	d.Set("user_id", clusterRole.UserID)
	d.Set("user_principal_id", clusterRole.UserPrincipalID)

	return []*schema.ResourceData{d}, nil
}

// ClusterRoleTemplateBindingStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster Role Template Binding.
func clusterRoleTemplateBindingStateRefreshFunc(client *managementClient.Client, clusterRoleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clu, err := client.ClusterRoleTemplateBinding.ByID(clusterRoleID)
		if err != nil {
			if IsNotFound(err) {
				return clu, "removed", nil
			}
			return nil, "", err
		}

		if clu.Removed != "" {
			return clu, "removed", nil
		}

		return clu, "active", nil
	}
}
