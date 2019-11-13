package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	loggingKafkaKind = "kafka"
)

//Schemas

func loggingKafkaConfigFields() map[string]*schema.Schema {
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
