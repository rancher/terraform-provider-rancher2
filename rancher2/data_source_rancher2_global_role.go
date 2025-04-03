package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2GlobalRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2GlobalRoleRead,
		Schema: map[string]*schema.Schema{
			"builtin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Builtin global role",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Global role policy description",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Global role policy name",
			},
			"new_user_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this role should be added to new users",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Global role policy rules",
				Elem: &schema.Resource{
					Schema: policyRuleFields(),
				},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Global role uuid",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the global role",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the global role",
			},
			"inherited_cluster_roles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Names of role templates whose permissions are granted by this global role in every cluster besides the local cluster",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRancher2GlobalRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	listOpts := NewListOpts(filters)

	globalRoles, err := client.GlobalRole.List(listOpts)
	if err != nil {
		return err
	}

	count := len(globalRoles.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] global role with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d global role with name \"%s\"", count, name)
	}

	return flattenGlobalRole(d, &globalRoles.Data[0])
}
