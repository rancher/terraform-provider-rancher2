package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Shemas

func certificateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"certs": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Certificate certs base64 encoded",
		},
		"key": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Certificate key base64 encoded",
		},
		"project_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Project ID to add certificate",
		},
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Certificate description",
		},
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Certificate name",
		},
		"namespace_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Namespace ID to add certificate",
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Annotations of the certificate",
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Labels of the certificate",
		},
	}

	return s
}
