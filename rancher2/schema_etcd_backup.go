package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func etcdBackupFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"backup_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesEtcdBackupConfigFields(),
			},
		},
		"filename": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"manual": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"namespace_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"annotations": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["annotations"],
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["labels"],
		},
	}

	return s
}
