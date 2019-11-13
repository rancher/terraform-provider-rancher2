package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2MultiClusterAppImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceID := "cattle-global-data:" + d.Id()

	d.SetId(resourceID)

	err := resourceRancher2MultiClusterAppRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
