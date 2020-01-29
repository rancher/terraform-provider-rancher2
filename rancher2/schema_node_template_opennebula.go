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
	DevPrefix    string `json:"devPrefix,omitempty" yaml:"devPrefix,omitempty"`
	DiskResize   string `json:"diskResize,omitempty" yaml:"diskResize,omitempty"`
	ImageName    string `json:"imageName,omitempty" yaml:"imageName,omitempty"`
	Memory       string `json:"memory,omitempty" yaml:"memory,omitempty"`
	NetworkName  string `json:"networkName,omitempty" yaml:"networkName,omitempty"`
	Password     string `json:"password,omitempty" yaml:"password,omitempty"`
	TemplateID   string `json:"templateId,omitempty" yaml:"templateId,omitempty"`
	User         string `json:"user,omitempty" yaml:"user,omitempty"`
	XMLRPCURL    string `json:"xmlrpcurl,omitempty" yaml:"xmlrpcurl,omitempty"`
	CPU          string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	DisableVnc   bool   `json:"disableVnc,omitempty" yaml:"disableVnc,omitempty"`
	ImageID      string `json:"imageId,omitempty" yaml:"imageId,omitempty"`
	ImageOwner   string `json:"imageOwner,omitempty" yaml:"imageOwner,omitempty"`
	NetworkID    string `json:"networkId,omitempty" yaml:"networkId,omitempty"`
	NetworkOwner string `json:"networkOwner,omitempty" yaml:"networkOwner,omitempty"`
	SSHUser      string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	TemplateName string `json:"templateName,omitempty" yaml:"templateName,omitempty"`
	Vcpu         string `json:"vcpu,omitempty" yaml:"vcpu,omitempty"`
}

//Schemas

func opennebulaConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"b2d_size": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"dev_prefix": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"disk_resize": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"memory": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"template_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"user": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"xml_rpc_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"cpu": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"disable_vnc": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"image_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_owner": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_owner": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"ssh_user": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "docker",
		},
		"template_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"vcpu": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}
