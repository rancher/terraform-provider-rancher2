package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenLoggingSyslogConfig(in *managementClient.SyslogConfig, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

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

func expandLoggingSyslogConfig(p []interface{}) (*managementClient.SyslogConfig, error) {
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
