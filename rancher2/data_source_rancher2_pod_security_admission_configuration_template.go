package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2PodSecurityAdmissionConfigurationTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2PodSecurityAdmissionConfigurationTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pod Security Admission Configuration Template name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pod Security Admission Configuration Template description",
			},
			// defaults is NOT required for data source
			"defaults": {
				Type:        schema.TypeList,
				Description: "defaults allows the user to define admission control mode for Pod Security",
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: podSecurityAdmissionConfigurationTemplateDefaultFields(),
				},
			},
			"exemptions": {
				Type:        schema.TypeList,
				Description: "exemptions allows the creation of pods for specific Usernames, RuntimeClassNames, and Namespaces that would otherwise be prohibited",
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: podSecurityAdmissionConfigurationTemplateExemptionFields(),
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Annotations of the Pod Security Admission Configuration Template",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Labels of the Pod Security Admission Configuration Template",
			},
		},
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
