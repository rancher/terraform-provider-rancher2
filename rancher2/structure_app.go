package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	projectClient "github.com/rancher/types/client/project/v3"
)

// Flatteners

func flattenApp(d *schema.ResourceData, in *projectClient.App) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	d.Set("project_id", in.ProjectID)
	d.Set("name", in.Name)
	d.Set("target_namespace", in.TargetNamespace)
	d.Set("external_id", in.ExternalID)
	d.Set("description", in.Description)
	d.Set("values_yaml", in.ValuesYaml)

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("answers", toMapInterface(in.Answers))
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

func expandApp(in *schema.ResourceData) *projectClient.App {
	obj := &projectClient.App{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID
	obj.Name = in.Get("name").(string)
	obj.TargetNamespace = in.Get("target_namespace").(string)
	obj.ExternalID = in.Get("external_id").(string)
	obj.Description = in.Get("description").(string)
	obj.ValuesYaml = in.Get("values_yaml").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("answers").(map[string]interface{}); ok && len(v) > 0 {
		obj.Answers = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
