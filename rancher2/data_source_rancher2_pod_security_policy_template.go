package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2PodSecurityPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2PodSecurityPolicyTemplateRead,
		Schema:      podSecurityPolicyTemplateFields(),
	}
}

func dataSourceRancher2PodSecurityPolicyTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	pspt, err := client.PodSecurityPolicyTemplate.ByID(name)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenPodSecurityPolicyTemplate(d, pspt))
}
