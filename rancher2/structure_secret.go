package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

// Flatteners

func flattenProjectSecret(d *schema.ResourceData, in *projectClient.Secret) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("data", toMapInterface(in.Data))
	if err != nil {
		return err
	}

	d.Set("project_id", in.ProjectID)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
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

func flattenNamespacedSecret(d *schema.ResourceData, in *projectClient.NamespacedSecret) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("data", toMapInterface(in.Data))
	if err != nil {
		return err
	}

	d.Set("project_id", in.ProjectID)
	d.Set("namespace_id", in.NamespaceId)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
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

func flattenSecret(d *schema.ResourceData, in interface{}) error {
	namespaceID := d.Get("namespace_id").(string)
	if len(namespaceID) > 0 {
		return flattenNamespacedSecret(d, in.(*projectClient.NamespacedSecret))
	}

	return flattenProjectSecret(d, in.(*projectClient.Secret))

}

// Expanders

func expandProjectSecret(in *schema.ResourceData) *projectClient.Secret {
	obj := &projectClient.Secret{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("data").(map[string]interface{}); ok && len(v) > 0 {
		obj.Data = toMapString(v)
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandNamespacedSecret(in *schema.ResourceData) *projectClient.NamespacedSecret {
	obj := &projectClient.NamespacedSecret{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("data").(map[string]interface{}); ok && len(v) > 0 {
		obj.Data = toMapString(v)
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	obj.NamespaceId = in.Get("namespace_id").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandSecret(in *schema.ResourceData) interface{} {
	namespaceID := in.Get("namespace_id").(string)
	if len(namespaceID) > 0 {
		return expandNamespacedSecret(in)
	}

	return expandProjectSecret(in)
}
