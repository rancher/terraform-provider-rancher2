package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2ConfigMapV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID, name := splitID(d.Id())
	d.Set("cluster_id", clusterID)
	d.Set("name", name)

	err := resourceRancher2ConfigMapV2Read(d, meta)
	if err != nil || d.Id() == "" {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
