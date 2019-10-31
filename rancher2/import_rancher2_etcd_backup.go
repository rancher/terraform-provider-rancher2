package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2EtcdBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	clusterDriver, err := client.EtcdBackup.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenEtcdBackup(d, clusterDriver)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
