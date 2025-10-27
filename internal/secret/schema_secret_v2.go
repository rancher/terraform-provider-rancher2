package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	"k8s.io/api/core/v1"
)

const (
	secretV2Kind         = "Secret"
	secretV2APIVersion   = "v1"
	secretV2APIType      = "secret"
	secretV2ClusterIDsep = "."
)

//Types

type SecretV2 struct {
	norman.Resource
	v1.Secret
	K8SType string `json:"_type,omitempty"`
}

func secretV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "K8s cluster ID",
		},
		"data": {
			Type:        schema.TypeMap,
			Required:    true,
			Sensitive:   true,
			Description: "Secret data map",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "K8s Secret name",
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "default",
			Description: "K8s Secret namespace",
		},
		"immutable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If set to true, ensures that data stored in the Secret cannot be updated",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     string(v1.SecretTypeOpaque),
			Description: "Used to facilitate programmatic handling of Secret data",
		},
		"resource_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
