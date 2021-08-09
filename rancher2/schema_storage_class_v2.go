package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	norman "github.com/rancher/norman/types"
	"k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
)

const (
	storageClassV2Kind                        = "StorageClass"
	storageClassV2APIVersion                  = "storage.k8s.io/v1"
	storageClassV2APIType                     = "storage.k8s.io.storageclass"
	storageClassV2ClusterIDsep                = "."
	storageClassV2ReclaimRecycle              = string(v1.PersistentVolumeReclaimRecycle)
	storageClassV2ReclaimDelete               = string(v1.PersistentVolumeReclaimDelete)
	storageClassV2ReclaimRetain               = string(v1.PersistentVolumeReclaimRetain)
	storageClassV2BindingImmediate            = string(storageV1.VolumeBindingImmediate)
	storageClassV2BindingWaitForFirstConsumer = string(storageV1.VolumeBindingWaitForFirstConsumer)
)

var (
	storageClassV2ReclaimList = []string{
		storageClassV2ReclaimRecycle,
		storageClassV2ReclaimDelete,
		storageClassV2ReclaimRetain,
	}
	storageClassV2BindingList = []string{
		storageClassV2BindingImmediate,
		storageClassV2BindingWaitForFirstConsumer,
	}
)

//Types

type StorageClassV2 struct {
	norman.Resource
	storageV1.StorageClass
}

func storageClassV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "K8s cluster ID",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "StorageClass name",
		},
		"k8s_provisioner": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "StorageClass provisioner",
		},
		"allow_volume_expansion": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "StorageClass allow_volume_expansion",
		},
		"mount_options": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "StorageClass mount options",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"parameters": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "StorageClass provisioner paramaters",
		},
		"reclaim_policy": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      storageClassV2ReclaimDelete,
			Description:  "StorageClass provisioner reclaim policy",
			ValidateFunc: validation.StringInSlice(storageClassV2ReclaimList, true),
		},
		"volume_binding_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      storageClassV2BindingImmediate,
			Description:  "StorageClass provisioner volume binding mode",
			ValidateFunc: validation.StringInSlice(storageClassV2BindingList, true),
		},
		"resource_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
