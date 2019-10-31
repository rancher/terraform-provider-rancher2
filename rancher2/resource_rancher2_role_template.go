package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		if IsNotFound(err) {
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
		if IsNotFound(err) {
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

	d.SetId("")
	return nil
}
