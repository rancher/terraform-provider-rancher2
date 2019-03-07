package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	loggingSyslogKind = "syslog"
)

var (
	syslogSeverityKinds = []string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}
	syslogProtocolKinds = []string{"tcp", "udp"}
)

//Schemas

func syslogConfigFields() map[string]*schema.Schema {
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

// Flatteners

func flattenSyslogConfig(in *managementClient.SyslogConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["endpoint"] = in.Endpoint

	if len(in.Certificate) > 0 {
		obj["certificate"] = in.Certificate
	}

	if len(in.ClientCert) > 0 {
		obj["client_cert"] = in.ClientCert
	}

	if len(in.ClientKey) > 0 {
		obj["client_key"] = in.ClientKey
	}

	if len(in.Program) > 0 {
		obj["program"] = in.Program
	}

	if len(in.Protocol) > 0 {
		obj["protocol"] = in.Protocol
	}

	if len(in.Severity) > 0 {
		obj["severity"] = in.Severity
	}

	obj["ssl_verify"] = in.SSLVerify

	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandSyslogConfig(p []interface{}) (*managementClient.SyslogConfig, error) {
	obj := &managementClient.SyslogConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["certificate"].(string); ok && len(v) > 0 {
		obj.Certificate = v
	}

	if v, ok := in["client_cert"].(string); ok && len(v) > 0 {
		obj.ClientCert = v
	}

	if v, ok := in["client_key"].(string); ok && len(v) > 0 {
		obj.ClientKey = v
	}

	if v, ok := in["program"].(string); ok && len(v) > 0 {
		obj.Program = v
	}

	if v, ok := in["protocol"].(string); ok && len(v) > 0 {
		obj.Protocol = v
	}

	if v, ok := in["severity"].(string); ok && len(v) > 0 {
		obj.Severity = v
	}

	if v, ok := in["ssl_verify"].(bool); ok {
		obj.SSLVerify = v
	}

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	return obj, nil
}
