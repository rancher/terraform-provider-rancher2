package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	loggingCustomTargetKind = "custom"
)

func loggingCustomTargetConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"content": {
			Type:      schema.TypeString,
			Required:  true,
			StateFunc: TrimSpace,
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			StateFunc: TrimSpace,
		},
		"client_cert": {
			Type:      schema.TypeString,
			Optional:  true,
			StateFunc: TrimSpace,
		},
		"client_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
			StateFunc: TrimSpace,
		},
	}

	return s
}
