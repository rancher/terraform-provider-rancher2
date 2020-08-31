package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterDriverEKSImport = "EKS"
)

//Schemas

func clusterEKSImportFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cloud_credential": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The EKS cloud_credential id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The EKS cluster name",
		},
		"region": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The EKS region",
		},
	}

	return s
}
