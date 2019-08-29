package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenLoggingKafkaConfig(in *managementClient.KafkaConfig, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	obj["topic"] = in.Topic

	if len(in.BrokerEndpoints) > 0 {
		obj["broker_endpoints"] = toArrayInterface(in.BrokerEndpoints)
	}

	if len(in.Certificate) > 0 {
		obj["certificate"] = in.Certificate
	}

	if len(in.ClientCert) > 0 {
		obj["client_cert"] = in.ClientCert
	}

	if len(in.ClientKey) > 0 {
		obj["client_key"] = in.ClientKey
	}

	if len(in.ZookeeperEndpoint) > 0 {
		obj["zookeeper_endpoint"] = in.ZookeeperEndpoint
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandLoggingKafkaConfig(p []interface{}) (*managementClient.KafkaConfig, error) {
	obj := &managementClient.KafkaConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["topic"].(string); ok && len(v) > 0 {
		obj.Topic = v
	}

	if v, ok := in["broker_endpoints"].([]interface{}); ok && len(v) > 0 {
		obj.BrokerEndpoints = toArrayString(v)
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

	if v, ok := in["zookeeper_endpoint"].(string); ok && len(v) > 0 {
		obj.ZookeeperEndpoint = v
	}

	return obj, nil
}
