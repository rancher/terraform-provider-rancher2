package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type vmwarevsphereCredentialConfig struct {
	Password    string `json:"password,omitempty" yaml:"password,omitempty"`
	Username    string `json:"username,omitempty" yaml:"username,omitempty"`
	Vcenter     string `json:"vcenter,omitempty" yaml:"vcenter,omitempty"`
	VcenterPort string `json:"vcenterPort,omitempty" yaml:"vcenterPort,omitempty"`
}

//Schemas

func cloudCredentialVsphereFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "vSphere password",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "vSphere username",
		},
		"vcenter": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "vSphere IP/hostname for vCenter",
		},
		"vcenter_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "443",
			Description: "vSphere Port for vCenter",
		},
	}

	return s
}
