package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	opennebulaConfigDriver = "opennebula"
)

//Types

type opennebulaConfig struct {
	b2dsize      string `json:"b2dSize,omitempty" yaml:"b2dSize,omitempty"`
	devprefix    string `json:"devPrefix,omitempty" yaml:"devPrefix,omitempty"`
	diskresize   string `json:"diskResize,omitempty" yaml:"diskResize,omitempty"`
	imagename    string `json:"imageName,omitempty" yaml:"imageName,omitempty"`
	memory       string `json:"memory,omitempty" yaml:"memory,omitempty"`
	networkname  string `json:"networkName,omitempty" yaml:"networkName,omitempty"`
	password     string `json:"password,omitempty" yaml:"password,omitempty"`
	templateid   string `json:"templateId,omitempty" yaml:"templateId,omitempty"`
	user         string `json:"user,omitempty" yaml:"user,omitempty"`
	xmlrpcurl    string `json:"xmlrpcurl,omitempty" yaml:"xmlrpcurl,omitempty"`
	cpu          string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	disablevnc   bool   `json:"disableVnc,omitempty" yaml:"disableVnc,omitempty"`
	imageid      string `json:"imageId,omitempty" yaml:"imageId,omitempty"`
	imageowner   string `json:"imageOwner,omitempty" yaml:"imageOwner,omitempty"`
	networkid    string `json:"networkId,omitempty" yaml:"networkId,omitempty"`
	networkowner string `json:"networkOwner,omitempty" yaml:"networkOwner,omitempty"`
	sshuser      string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	templatename string `json:"templateName,omitempty" yaml:"templateName,omitempty"`
	vcpu         string `json:"vcpu,omitempty" yaml:"vcpu,omitempty"`
}

//Schemas

func opennebulaConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"b2dsize": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"devprefix": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"diskresize": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"imagename": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"memory": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"networkname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"templateid": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"user": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"xmlrpcurl": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"cpu": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"disablevnc": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"imageid": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"imageowner": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"networkid": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"networkowner": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  false,
		},
		"sshuser": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "docker",
		},
		"templatename": &schema.Schema{
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
