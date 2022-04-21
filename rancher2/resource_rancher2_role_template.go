package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2RoleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2RoleTemplateCreate,
		Read:   resourceRancher2RoleTemplateRead,
		Update: resourceRancher2RoleTemplateUpdate,
		Delete: resourceRancher2RoleTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2RoleTemplateImport,
		},

		Schema: roleTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2RoleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	roleTemplate := expandRoleTemplate(d)

	log.Printf("[INFO] Creating role template")

	newRoleTemplate, err := client.RoleTemplate.Create(roleTemplate)
	if err != nil {
		return err
	}

	d.SetId(newRoleTemplate.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    roleTemplateStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for role template (%s) to be active: %s", d.Id(), waitErr)
	}

	return resourceRancher2RoleTemplateRead(d, meta)
}

func resourceRancher2RoleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing role template ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	roleTemplate, err := client.RoleTemplate.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] role template ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenRoleTemplate(d, roleTemplate)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2RoleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating role template ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	roleTemplate, err := client.RoleTemplate.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"administrative":  d.Get("administrative").(bool),
		"context":         d.Get("context").(string),
		"description":     d.Get("description").(string),
		"external":        d.Get("external").(bool),
		"hidden":          d.Get("hidden").(bool),
		"locked":          d.Get("locked").(bool),
		"name":            d.Get("name").(string),
		"roleTemplateIds": toArrayString(d.Get("role_template_ids").([]interface{})),
		"rules":           expandPolicyRules(d.Get("rules").([]interface{})),
		"annotations":     toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":          toMapString(d.Get("labels").(map[string]interface{})),
	}

	switch update["context"] {
	case roleTemplateContextCluster:
		update["clusterCreatorDefault"] = d.Get("default_role").(bool)
		update["projectCreatorDefault"] = false
	case roleTemplateContextProject:
		update["clusterCreatorDefault"] = false
		update["projectCreatorDefault"] = d.Get("default_role").(bool)
	}

	_, err = client.RoleTemplate.Update(roleTemplate, update)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    roleTemplateStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for role template (%s) to be updated: %s", d.Id(), waitErr)
	}

	return resourceRancher2RoleTemplateRead(d, meta)
}

func resourceRancher2RoleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting role template ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	roleTemplate, err := client.RoleTemplate.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Role template ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	if !roleTemplate.Builtin {
		err = client.RoleTemplate.Delete(roleTemplate)
		if err != nil {
			return fmt.Errorf("Error removing role template: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    roleTemplateStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for role template (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func roleTemplateStateRefreshFunc(client *managementClient.Client, roleTemplateID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.RoleTemplate.ByID(roleTemplateID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "exists", nil
	}
}
