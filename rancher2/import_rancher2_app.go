package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRancher2AppImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID, resourceID := splitAppID(d.Id())

	d.SetId(resourceID)
	d.Set("project_id", projectID)

	err := resourceRancher2AppRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
