package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2App() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2AppRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the app",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID to add app",
			},
			"target_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Namespace name to add app",
			},
			"answers": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Answers of the app",
			},
			"catalog_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Catalog name of the app",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External ID of the app",
			},
			"revision_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "App revision id",
			},
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template name of the app",
			},
			"template_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template version of the app",
			},
			"values_yaml": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "values.yaml file content of the app",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Annotations of the app",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the app",
			},
		},
	}
}

func dataSourceRancher2AppRead(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	targetNamespace := d.Get("target_namespace").(string)

	filters := map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	}

	if len(targetNamespace) > 0 {
		filters["targetNamespace"] = targetNamespace
	}

	listOpts := NewListOpts(filters)

	client, err := meta.(*Config).ProjectClient(projectID)
	if err != nil {
		return err
	}

	apps, err := client.App.List(listOpts)
	if err != nil {
		return err
	}

	count := len(apps.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] app with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d app with name \"%s\" on project ID \"%s\"", count, name, projectID)
	}

	return flattenApp(d, &apps.Data[0])
}
