package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2User() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2UserRead,

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

func dataSourceRancher2UserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	username := d.Get("username").(string)
	name := d.Get("name").(string)
	externalUser := d.Get("is_external").(bool)

	if len(username) == 0 && len(name) == 0 {
		log.Printf("[WARN] username and name filters are nil")
		return nil
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	count := len(users.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] user with username \"%s\" and/or name \"%s\" not found", username, name)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d users username \"%s\" and/or name \"%s\"", count, username, name)
	}

	return diag.FromErr(flattenUser(d, &users.Data[0]))
}
