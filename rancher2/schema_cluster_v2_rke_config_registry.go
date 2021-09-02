package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigRegistryConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"hostname": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Registry hostname",
		},
		"auth_config_secret_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Registry auth config secret name",
		},
		"tls_secret_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Registry TLS secret name. TLS is a pair of Cert/Key",
		},
		"ca_bundle": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Registry CA bundle",
		},
		"insecure": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Registry insecure connectivity",
		},
	}

	return s
}

func clusterV2RKEConfigRegistryMirrorFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoints": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Registry mirror endpoints",
		},
		"hostname": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Registry hostname",
		},
		"rewrites": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Registry mirror rewrites",
		},
	}

	return s
}

func clusterV2RKEConfigRegistryFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"configs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Registry config",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigRegistryConfigFields(),
			},
		},
		"mirrors": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Registry mirrors",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigRegistryMirrorFields(),
			},
		},
	}

	return s
}
