package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

// Flatteners

func flattenRegistryCredential(in map[string]projectClient.RegistryCredential, p []interface{}) []interface{} {
	if len(in) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(in))
	lenP := len(p)
	i := 0
	for key := range in {
		var obj map[string]interface{}
		if lenP <= i {
			obj = make(map[string]interface{})
		} else {
			obj = p[i].(map[string]interface{})
		}

		obj["address"] = key

		if len(in[key].Password) > 0 {
			obj["password"] = in[key].Password
		}

		if len(in[key].Username) > 0 {
			obj["username"] = in[key].Username
		}

		out[i] = obj
		i++
	}

	return out
}

func flattenDockerCredential(d *schema.ResourceData, in *projectClient.DockerCredential) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	d.Set("project_id", in.ProjectID)
	d.Set("name", in.Name)
	d.Set("description", in.Description)

	v, ok := d.Get("registries").([]interface{})
	if !ok {
		v = []interface{}{}
	}
	registryCredentials := flattenRegistryCredential(in.Registries, v)
	err := d.Set("registries", registryCredentials)
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

func flattenNamespacedDockerCredential(d *schema.ResourceData, in *projectClient.NamespacedDockerCredential) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	d.Set("project_id", in.ProjectID)
	d.Set("name", in.Name)
	d.Set("description", in.Description)
	d.Set("namespace_id", in.NamespaceId)

	v, ok := d.Get("registries").([]interface{})
	if !ok {
		v = []interface{}{}
	}
	registryCredentials := flattenRegistryCredential(in.Registries, v)
	err := d.Set("registries", registryCredentials)
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

func flattenRegistry(d *schema.ResourceData, in interface{}) error {
	namespaceID := d.Get("namespace_id").(string)
	if len(namespaceID) > 0 {
		return flattenNamespacedDockerCredential(d, in.(*projectClient.NamespacedDockerCredential))
	}

	return flattenDockerCredential(d, in.(*projectClient.DockerCredential))

}

// Expanders

func expandRegistryCredential(p []interface{}) map[string]projectClient.RegistryCredential {
	if len(p) == 0 || p[0] == nil {
		return map[string]projectClient.RegistryCredential{}
	}

	obj := make(map[string]projectClient.RegistryCredential)

	for i := range p {
		in := p[i].(map[string]interface{})
		aux := projectClient.RegistryCredential{}
		key := in["address"].(string)

		if v, ok := in["password"].(string); ok && len(v) > 0 {
			aux.Password = v
		}

		if v, ok := in["username"].(string); ok && len(v) > 0 {
			aux.Username = v
		}

		obj[key] = aux
	}

	return obj
}

func expandDockerCredential(in *schema.ResourceData) *projectClient.DockerCredential {
	obj := &projectClient.DockerCredential{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID
	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	if v, ok := in.Get("registries").([]interface{}); ok && len(v) > 0 {
		obj.Registries = expandRegistryCredential(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandNamespacedDockerCredential(in *schema.ResourceData) *projectClient.NamespacedDockerCredential {
	obj := &projectClient.NamespacedDockerCredential{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID
	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)
	obj.NamespaceId = in.Get("namespace_id").(string)

	if v, ok := in.Get("registries").([]interface{}); ok && len(v) > 0 {
		obj.Registries = expandRegistryCredential(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandRegistry(in *schema.ResourceData) interface{} {
	namespaceID := in.Get("namespace_id").(string)
	if len(namespaceID) > 0 {
		return expandNamespacedDockerCredential(in)
	}

	return expandDockerCredential(in)
}
