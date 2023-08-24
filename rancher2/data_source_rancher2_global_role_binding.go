package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2GlobalRoleBinding() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2GlobalRoleBindingRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"global_role_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_principal_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceRancher2GlobalRoleBindingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	globalRole := d.Get("global_role_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	if len(globalRole) > 0 {
		filters["globalRoleId"] = globalRole
	}
	listOpts := NewListOpts(filters)

	globalRoleBindings, err := client.GlobalRoleBinding.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(globalRoleBindings.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] global role binding with name \"%s\" not found", name)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d global role binding with name \"%s\"", count, name)
	}

	return diag.FromErr(flattenGlobalRoleBinding(d, &globalRoleBindings.Data[0]))
}
