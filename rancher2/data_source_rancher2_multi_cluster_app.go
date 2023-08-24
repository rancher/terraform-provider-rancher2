package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2MultiClusterApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2MultiClusterAppRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Multi cluster app name",
			},
			"targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multi cluster app targets",
				Elem: &schema.Resource{
					Schema: targetFields(),
				},
			},
			"catalog_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Multi cluster app catalog name",
			},
			"answers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multi cluster app answers",
				Elem: &schema.Resource{
					Schema: answerFields(),
				},
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multi cluster app members",
				Elem: &schema.Resource{
					Schema: memberFields(),
				},
			},
			"revision_history_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Multi cluster app revision history limit",
			},
			"revision_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Multi cluster app revision id",
			},
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multi cluster app members",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Multi cluster app template version",
			},
			"template_version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Multi cluster app template version ID",
			},
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Multi cluster app template name",
			},
			"upgrade_strategy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multi cluster app upgrade strategy",
				Elem: &schema.Resource{
					Schema: upgradeStrategyFields(),
				},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2MultiClusterAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}

	listOpts := NewListOpts(filters)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	multiClusterApps, err := client.MultiClusterApp.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(multiClusterApps.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] multi cluster app with name \"%s\" not found", name)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d multi cluster app with name \"%s\"", count, name)
	}

	templateVersion, err := client.TemplateVersion.ByID(multiClusterApps.Data[0].TemplateVersionID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenMultiClusterApp(d, &multiClusterApps.Data[0], templateVersion.ExternalID))
}
