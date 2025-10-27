package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenFeature(d *schema.ResourceData, in *managementClient.Feature) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening feature: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	if in.Value != nil {
		d.Set("value", *in.Value)
	}

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

func expandFeature(in *schema.ResourceData) (*managementClient.Feature, error) {
	obj := &managementClient.Feature{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding feature: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}
	if v, ok := in.Get("value").(bool); ok {
		obj.Value = &v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
