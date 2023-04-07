package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	harvesterConfigDriver = "harvester"

	harvesterDiskBusVIRTIO = "virtio"
	harvesterDiskBusSCSI   = "scsi"
	harvesterDiskBusSATA   = "sata"

	harvesterNetworkModelVIRTIO  = "virtio"
	harvesterNetworkModelE1000   = "e1000"
	harvesterNetworkModelE1000E  = "e1000e"
	harvesterNetworkModelNE2KPCO = "ne2k_pco"
	harvesterNetworkModelPCNET   = "pcnet"
	harvesterNetworkModelRTL8139 = "rtl8139"
)

//Types

type harvesterConfig struct {
	VMNamespace  string `json:"vmNamespace,omitempty" yaml:"vmNamespace,omitempty"`
	CPUCount     string `json:"cpuCount,omitempty" yaml:"cpuCount,omitempty"`
	MemorySize   string `json:"memorySize,omitempty" yaml:"memorySize,omitempty"`
	DiskSize     string `json:"diskSize,omitempty" yaml:"diskSize,omitempty"`
	DiskBus      string `json:"diskBus,omitempty" yaml:"diskBus,omitempty"`
	ImageName    string `json:"imageName,omitempty" yaml:"imageName,omitempty"`
	SSHUser      string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	SSHPassword  string `json:"sshPassword,omitempty" yaml:"sshPassword,omitempty"`
	NetworkName  string `json:"networkName,omitempty" yaml:"networkName,omitempty"`
	NetworkModel string `json:"networkModel,omitempty" yaml:"networkModel,omitempty"`
	UserData     string `json:"userData,omitempty" yaml:"userData,omitempty"`
	NetworkData  string `json:"networkData,omitempty" yaml:"networkData,omitempty"`
}

//Schemas

func harvesterConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"vm_namespace": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Virtual machine namespace",
		},
		"cpu_count": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "2",
			Description: "CPU count",
		},
		"memory_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "4",
			Description: "Memory size (in GiB)",
		},
		"disk_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "40",
			Description: "Disk size (in GiB)",
		},
		"disk_bus": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  harvesterDiskBusVIRTIO,
			ValidateFunc: validation.StringInSlice([]string{
				harvesterDiskBusVIRTIO,
				harvesterDiskBusSATA,
				harvesterDiskBusSCSI,
			}, false),
			Description: "Disk bus",
		},
		"image_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Image name",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SSH username",
		},
		"ssh_password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SSH password",
		},
		"network_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Network name",
		},
		"network_model": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  harvesterNetworkModelVIRTIO,
			ValidateFunc: validation.StringInSlice([]string{
				harvesterNetworkModelVIRTIO,
				harvesterNetworkModelE1000,
				harvesterNetworkModelE1000E,
				harvesterNetworkModelNE2KPCO,
				harvesterNetworkModelPCNET,
				harvesterNetworkModelRTL8139,
			}, false),
			Description: "Network model",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "UserData content of cloud-init, base64 is supported",
		},
		"network_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "NetworkData content of cloud-init, base64 is supported",
		},
	}

	return s
}
