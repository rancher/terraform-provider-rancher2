package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	digitaloceanConfigDriver = "digitalocean"
)

//Types

type digitaloceanConfig struct {
	AccessToken       string `json:"accessToken,omitempty" yaml:"accessToken,omitempty"`
	Backups           bool   `json:"backups,omitempty" yaml:"backups,omitempty"`
	Image             string `json:"image,omitempty" yaml:"image,omitempty"`
	IPV6              bool   `json:"ipv6,omitempty" yaml:"ipv6,omitempty"`
	Monitoring        bool   `json:"monitoring,omitempty" yaml:"monitoring,omitempty"`
	PrivateNetworking bool   `json:"privateNetworking,omitempty" yaml:"privateNetworking,omitempty"`
	Region            string `json:"region,omitempty" yaml:"region,omitempty"`
	Size              string `json:"size,omitempty" yaml:"size,omitempty"`
	SSHKeyFingerprint string `json:"sshKeyFingerprint,omitempty" yaml:"sshKeyFingerprint,omitempty"`
	SSHKeyPath        string `json:"sshKeyPath,omitempty" yaml:"sshKeyPath,omitempty"`
	SSHPort           string `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser           string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	Tags              string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Userdata          string `json:"userdata,omitempty" yaml:"userdata,omitempty"`
}

//Schemas

func digitaloceanConfigFields() map[string]*schema.Schema {
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
		"ssh_key_fingerprint": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SSH key fingerprint",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH private key path",
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
