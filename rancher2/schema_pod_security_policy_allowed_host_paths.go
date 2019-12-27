package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicyAllowedHostPathsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"path_prefix": {
			Type:        schema.TypeString,
			Description: "pathPrefix is the path prefix that the host volume must match. It does not support `*`. Trailing slashes are trimmed when validating the path prefix with a host path.",
			Required:    true,
		},
		"read_only": {
			Type:        schema.TypeBool,
			Description: "when set to true, will allow host volumes matching the pathPrefix only if all volume mounts are readOnly.",
			Optional:    true,
		},
	}

	return s
}
