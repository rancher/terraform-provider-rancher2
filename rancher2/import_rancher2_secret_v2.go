package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2SecretV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID, _ := splitID(d.Id())
	d.Set("cluster_id", clusterID)

	err := resourceRancher2SecretV2Read(d, meta)
	if err != nil || d.Id() == "" {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
