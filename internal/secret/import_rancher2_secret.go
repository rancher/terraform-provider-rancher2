package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2SecretImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	namespaceID, projectID, resourceID := splitRegistryID(d.Id())

	d.SetId(resourceID)
	d.Set("project_id", projectID)
	d.Set("namespace_id", namespaceID)

	err := resourceRancher2SecretRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
