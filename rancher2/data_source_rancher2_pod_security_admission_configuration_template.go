package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2PodSecurityAdmissionConfigurationTemplate() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceRancher2PodSecurityAdmissionConfigurationTemplateRead,
		Schema: podSecurityAdmissionConfigurationTemplateFields(),
	}
}

func dataSourceRancher2PodSecurityAdmissionConfigurationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	psact, err := client.PodSecurityAdmissionConfigurationTemplate.ByID(name)
	if err != nil {
		return err
	}

	return flattenPodSecurityAdmissionConfigurationTemplate(d, psact)
}
