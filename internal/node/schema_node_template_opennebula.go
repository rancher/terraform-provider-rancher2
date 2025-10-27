package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	opennebulaConfigDriver = "opennebula"
)

//Types

type opennebulaConfig struct {
	B2dSize      string `json:"b2dSize,omitempty" yaml:"b2dSize,omitempty"`
	CPU          string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	DevPrefix    string `json:"devPrefix,omitempty" yaml:"devPrefix,omitempty"`
	DiskResize   string `json:"diskResize,omitempty" yaml:"diskResize,omitempty"`
	DisableVnc   bool   `json:"disableVnc,omitempty" yaml:"disableVnc,omitempty"`
	ImageName    string `json:"imageName,omitempty" yaml:"imageName,omitempty"`
	ImageID      string `json:"imageId,omitempty" yaml:"imageId,omitempty"`
	ImageOwner   string `json:"imageOwner,omitempty" yaml:"imageOwner,omitempty"`
	Memory       string `json:"memory,omitempty" yaml:"memory,omitempty"`
	NetworkID    string `json:"networkId,omitempty" yaml:"networkId,omitempty"`
	NetworkName  string `json:"networkName,omitempty" yaml:"networkName,omitempty"`
	NetworkOwner string `json:"networkOwner,omitempty" yaml:"networkOwner,omitempty"`
	Password     string `json:"password,omitempty" yaml:"password,omitempty"`
	SSHUser      string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	TemplateID   string `json:"templateId,omitempty" yaml:"templateId,omitempty"`
	TemplateName string `json:"templateName,omitempty" yaml:"templateName,omitempty"`
	User         string `json:"user,omitempty" yaml:"user,omitempty"`
	Vcpu         string `json:"vcpu,omitempty" yaml:"vcpu,omitempty"`
	XMLRPCURL    string `json:"xmlrpcurl,omitempty" yaml:"xmlrpcurl,omitempty"`
}

//Schemas

func opennebulaConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"user": {
			Type:     schema.TypeString,
			Required: true,
		},
		"xml_rpc_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"b2d_size": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dev_prefix": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"disable_vnc": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"disk_resize": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_owner": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"memory": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_owner": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"template_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"template_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ssh_user": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "docker",
		},
		"vcpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}
