package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2NodeDriverImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	nodeDriver, err := client.NodeDriver.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenNodeDriver(d, nodeDriver)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
