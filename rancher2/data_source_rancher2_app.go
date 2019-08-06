package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	projectClient "github.com/rancher/types/client/project/v3"
)

func dataSourceRancher2App() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2AppRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID to add app",
			},
			"target_namespace": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Namespace name to add app",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the app",
			},
			"external_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External ID of the app, like catalog://?catalog=helm&template=cluster-autoscaler&version=3.1.0",
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"answers": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Answers of the app",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the app",
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2AppRead(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)
	targetNamespace := d.Get("target_namespace").(string)

	filters := map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	}

	if len(targetNamespace) > 0 {
		filters["targetNamespace"] = targetNamespace
	}

	apps, err := meta.(*Config).GetAppByFilters(filters)
	if err != nil {
		return err
	}

	err = dataSourceRancher2AppCheck(len(apps.(*projectClient.AppCollection).Data), projectID, name)
	if err != nil {
		return err
	}

	return flattenApp(d, &apps.(*projectClient.AppCollection).Data[0])
}

func dataSourceRancher2AppCheck(i int, projectID, name string) error {
	if i <= 0 {
		return fmt.Errorf("[ERROR] app with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if i > 1 {
		return fmt.Errorf("[ERROR] found %d app with name \"%s\" on project ID \"%s\"", i, name, projectID)
	}
	return nil
}
