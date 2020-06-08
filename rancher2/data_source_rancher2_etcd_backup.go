package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2EtcdBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2EtcdBackupRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRKEConfigServicesEtcdBackupConfigFields(),
				},
			},
			"filename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manual": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: descriptions["annotations"],
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: descriptions["labels"],
			},
		},
	}
}

func dataSourceRancher2EtcdBackupRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
		"name":      name,
	}
	listOpts := NewListOpts(filters)

	etcdBackups, err := client.EtcdBackup.List(listOpts)
	if err != nil {
		return err
	}

	count := len(etcdBackups.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] etcd backup with name \"%s\" on cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d etcd backup with name \"%s\" on cluster ID \"%s\"", count, name, clusterID)
	}

	return flattenEtcdBackup(d, &etcdBackups.Data[0])
}
