package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	tokenDefaultSessionDesc = "Terraform token temp token"
	tokenDefaultTTL         = "60000"
)

//Schemas

func tokenFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token access key",
		},
		"cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Cluster ID for scoped token",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Token description",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Token enabled",
		},
		"expired": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Token expired",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token name",
		},
		"renew": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
			Description: "Renew expired or disabled token",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "Token secret key",
		},
		"token": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "Token value",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Token time to live in seconds",
		},
		"user_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token user ID",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
