package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2User() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2UserRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_external": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"principal_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

func dataSourceRancher2UserRead(d *schema.ResourceData, meta interface{}) error {
	username := d.Get("username").(string)
	name := d.Get("name").(string)
	externalUser := d.Get("is_external").(bool)

	if len(username) == 0 && len(name) == 0 {
		log.Printf("[WARN] username and name filters are nil")
		return nil
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}

	if len(username) > 0 {
		filters["username"] = username
	}

	if len(name) > 0 {
		filters["name"] = name
		if externalUser {
			filters["username"] = ""
		}
	}

	listOpts := NewListOpts(filters)

	users, err := client.User.List(listOpts)
	if err != nil {
		return err
	}

	count := len(users.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] user with username \"%s\" and/or name \"%s\" not found", username, name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d users username \"%s\" and/or name \"%s\"", count, username, name)
	}

	return flattenUser(d, &users.Data[0])
}
