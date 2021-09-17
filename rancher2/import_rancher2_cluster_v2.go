package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2ClusterV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRancher2ClusterV2Read(d, meta)
	if err != nil || d.Id() == "" {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
