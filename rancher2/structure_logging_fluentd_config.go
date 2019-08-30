package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenLoggingFluentdConfigFluentServer(in []managementClient.FluentServer, p []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(in))

	lenP := len(p)

	for i := range in {
		var obj map[string]interface{}
		if lenP <= i {
			obj = make(map[string]interface{})
		} else {
			obj = p[i].(map[string]interface{})
		}

		if len(in[i].Endpoint) > 0 {
			obj["endpoint"] = in[i].Endpoint
		}

		if len(in[i].Hostname) > 0 {
			obj["hostname"] = in[i].Hostname
		}

		if len(in[i].Password) > 0 {
			obj["password"] = in[i].Password
		}

		if len(in[i].SharedKey) > 0 {
			obj["shared_key"] = in[i].SharedKey
		}

		obj["standby"] = in[i].Standby

		if len(in[i].Username) > 0 {
			obj["username"] = in[i].Username
		}

		if in[i].Weight > 0 {
			obj["weight"] = int(in[i].Weight)
		}

		result[i] = obj
	}

	return result, nil
}

func flattenLoggingFluentdConfig(in *managementClient.FluentForwarderConfig, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.FluentServers != nil {
		v, ok := obj["fluent_servers"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		servers, err := flattenLoggingFluentdConfigFluentServer(in.FluentServers, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["fluent_servers"] = servers
	}

	if len(in.Certificate) > 0 {
		obj["certificate"] = in.Certificate
	}

	obj["compress"] = in.Compress

	obj["enable_tls"] = in.EnableTLS

	return []interface{}{obj}, nil
}

// Expanders

func expandLoggingFluentdConfigFluentServer(p []interface{}) ([]managementClient.FluentServer, error) {
	result := []managementClient.FluentServer{}

	if len(p) == 0 || p[0] == nil {
		return result, nil
	}

	for i := range p {
		obj := managementClient.FluentServer{}
		in := p[i].(map[string]interface{})

		if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
			obj.Endpoint = v
		}

		if v, ok := in["hostname"].(string); ok && len(v) > 0 {
			obj.Hostname = v
		}

		if v, ok := in["password"].(string); ok && len(v) > 0 {
			obj.Password = v
		}

		if v, ok := in["shared_key"].(string); ok && len(v) > 0 {
			obj.SharedKey = v
		}

		if v, ok := in["standby"].(bool); ok {
			obj.Standby = v
		}

		if v, ok := in["username"].(string); ok && len(v) > 0 {
			obj.Username = v
		}

		if v, ok := in["weight"].(int); ok {
			obj.Weight = int64(v)
		}

		result = append(result, obj)
	}
	return result, nil
}

func expandLoggingFluentdConfig(p []interface{}) (*managementClient.FluentForwarderConfig, error) {
	obj := &managementClient.FluentForwarderConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["fluent_servers"].([]interface{}); ok && len(v) > 0 {
		servers, err := expandLoggingFluentdConfigFluentServer(v)
		if err != nil {
			return obj, err
		}
		obj.FluentServers = servers
	}

	if v, ok := in["certificate"].(string); ok && len(v) > 0 {
		obj.Certificate = v
	}

	if v, ok := in["compress"].(bool); ok {
		obj.Compress = v
	}

	if v, ok := in["enable_tls"].(bool); ok {
		obj.EnableTLS = v
	}

	return obj, nil
}
