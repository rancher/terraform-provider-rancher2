package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenRoleTemplateDefault(in *managementClient.RoleTemplate) bool {
	switch in.Context {
	case roleTemplateContextCluster:
		return in.ClusterCreatorDefault
	case roleTemplateContextProject:
		return in.ProjectCreatorDefault
	default:
		return false
	}
}

func flattenRoleTemplate(d *schema.ResourceData, in *managementClient.RoleTemplate) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening role template: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("administrative", in.Administrative)
	d.Set("builtin", in.Builtin)
	d.Set("default_role", flattenRoleTemplateDefault(in))
	d.Set("context", in.Context)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	d.Set("external", in.External)
	d.Set("hidden", in.Hidden)
	d.Set("locked", in.Locked)

	err := d.Set("role_template_ids", toArrayInterface(in.RoleTemplateIDs))
	if err != nil {
		return err
	}

	err = d.Set("rules", flattenPolicyRules(in.Rules))
	if err != nil {
		return err
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

func expandRoleTemplateDefault(def bool, in *managementClient.RoleTemplate) {
	switch in.Context {
	case roleTemplateContextCluster:
		in.ClusterCreatorDefault = def
	case roleTemplateContextProject:
		in.ProjectCreatorDefault = def
	}
}

func expandRoleTemplate(in *schema.ResourceData) *managementClient.RoleTemplate {
	obj := &managementClient.RoleTemplate{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.Administrative = in.Get("administrative").(bool)
	obj.Context = in.Get("context").(string)
	expandRoleTemplateDefault(in.Get("default_role").(bool), obj)

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	obj.External = in.Get("external").(bool)
	obj.Hidden = in.Get("hidden").(bool)
	obj.Locked = in.Get("locked").(bool)

	if v, ok := in.Get("role_template_ids").([]interface{}); ok && len(v) > 0 {
		obj.RoleTemplateIDs = toArrayString(v)
	}

	if v, ok := in.Get("rules").([]interface{}); ok && len(v) > 0 {
		obj.Rules = expandPolicyRules(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
