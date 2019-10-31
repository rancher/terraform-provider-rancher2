package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2AppImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID, appID, err := splitAppID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(appID)
	d.Set("project_id", projectID)

	err = resourceRancher2AppRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
