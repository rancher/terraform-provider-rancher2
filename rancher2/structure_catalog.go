package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenCatalog(d *schema.ResourceData, in *managementClient.Catalog) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("url", in.URL)
	d.Set("description", in.Description)
	d.Set("kind", in.Kind)
	d.Set("branch", in.Branch)

	err := d.Set("annotations", toMapInterface(in.Annotations))
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

func expandCatalog(in *schema.ResourceData) *managementClient.Catalog {
	obj := &managementClient.Catalog{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.URL = in.Get("url").(string)
	obj.Description = in.Get("description").(string)
	obj.Kind = in.Get("kind").(string)
	obj.Branch = in.Get("branch").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
