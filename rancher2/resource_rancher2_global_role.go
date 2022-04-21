package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2GlobalRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2GlobalRoleCreate,
		Read:   resourceRancher2GlobalRoleRead,
		Update: resourceRancher2GlobalRoleUpdate,
		Delete: resourceRancher2GlobalRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2GlobalRoleImport,
		},

		Schema: globalRoleFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2GlobalRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalRole := expandGlobalRole(d)

	log.Printf("[INFO] Creating global role")

	newGlobalRole, err := client.GlobalRole.Create(globalRole)
	if err != nil {
		return err
	}

	d.SetId(newGlobalRole.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    globalRoleStateRefreshFunc(client, newGlobalRole.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for Global Role (%s) to be created: %s", newGlobalRole.ID, waitErr)
	}

	return resourceRancher2GlobalRoleRead(d, meta)
}

func resourceRancher2GlobalRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing global role ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalRole, err := client.GlobalRole.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] global role ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenGlobalRole(d, globalRole)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2GlobalRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Updating global role ID %s", id)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalRole, err := client.GlobalRole.ByID(id)
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"description":    d.Get("description").(string),
		"name":           d.Get("name").(string),
		"newUserDefault": d.Get("new_user_default").(bool),
		"rules":          expandPolicyRules(d.Get("rules").([]interface{})),
		"annotations":    toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":         toMapString(d.Get("labels").(map[string]interface{})),
	}

	_, err = client.GlobalRole.Update(globalRole, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalRoleStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for Global Role (%s) to be updated: %s", id, waitErr)
	}

	return resourceRancher2GlobalRoleRead(d, meta)
}

func resourceRancher2GlobalRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting global role ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalRole, err := client.GlobalRole.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Global role ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	if !globalRole.Builtin {
		err = client.GlobalRole.Delete(globalRole)
		if err != nil {
			return fmt.Errorf("Error removing global role: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    globalRoleStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for Global Role (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func globalRoleStateRefreshFunc(client *managementClient.Client, globalRoleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.GlobalRole.ByID(globalRoleID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}
