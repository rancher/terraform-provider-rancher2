package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	loggingSyslogKind = "syslog"
)

var (
	syslogSeverityKinds = []string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}
	syslogProtocolKinds = []string{"tcp", "udp"}
)

//Schemas

func loggingSyslogConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_cert": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"program": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"protocol": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "udp",
			ValidateFunc: validation.StringInSlice(syslogProtocolKinds, true),
		},
		"severity": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "notice",
			ValidateFunc: validation.StringInSlice(syslogSeverityKinds, true),
		},
		"ssl_verify": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"token": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}

	return s
}
