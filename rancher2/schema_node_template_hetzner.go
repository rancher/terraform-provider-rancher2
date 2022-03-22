package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	hetznerConfigDriver = "hetzner"
)

//Types

type hetznerConfig struct {
	APIToken          string            `json:"apiToken,omitempty" yaml:"apiToken,omitempty"`
	Image             string            `json:"image,omitempty" yaml:"image,omitempty"`
	ImageID           string            `json:"imageId,omitempty" yaml:"imageId,omitempty"`
	ServerLabels      map[string]string `json:"serverLabels,omitempty" yaml:"serverLabels,omitempty"`
	ServerLocation    string            `json:"serverLocation,omitempty" yaml:"serverLocation,omitempty"`
	ServerType        string            `json:"serverType,omitempty" yaml:"serverType,omitempty"`
	Networks          []string          `json:"networks,omitempty" yaml:"networks,omitempty"`
	UsePrivateNetwork bool              `json:"usePrivateNetwork,omitempty" yaml:"usePrivateNetwork,omitempty"`
	UserData          string            `json:"userData,omitempty" yaml:"userData,omitempty"`
	Volumes           []string          `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Firewalls         []string          `json:"firewalls,omitempty" yaml:"firewalls,omitempty"`
	AdditionalKeys    []string          `json:"additionalKey,omitempty" yaml:"additionalKey,omitempty"`
	PlacementGroup    string            `json:"placementGroup,omitempty" yaml:"placementGroup,omitempty"`
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
		"image_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "15512617",
			Description: "Hetzner Cloud server image id",
		},
		"server_labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Map of the labels which will be assigned to the server",
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
		"use_private_network": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use private network",
		},
		"userdata": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path to file with cloud-init user-data",
		},
		"volumes": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comma-separated list of volume IDs or names which should be attached to the server",
		},
		"firewalls": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "List of firewall IDs or name which should be attached to the server",
		},
		"additional_keys": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "List of ssh keys which should be used to provision the machine with",
		},
		"placement_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Placement group string",
		},
	}

	return s
}
