package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterRoleTemplateBinding() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterRoleTemplateBindingRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_principal_id": {
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

func dataSourceRancher2ClusterRoleTemplateBindingRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	roleTemplateID := d.Get("role_template_id").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
		"name":      name,
	}
	if len(roleTemplateID) > 0 {
		filters["roleTemplateId"] = roleTemplateID
	}
	listOpts := NewListOpts(filters)

	clusterRoleTemplateBindings, err := client.ClusterRoleTemplateBinding.List(listOpts)
	if err != nil {
		return err
	}

	count := len(clusterRoleTemplateBindings.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster role template binding with name \"%s\" on cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster role template binding with name \"%s\" on cluster ID \"%s\"", count, name, clusterID)
	}

	return flattenClusterRoleTemplateBinding(d, &clusterRoleTemplateBindings.Data[0])
}
