package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenEtcdBackup(d *schema.ResourceData, in *managementClient.EtcdBackup) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	backupConfig := d.Get("backup_config").([]interface{})
	err := d.Set("backup_config", flattenClusterRKEConfigServicesEtcdBackupConfig(in.BackupConfig, backupConfig))
	if err != nil {
		return err
	}

	d.Set("cluster_id", in.ClusterID)
	d.Set("filename", in.Filename)
	d.Set("manual", in.Manual)
	d.Set("name", in.Name)
	d.Set("namespace_id", in.NamespaceId)

	err = d.Set("annotations", toMapInterface(in.Annotations))
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

func expandEtcdBackup(in *schema.ResourceData) (*managementClient.EtcdBackup, error) {
	obj := &managementClient.EtcdBackup{}
	if in == nil {
		return nil, nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("backup_config").([]interface{}); ok && len(v) > 0 {
		backupConfig, err := expandClusterRKEConfigServicesEtcdBackupConfig(v)
		if err != nil {
			return nil, err
		}
		obj.BackupConfig = backupConfig
	}

	obj.ClusterID = in.Get("cluster_id").(string)

	if v, ok := in.Get("filename").(string); ok && len(v) > 0 {
		obj.Filename = v
	}

	if v, ok := in.Get("manual").(bool); ok {
		obj.Manual = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("namespace_id").(string); ok && len(v) > 0 {
		obj.NamespaceId = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
