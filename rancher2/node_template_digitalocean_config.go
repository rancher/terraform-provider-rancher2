package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	digitaloceanConfigDriver = "digitalocean"
)

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
			Required:    true,
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

// Flatteners

func flattenDigitaloceanConfig(in *digitaloceanConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessToken) > 0 {
		obj["access_token"] = in.AccessToken
	}

	obj["backups"] = in.Backups

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["ipv6"] = in.IPV6
	obj["monitoring"] = in.Monitoring
	obj["private_networking"] = in.PrivateNetworking

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.Size) > 0 {
		obj["size"] = in.Size
	}

	if len(in.SSHKeyFingerprint) > 0 {
		obj["ssh_key_fingerprint"] = in.SSHKeyFingerprint
	}

	if len(in.SSHKeyPath) > 0 {
		obj["ssh_key_path"] = in.SSHKeyPath
	}

	if len(in.SSHPort) > 0 {
		obj["ssh_port"] = in.SSHPort
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if len(in.Userdata) > 0 {
		obj["userdata"] = in.Userdata
	}

	return []interface{}{obj}
}

// Expanders

func expandDigitaloceanConfig(p []interface{}) *digitaloceanConfig {
	obj := &digitaloceanConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_token"].(string); ok && len(v) > 0 {
		obj.AccessToken = v
	}

	if v, ok := in["backups"].(bool); ok {
		obj.Backups = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["ipv6"].(bool); ok {
		obj.IPV6 = v
	}

	if v, ok := in["monitoring"].(bool); ok {
		obj.Monitoring = v
	}
	if v, ok := in["private_networking"].(bool); ok {
		obj.PrivateNetworking = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["size"].(string); ok && len(v) > 0 {
		obj.Size = v
	}

	if v, ok := in["ssh_key_fingerprint"].(string); ok && len(v) > 0 {
		obj.SSHKeyFingerprint = v
	}

	if v, ok := in["ssh_key_path"].(string); ok && len(v) > 0 {
		obj.SSHKeyPath = v
	}

	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.Userdata = v
	}

	return obj
}
