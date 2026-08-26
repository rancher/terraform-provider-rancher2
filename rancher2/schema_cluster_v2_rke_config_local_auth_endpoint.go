package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2LocalAuthEndpointFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ca_certs": {
			Type:     schema.TypeString,
			Optional: true,
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
	}

	return s
}
