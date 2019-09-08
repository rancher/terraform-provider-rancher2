package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	hetznerConfigDriver = "hetzner"
)

//Types

type hetznerConfig struct {
	APIToken       string `json:"apiToken,omitempty" yaml:"apiToken,omitempty"`
	Image          string `json:"image,omitempty" yaml:"image,omitempty"`
	ServerLocation string `json:"serverLocation,omitempty" yaml:"serverLocation,omitempty"`
	ServerType     string `json:"serverType,omitempty" yaml:"serverType,omitempty"`
	Networks       string `json:"networks,omitempty" yaml:"networks,omitempty"`
	Volumes        string `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Userdata       string `json:"userData,omitempty" yaml:"userData,omitempty"`
}

//Schemas

func hetznerConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"api_token": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Hetzner Cloud project API token",
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ubuntu-18.04",
			Description: "Hetzner Cloud server image",
		},
		"server_location": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "nbg1",
			Description: "Hetzner Cloud datacenter",
		},
		"server_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "cx11",
			Description: "Hetzner Cloud server type",
		},
		"networks": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comma-separated list of network IDs or names which should be attached to the server private network interface",
		},
		"volumes": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comma-separated list of volume IDs or names which should be attached to the server",
		},
		"userdata": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path to file with cloud-init user-data",
		},
	}

	return s
}
