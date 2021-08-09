package rancher2

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2StorageClassV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2StorageClassV2Read,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "K8s cluster ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "StorageClass name",
			},
			"k8s_provisioner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "StorageClass provisioner",
			},
			"allow_volume_expansion": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "StorageClass allow_volume_expansion",
			},
			"mount_options": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "StorageClass mount options",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parameters": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "StorageClass provisioner paramaters",
			},
			"reclaim_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "StorageClass provisioner reclaim policy",
			},
			"volume_binding_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "StorageClass provisioner volume binding mode",
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2StorageClassV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	storageClass, err := getStorageClassV2ByID(meta.(*Config), clusterID, name)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] StorageClass V2 %s not found at cluster %s", name, clusterID)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenStorageClassV2(d, storageClass)
}
