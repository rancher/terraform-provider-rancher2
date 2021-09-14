package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func machineConfigV2DigitaloceanFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Digital Ocean access token",
		},
		"backups": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable backups for droplet",
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ubuntu-16-04-x64",
			Description: "Digital Ocean Image",
		},
		"ipv6": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable ipv6 for droplet",
		},
		"monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable monitoring for droplet",
		},
		"private_networking": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable private networking for droplet",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "nyc3",
			Description: "Digital Ocean region",
		},
		"size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "s-1vcpu-1gb",
			Description: "Digital Ocean size",
		},
		"ssh_key_contents": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SSH private key contents",
		},
		"ssh_key_fingerprint": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SSH key fingerprint",
		},
		"ssh_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "22",
			Description: "SSH port",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "root",
			Description: "SSH username",
		},
		"tags": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comma-separated list of tags to apply to the Droplet",
		},
		"userdata": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-user",
			Description: "Path to file with cloud-init user-data",
		},
	}

	return s
}
