package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenUser(d *schema.ResourceData, in *managementClient.User) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("username", in.Username)
	d.Set("enabled", in.Enabled)

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}

	err := d.Set("principal_ids", toArrayInterface(in.PrincipalIDs))
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

func expandUser(in *schema.ResourceData) *managementClient.User {
	obj := &managementClient.User{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Password = in.Get("password").(string)
	obj.Username = in.Get("username").(string)
	enabled := in.Get("enabled").(bool)
	obj.Enabled = &enabled

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.PrincipalIDs = toArrayString(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
