package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigPrivateRegistriesECRCredentialsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"aws_access_key_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"aws_secret_access_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"aws_session_token": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
	return s
}

func clusterRKEConfigPrivateRegistriesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ecr_credential_plugin": {
			Type:        schema.TypeList,
			Description: "ECR credential plugin config",
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigPrivateRegistriesECRCredentialsFields(),
			},
		},
		"is_default": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"user": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
	return s
}
