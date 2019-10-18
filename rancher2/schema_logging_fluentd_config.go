package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	loggingFluentdKind = "fluentd"
)

//Schemas

func loggingFluentdConfigFluentServerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"shared_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"standby": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"username": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"weight": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	return s
}

func loggingFluentdConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"fluent_servers": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: loggingFluentdConfigFluentServerFields(),
			},
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"compress": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"enable_tls": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	return s
}
