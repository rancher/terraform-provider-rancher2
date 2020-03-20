package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicyRuntimeClassFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"allowed_runtime_class_names": {
			Type:        schema.TypeList,
			Description: "allowedRuntimeClassNames is a whitelist of RuntimeClass names that may be specified on a pod. A value of \"*\" means that any RuntimeClass name is allowed, and must be the only item in the list. An empty list requires the RuntimeClassName field to be unset.",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"default_runtime_class_name": {
			Type:        schema.TypeString,
			Description: "defaultRuntimeClassName is the default RuntimeClassName to set on the pod. The default MUST be allowed by the allowedRuntimeClassNames list. A value of nil does not mutate the Pod.",
			Optional:    true,
		},
	}

	return s
}
