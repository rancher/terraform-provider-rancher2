package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2LocalAuthEndpointFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ca_certs": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"use_internal_ca_certs"},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"fqdn": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"use_internal_ca_certs": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"ca_certs"},
		},
	}

	return s
}