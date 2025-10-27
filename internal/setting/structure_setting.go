package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenSetting(d *schema.ResourceData, in *managementClient.Setting) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening setting: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("value", in.Value)

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

func expandSetting(in *schema.ResourceData) (*managementClient.Setting, error) {
	obj := &managementClient.Setting{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding setting: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.Value = in.Get("value").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
