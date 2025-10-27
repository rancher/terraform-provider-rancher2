package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"k8s.io/api/core/v1"
)

// Flatteners

func flattenSecretV2(d *schema.ResourceData, in *SecretV2) error {
	if in == nil {
		return nil
	}

	if len(in.ID) > 0 {
		d.SetId(d.Get("cluster_id").(string) + secretV2ClusterIDsep + in.ID)
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
	if len(in.K8SType) > 0 {
		d.Set("type", string(in.K8SType))
	}

	if len(in.Data) > 0 {
		result := make(map[string]string, len(in.Data))
		for k, v := range in.Data {
			result[k] = string(v)
		}
		err = d.Set("data", toMapInterface(result))
		if err != nil {
			return err
		}
	}

	return nil
}

// Expanders

func expandSecretV2(in *schema.ResourceData) *SecretV2 {
	if in == nil {
		return nil
	}
	obj := &SecretV2{}

	if len(in.Id()) > 0 {
		_, obj.ID = splitID(in.Id())
	}
	obj.TypeMeta.Kind = secretV2Kind
	obj.TypeMeta.APIVersion = secretV2APIVersion

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
	if v, ok := in.Get("type").(string); ok && len(v) > 0 {
		obj.Resource.Type = v
		obj.K8SType = v
	}
	if v, ok := in.Get("data").(map[string]interface{}); ok && len(v) > 0 {
		obj.StringData = toMapString(v)
	}

	return obj
}
