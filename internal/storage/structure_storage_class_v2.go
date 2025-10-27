package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
)

// Flatteners

func flattenStorageClassV2(d *schema.ResourceData, in *StorageClassV2) error {
	if in == nil {
		return nil
	}

	if len(in.ID) > 0 {
		d.SetId(d.Get("cluster_id").(string) + storageClassV2ClusterIDsep + in.ID)
	}
	d.Set("name", in.ObjectMeta.Name)
	err := d.Set("annotations", toMapInterface(in.ObjectMeta.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.ObjectMeta.Labels))
	if err != nil {
		return err
	}
	d.Set("resource_version", in.ObjectMeta.ResourceVersion)
	d.Set("k8s_provisioner", in.Provisioner)

	if in.AllowVolumeExpansion != nil {
		d.Set("allow_volume_expansion", *in.AllowVolumeExpansion)
	}
	d.Set("mount_options", toArrayInterfaceSorted(in.MountOptions))
	if err != nil {
		return err
	}
	if in.Parameters != nil && len(in.Parameters) > 0 {
		d.Set("parameters", toMapInterface(in.Parameters))
		if err != nil {
			return err
		}
	}
	if in.ReclaimPolicy != nil && len(*in.ReclaimPolicy) > 0 {
		d.Set("reclaim_policy", string(*in.ReclaimPolicy))
	}
	if in.VolumeBindingMode != nil && len(*in.VolumeBindingMode) > 0 {
		d.Set("volume_binding_mode", string(*in.VolumeBindingMode))
	}

	return nil
}

// Expanders

func expandStorageClassV2(in *schema.ResourceData) *StorageClassV2 {
	if in == nil {
		return nil
	}
	obj := &StorageClassV2{}

	if len(in.Id()) > 0 {
		_, obj.ID = splitID(in.Id())
	}
	obj.TypeMeta.Kind = storageClassV2Kind
	obj.TypeMeta.APIVersion = storageClassV2APIVersion

	obj.ObjectMeta.Name = in.Get("name").(string)
	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Annotations = toMapString(v)
	}
	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Labels = toMapString(v)
	}
	if v, ok := in.Get("resource_version").(string); ok {
		obj.ObjectMeta.ResourceVersion = v
	}
	if v, ok := in.Get("k8s_provisioner").(string); ok {
		obj.Provisioner = v
	}
	if v, ok := in.Get("allow_volume_expansion").(bool); ok {
		obj.AllowVolumeExpansion = &v
	}
	if v, ok := in.Get("mount_options").([]interface{}); ok {
		obj.MountOptions = toArrayStringSorted(v)
	}
	if v, ok := in.Get("parameters").(map[string]interface{}); ok && len(v) > 0 {
		obj.Parameters = toMapString(v)
	}
	if v, ok := in.Get("reclaim_policy").(string); ok {
		policy := v1.PersistentVolumeReclaimPolicy(v)
		obj.ReclaimPolicy = &policy
	}
	if v, ok := in.Get("volume_binding_mode").(string); ok {
		mode := storageV1.VolumeBindingMode(v)
		obj.VolumeBindingMode = &mode
	}

	return obj
}
