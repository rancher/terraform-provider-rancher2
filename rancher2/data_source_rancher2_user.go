package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRancher2User() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2UserRead,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"principal_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2UserRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	username := d.Get("username").(string)

	filters := map[string]interface{}{
		"username": username,
	}
	listOpts := NewListOpts(filters)

	users, err := client.User.List(listOpts)
	if err != nil {
		return err
	}

	count := len(users.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] username \"%s\" not found", username)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d username \"%s\"", count, username)
	}

	return flattenUser(d, &users.Data[0])
}
