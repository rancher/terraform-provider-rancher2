package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ProjectAlertGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ProjectAlertGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert group name",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert group project ID",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert group description",
			},
			"group_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert group interval seconds",
			},
			"group_wait_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert group wait seconds",
			},
			"recipients": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Alert group recipients",
				Elem: &schema.Resource{
					Schema: recipientFields(),
				},
			},
			"repeat_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert group repeat interval seconds",
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

func dataSourceRancher2ProjectAlertGroupRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name":      name,
		"projectId": projectID,
	}
	listOpts := NewListOpts(filters)

	alertGroups, err := client.ProjectAlertGroup.List(listOpts)
	if err != nil {
		return err
	}

	count := len(alertGroups.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] project alert group with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d project alert group with name \"%s\" on project ID \"%s\"", count, name, projectID)
	}

	return flattenProjectAlertGroup(d, &alertGroups.Data[0])
}
