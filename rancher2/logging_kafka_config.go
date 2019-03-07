package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	loggingKafkaKind = "kafka"
)

//Schemas

func kafkaConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"topic": {
			Type:     schema.TypeString,
			Required: true,
		},
		"broker_endpoints": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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
		"zookeeper_endpoint": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}

// Flatteners

func flattenKafkaConfig(in *managementClient.KafkaConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
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

func expandKafkaConfig(p []interface{}) (*managementClient.KafkaConfig, error) {
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
