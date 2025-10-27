package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

func dataSourceRancher2Secret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2SecretRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID to add secret",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the secret",
			},
			"data": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Secret data base64 encoded",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the secret",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace ID to add secret",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the secret",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the secret",
			},
		},
	}
}

func dataSourceRancher2SecretRead(d *schema.ResourceData, meta interface{}) error {
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

	secrets, err := meta.(*Config).GetSecretByFilters(filters)
	if err != nil {
		return err
	}

	switch t := secrets.(type) {
	case *projectClient.NamespacedSecretCollection:
		err = dataSourceRancher2SecretCheck(len(secrets.(*projectClient.NamespacedSecretCollection).Data), projectID, name)
		if err != nil {
			return err
		}
		return flattenSecret(d, &secrets.(*projectClient.NamespacedSecretCollection).Data[0])
	case *projectClient.SecretCollection:
		err = dataSourceRancher2SecretCheck(len(secrets.(*projectClient.SecretCollection).Data), projectID, name)
		if err != nil {
			return err
		}
		return flattenSecret(d, &secrets.(*projectClient.SecretCollection).Data[0])
	default:
		return fmt.Errorf("[ERROR] secret type %s isn't supported", t)
	}
}

func dataSourceRancher2SecretCheck(i int, projectID, name string) error {
	if i <= 0 {
		return fmt.Errorf("[ERROR] secret with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if i > 1 {
		return fmt.Errorf("[ERROR] found %d secret with name \"%s\" on project ID \"%s\"", i, name, projectID)
	}
	return nil
}
