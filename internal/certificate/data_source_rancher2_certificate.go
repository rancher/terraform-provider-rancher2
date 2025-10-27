package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

func dataSourceRancher2Certificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2CertificateRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID to add certificate",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the certificate",
			},
			"certs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate certs base64 encoded",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the certificate",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace ID to add certificate",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the certificate",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the certificate",
			},
		},
	}
}

func dataSourceRancher2CertificateRead(d *schema.ResourceData, meta interface{}) error {
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

	certs, err := meta.(*Config).GetCertificateByFilters(filters)
	if err != nil {
		return err
	}

	switch t := certs.(type) {
	case *projectClient.NamespacedCertificateCollection:
		err = dataSourceRancher2CertificateCheck(len(certs.(*projectClient.NamespacedCertificateCollection).Data), projectID, name)
		if err != nil {
			return err
		}
		return flattenCertificate(d, &certs.(*projectClient.NamespacedCertificateCollection).Data[0])
	case *projectClient.CertificateCollection:
		err = dataSourceRancher2CertificateCheck(len(certs.(*projectClient.CertificateCollection).Data), projectID, name)
		if err != nil {
			return err
		}
		return flattenCertificate(d, &certs.(*projectClient.CertificateCollection).Data[0])
	default:
		return fmt.Errorf("[ERROR] certificate type %s isn't supported", t)
	}
}

func dataSourceRancher2CertificateCheck(i int, projectID, name string) error {
	if i <= 0 {
		return fmt.Errorf("[ERROR] certificate with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if i > 1 {
		return fmt.Errorf("[ERROR] found %d certificate with name \"%s\" on project ID \"%s\"", i, name, projectID)
	}
	return nil
}
