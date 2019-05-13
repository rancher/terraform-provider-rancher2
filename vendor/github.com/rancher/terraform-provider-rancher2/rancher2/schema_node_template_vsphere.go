package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	vmwarevsphereConfigDriver = "vmwarevsphere"
)

var (
	vmwarevsphereConfigVappIpallocationpolicies = []string{"dhcp", "fixed", "transient", "fixedAllocated"}
	vmwarevsphereConfigVappIpprotocols          = []string{"IPv4", "IPv6"}
	vmwarevsphereConfigVappTransports           = []string{"iso", "com.vmware.guestInfo"}
)

//Types

type vmwarevsphereConfig struct {
	Boot2dockerURL         string   `json:"boot2dockerUrl,omitempty" yaml:"boot2dockerUrl,omitempty"`
	Cfgparam               []string `json:"cfgparam,omitempty" yaml:"cfgparam,omitempty"`
	Cloudinit              string   `json:"cloudinit,omitempty" yaml:"cloudinit,omitempty"`
	CPUCount               string   `json:"cpuCount,omitempty" yaml:"cpuCount,omitempty"`
	Datacenter             string   `json:"datacenter,omitempty" yaml:"datacenter,omitempty"`
	Datastore              string   `json:"datastore,omitempty" yaml:"datastore,omitempty"`
	DiskSize               string   `json:"diskSize,omitempty" yaml:"diskSize,omitempty"`
	Folder                 string   `json:"folder,omitempty" yaml:"folder,omitempty"`
	Hostsystem             string   `json:"hostsystem,omitempty" yaml:"hostsystem,omitempty"`
	MemorySize             string   `json:"memorySize,omitempty" yaml:"memorySize,omitempty"`
	Network                []string `json:"network,omitempty" yaml:"network,omitempty"`
	Password               string   `json:"password,omitempty" yaml:"password,omitempty"`
	Pool                   string   `json:"pool,omitempty" yaml:"pool,omitempty"`
	Username               string   `json:"username,omitempty" yaml:"username,omitempty"`
	VappIpallocationpolicy string   `json:"vappIpallocationpolicy,omitempty" yaml:"vappIpallocationpolicy,omitempty"`
	VappIpprotocol         string   `json:"vappIpprotocol,omitempty" yaml:"vappIpprotocol,omitempty"`
	VappProperty           []string `json:"vappProperty,omitempty" yaml:"vappProperty,omitempty"`
	VappTransport          string   `json:"vappTransport,omitempty" yaml:"vappTransport,omitempty"`
	Vcenter                string   `json:"vcenter,omitempty" yaml:"vcenter,omitempty"`
	VcenterPort            string   `json:"vcenterPort,omitempty" yaml:"vcenterPort,omitempty"`
}

//Schemas

func vsphereConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"boot2docker_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "https://releases.rancher.com/os/latest/rancheros-vmware.iso",
		},
		"cfgparam": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cloudinit": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"cpu_count": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "2",
		},
		"datacenter": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"datastore": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"disk_size": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "20480",
		},
		"folder": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"hostsystem": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"memory_size": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "2048",
		},
		"network": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "vSphere password",
		},
		"pool": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere username",
		},
		"vapp_ip_allocation_policy": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(vmwarevsphereConfigVappIpallocationpolicies, true),
		},
		"vapp_ip_protocol": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(vmwarevsphereConfigVappIpprotocols, true),
		},
		"vapp_property": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"vapp_transport": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(vmwarevsphereConfigVappTransports, true),
		},
		"vcenter": {
			Type:        schema.TypeString,
			Optional:    true,
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
