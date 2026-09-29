package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Types

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
		// use_internal_ca_certs makes the provider populate ca_certs from the
		// cluster's own internally generated CA certificate. Only use this when:
		// - The cluster sits behind a Layer 4 (TCP passthrough) load balancer.
		// - The endpoint FQDN is included in the API server certificate's SANs
		//   (configured via tls-san in machine_global_config).
		// Requires fqdn to be set, and is mutually exclusive with ca_certs.
		// use_internal_ca_certs is never derived from server state: its value in
		// Terraform state always matches what was last configured, never a
		// server-side comparison (see rancher2/resource_rancher2_cluster_v2.go).
		"use_internal_ca_certs": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	return s
}

// clusterV2RKEConfigLocalAuthEndpointFields returns the same fields as
// clusterV2LocalAuthEndpointFields, but without the use_internal_ca_certs field.
// This is used for the deprecated RKEConfig.local_auth_endpoint,
// which does not support that field.
func clusterV2RKEConfigLocalAuthEndpointFields() map[string]*schema.Schema {
	s := clusterV2LocalAuthEndpointFields()
	delete(s, "use_internal_ca_certs")
	return s
}
