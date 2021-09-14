package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	machineConfigV2VmwarevsphereCreationTypeDefault = "legacy"
)

var (
	machineConfigV2VmwarevsphereCreationType             = []string{"vm", "template", "library", "legacy"}
	machineConfigV2VmwarevsphereVappIpallocationpolicies = []string{"dhcp", "fixed", "transient", "fixedAllocated"}
	machineConfigV2VmwarevsphereVappIpprotocols          = []string{"IPv4", "IPv6"}
	machineConfigV2VmwarevsphereVappTransports           = []string{"iso", "com.vmware.guestInfo"}
)

//Schemas

func machineConfigV2VmwarevsphereFields() map[string]*schema.Schema {
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
			Default:      machineConfigV2VmwarevsphereCreationTypeDefault,
			ValidateFunc: validation.StringInSlice(machineConfigV2VmwarevsphereCreationType, true),
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
			ValidateFunc: validation.StringInSlice(machineConfigV2VmwarevsphereVappIpallocationpolicies, true),
		},
		"vapp_ip_protocol": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "vSphere vApp IP protocol for this deployment. Supported values are: IPv4 and IPv6",
			ValidateFunc: validation.StringInSlice(machineConfigV2VmwarevsphereVappIpprotocols, true),
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
			ValidateFunc: validation.StringInSlice(machineConfigV2VmwarevsphereVappTransports, true),
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
