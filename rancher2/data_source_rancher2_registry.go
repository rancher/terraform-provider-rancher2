package rancher2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

func dataSourceRancher2Registry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2RegistryRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID to add docker registry",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the docker registry",
			},
			"registries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: registryCredentialFields(),
				},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the docker registry",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace ID to add docker registry",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the docker registry",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the docker registry",
			},
		},
	}
}

func dataSourceRancher2RegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)
	namespaceID := d.Get("namespace_id").(string)

	filters := map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	}

	if len(namespaceID) > 0 {
		filters["namespaceId"] = namespaceID
	}

	registries, err := meta.(*Config).GetRegistryByFilters(filters)
	if err != nil {
		return diag.FromErr(err)
	}

	switch t := registries.(type) {
	case *projectClient.NamespacedDockerCredentialCollection:
		err = dataSourceRancher2RegistryCheck(len(registries.(*projectClient.NamespacedDockerCredentialCollection).Data), projectID, name)
		if err != nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(flattenRegistry(d, &registries.(*projectClient.NamespacedDockerCredentialCollection).Data[0]))
	case *projectClient.DockerCredentialCollection:
		err = dataSourceRancher2RegistryCheck(len(registries.(*projectClient.DockerCredentialCollection).Data), projectID, name)
		if err != nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(flattenRegistry(d, &registries.(*projectClient.DockerCredentialCollection).Data[0]))
	default:
		return diag.Errorf("[ERROR] Registry type %s isn't supported", t)
	}
}

func dataSourceRancher2RegistryCheck(i int, projectID, name string) error {
	if i <= 0 {
		return fmt.Errorf("[ERROR] registry with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if i > 1 {
		return fmt.Errorf("[ERROR] found %d registry with name \"%s\" on project ID \"%s\"", i, name, projectID)
	}
	return nil
}
