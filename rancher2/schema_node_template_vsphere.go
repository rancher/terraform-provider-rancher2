package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	vmwarevsphereConfigDriver              = "vmwarevsphere"
	vmwarevsphereConfigCreationTypeDefault = "legacy"
)

var (
	vmwarevsphereConfigCreationType             = []string{"vm", "template", "library", "legacy"}
	vmwarevsphereConfigVappIpallocationpolicies = []string{"dhcp", "fixed", "transient", "fixedAllocated"}
	vmwarevsphereConfigVappIpprotocols          = []string{"IPv4", "IPv6"}
	vmwarevsphereConfigVappTransports           = []string{"iso", "com.vmware.guestInfo"}
)

//Types

type vmwarevsphereConfig struct {
	Boot2dockerURL         string   `json:"boot2dockerUrl,omitempty" yaml:"boot2dockerUrl,omitempty"`
	Cfgparam               []string `json:"cfgparam,omitempty" yaml:"cfgparam,omitempty"`
	CloneFrom              string   `json:"cloneFrom,omitempty" yaml:"cloneFrom,omitempty"`
	CloudConfig            string   `json:"cloudConfig,omitempty" yaml:"cloudConfig,omitempty"`
	Cloudinit              string   `json:"cloudinit,omitempty" yaml:"cloudinit,omitempty"`
	ContentLibrary         string   `json:"contentLibrary,omitempty" yaml:"contentLibrary,omitempty"`
	CPUCount               string   `json:"cpuCount,omitempty" yaml:"cpuCount,omitempty"`
	CreationType           string   `json:"creationType,omitempty" yaml:"creationType,omitempty"`
	CustomAttributes       []string `json:"customAttribute,omitempty" yaml:"customAttribute,omitempty"`
	Datacenter             string   `json:"datacenter,omitempty" yaml:"datacenter,omitempty"`
	Datastore              string   `json:"datastore,omitempty" yaml:"datastore,omitempty"`
	DatastoreCluster       string   `json:"datastoreCluster,omitempty" yaml:"datastoreCluster,omitempty"`
	DiskSize               string   `json:"diskSize,omitempty" yaml:"diskSize,omitempty"`
	Folder                 string   `json:"folder,omitempty" yaml:"folder,omitempty"`
	Hostsystem             string   `json:"hostsystem,omitempty" yaml:"hostsystem,omitempty"`
	MemorySize             string   `json:"memorySize,omitempty" yaml:"memorySize,omitempty"`
	Network                []string `json:"network,omitempty" yaml:"network,omitempty"`
	Password               string   `json:"password,omitempty" yaml:"password,omitempty"`
	Pool                   string   `json:"pool,omitempty" yaml:"pool,omitempty"`
	SSHPassword            string   `json:"sshPassword,omitempty" yaml:"sshPassword,omitempty"`
	SSHPort                string   `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser                string   `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	SSHUserGroup           string   `json:"sshUserGroup,omitempty" yaml:"sshUserGroup,omitempty"`
	Tags                   []string `json:"tag,omitempty" yaml:"tag,omitempty"`
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
		"boot2docker_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "https://releases.rancher.com/os/latest/rancheros-vmware.iso",
			Description: "vSphere URL for boot2docker image",
		},
		"cfgparam": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "vSphere vm configuration parameters (used for guestinfo)",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"clone_from": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "If you choose creation type clone a name of what you want to clone is required",
		},
		"cloud_config": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filepath to a cloud-config yaml file to put into the ISO user-data",
		},
		"cloudinit": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere cloud-init filepath or url to add to guestinfo",
		},
		"content_library": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "If you choose to clone from a content library template specify the name of the library",
		},
		"cpu_count": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "2",
			Description: "vSphere CPU number for docker VM",
		},
		"creation_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      vmwarevsphereConfigCreationTypeDefault,
			ValidateFunc: validation.StringInSlice(vmwarevsphereConfigCreationType, true),
			Description:  "Creation type when creating a new virtual machine. Supported values: vm, template, library, legacy",
		},
		"custom_attributes": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "vSphere custom attributes, format key/value e.g. '200=my custom value'",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"datacenter": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere datacenter for virtual machine",
		},
		"datastore": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere datastore for virtual machine",
		},
		"datastore_cluster": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere datastore cluster for virtual machine",
		},
		"disk_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "20480",
			Description: "vSphere size of disk for docker VM (in MB)",
		},
		"folder": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere folder for the docker VM. This folder must already exist in the datacenter",
		},
		"hostsystem": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere compute resource where the docker VM will be instantiated. This can be omitted if using a cluster with DRS",
		},
		"memory_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "2048",
			Description: "vSphere size of memory for docker VM (in MB)",
		},
		"network": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "vSphere network where the virtual machine will be attached",
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
		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere resource pool for docker VM",
		},
		"ssh_password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Default:     "tcuser",
			Description: "If using a non-B2D image you can specify the ssh password",
		},
		"ssh_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "22",
			Description: "If using a non-B2D image you can specify the ssh port",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker",
			Description: "If using a non-B2D image you can specify the ssh user",
		},
		"ssh_user_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "staff",
			Description: "If using a non-B2D image the uploaded keys will need chown'ed, defaults to staff e.g. docker:staff",
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "vSphere tags id e.g. urn:xxx",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "vSphere username",
		},
		"vapp_ip_allocation_policy": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "vSphere vApp IP allocation policy. Supported values are: dhcp, fixed, transient and fixedAllocated",
			ValidateFunc: validation.StringInSlice(vmwarevsphereConfigVappIpallocationpolicies, true),
		},
		"vapp_ip_protocol": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "vSphere vApp IP protocol for this deployment. Supported values are: IPv4 and IPv6",
			ValidateFunc: validation.StringInSlice(vmwarevsphereConfigVappIpprotocols, true),
		},
		"vapp_property": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "vSphere vApp properties",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"vapp_transport": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "vSphere OVF environment transports to use for properties. Supported values are: iso and com.vmware.guestInfo",
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
