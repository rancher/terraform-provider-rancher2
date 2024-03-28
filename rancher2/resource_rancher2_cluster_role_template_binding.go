package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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

		Schema: clusterRoleTemplateBindingFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ClusterRoleTemplateBindingCreate(d *schema.ResourceData, meta interface{}) error {
	clusterRole := expandClusterRoleTemplateBinding(d)

	err := meta.(*Config).ClusterExist(clusterRole.ClusterID)
	if err != nil {
		return err
	}

	err = meta.(*Config).RoleTemplateExist(clusterRole.RoleTemplateID)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Cluster Role Template Binding %s", clusterRole.Name)

	newClusterRole, err := client.ClusterRoleTemplateBinding.Create(clusterRole)
	if err != nil {
		return err
	}

	d.SetId(newClusterRole.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    clusterRoleTemplateBindingStateRefreshFunc(client, newClusterRole.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster role template binding (%s) to be created: %s", newClusterRole.ID, waitErr)
	}

	return resourceRancher2ClusterRoleTemplateBindingRead(d, meta)
}

func resourceRancher2ClusterRoleTemplateBindingRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster Role Template Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		clusterRole, err := client.ClusterRoleTemplateBinding.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster Role Template Binding ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenClusterRoleTemplateBinding(d, clusterRole); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
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
		"groupId":          d.Get("group_id").(string),
		"groupPrincipalId": d.Get("group_principal_id").(string),
		"roleTemplateId":   d.Get("role_template_id").(string),
		"userId":           d.Get("user_id").(string),
		"userPrincipalId":  d.Get("user_principal_id").(string),
		"annotations":      toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":           toMapString(d.Get("labels").(map[string]interface{})),
	}

	newClusterRole, err := client.ClusterRoleTemplateBinding.Update(clusterRole, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    clusterRoleTemplateBindingStateRefreshFunc(client, newClusterRole.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
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
		if IsNotFound(err) || IsForbidden(err) {
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
		Timeout:    d.Timeout(schema.TimeoutDelete),
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

// clusterRoleTemplateBindingStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster Role Template Binding.
func clusterRoleTemplateBindingStateRefreshFunc(client *managementClient.Client, clusterRoleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ClusterRoleTemplateBinding.ByID(clusterRoleID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, "active", nil
	}
}
