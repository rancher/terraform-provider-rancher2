package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2PodSecurityAdmissionConfigurationTemplateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRancher2PodSecurityAdmissionConfigurationTemplateRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
