package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	"k8s.io/api/core/v1"
)

const (
	configMapV2Kind             = "ConfigMap"
	configMapV2APIVersion       = "v1"
	configMapV2APIType          = "configmap"
	configMapV2ClusterIDsep     = "."
	configMapV2ActiveCondition  = "Updated"
	configMapV2CreatedCondition = "Created"
)

//Types

type ConfigMapV2 struct {
	norman.Resource
	v1.ConfigMap
}

// Flatteners

func flattenConfigMapV2(d *schema.ResourceData, in *ConfigMapV2) error {
	if in == nil {
		return nil
	}

	if len(in.ID) > 0 {
		d.SetId(d.Get("cluster_id").(string) + configMapV2ClusterIDsep + in.ID)
	}
	d.Set("name", in.ObjectMeta.Name)
	d.Set("namespace", in.ObjectMeta.Namespace)
	err := d.Set("annotations", toMapInterface(in.ObjectMeta.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.ObjectMeta.Labels))
	if err != nil {
		return err
	}
	d.Set("resource_version", in.ObjectMeta.ResourceVersion)

	if in.Immutable != nil {
		d.Set("immutable", *in.Immutable)
	}

	if len(in.Data) > 0 {
		err = d.Set("data", toMapInterface(in.Data))
		if err != nil {
			return err
		}
	}

	return nil
}

// Expanders

func expandConfigMapV2(in *schema.ResourceData) *ConfigMapV2 {
	if in == nil {
		return nil
	}
	obj := &ConfigMapV2{}

	if len(in.Id()) > 0 {
		_, obj.ID = splitID(in.Id())
	}
	obj.TypeMeta.Kind = configMapV2Kind
	obj.TypeMeta.APIVersion = configMapV2APIVersion

	obj.ObjectMeta.Name = in.Get("name").(string)
	obj.ObjectMeta.Namespace = in.Get("namespace").(string)
	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Annotations = toMapString(v)
	}
	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Labels = toMapString(v)
	}
	if v, ok := in.Get("resource_version").(string); ok {
		obj.ObjectMeta.ResourceVersion = v
	}
	if v, ok := in.Get("immutable").(bool); ok {
		obj.Immutable = &v
	}
	if v, ok := in.Get("data").(map[string]interface{}); ok && len(v) > 0 {
		obj.Data = toMapString(v)
	}

	return obj
}
