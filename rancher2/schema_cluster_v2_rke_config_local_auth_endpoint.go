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
			Computed: true,
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
		// (Optional) Use the cluster's self-signed CA certificate
		// for the authorized cluster endpoint. Only use this when:
		// - Using a Layer 4 load balancer (TCP passthrough)
		// - The cluster's API server certificate includes the FQDN
		//   in its SANs (configured via tls-san in machine_global_config)
		// Mutually exclusive with ca_certs.
		"use_internal_ca_certs": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	return s
}
