package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenGlobalRole(d *schema.ResourceData, in *managementClient.GlobalRole) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening global role: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("builtin", in.Builtin)
	d.Set("new_user_default", in.NewUserDefault)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	err := d.Set("rules", flattenPolicyRules(in.Rules))
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

func expandGlobalRole(in *schema.ResourceData) *managementClient.GlobalRole {
	obj := &managementClient.GlobalRole{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.NewUserDefault = in.Get("new_user_default").(bool)

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
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
