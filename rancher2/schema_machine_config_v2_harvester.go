package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func machineConfigV2HarvesterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"vm_namespace": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Virtual machine namespace",
		},
		"vm_affinity": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "VM affinity, base64 is supported",
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
			Type:     schema.TypeString,
			Optional: true,
			Default:  "0",
			ConflictsWith: []string{
				"harvester_config.0.disk_info",
			},
			Description: "Disk size (in GiB)",
			Deprecated:  "Use disk_info instead",
		},
		"disk_bus": {
			Type:     schema.TypeString,
			Optional: true,
			ConflictsWith: []string{
				"harvester_config.0.disk_info",
			},
			Description: "Disk bus",
			Deprecated:  "Use disk_info instead",
		},
		"image_name": {
			Type:     schema.TypeString,
			Optional: true,
			ConflictsWith: []string{
				"harvester_config.0.disk_info",
			},
			Description: "Image name",
			Deprecated:  "Use disk_info instead",
		},
		"disk_info": {
			Type:     schema.TypeString,
			Optional: true,
			AtLeastOneOf: []string{
				"harvester_config.0.image_name",
			},
			ConflictsWith: []string{
				"harvester_config.0.disk_size",
				"harvester_config.0.disk_bus",
				"harvester_config.0.image_name",
			},
			Description: "A JSON string specifying info for the disks e.g. `{\"disks\":[{\"imageName\":\"harvester-public/image-57hzg\",\"bootOrder\":1,\"size\":40},{\"storageClassName\":\"node-driver-test\",\"bootOrder\":2,\"size\":1}]}`",
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
			Type:     schema.TypeString,
			Optional: true,
			ConflictsWith: []string{
				"harvester_config.0.network_info",
			},
			Description: "Network name",
			Deprecated:  "Use network_info instead",
		},
		"network_model": {
			Type:     schema.TypeString,
			Optional: true,
			ConflictsWith: []string{
				"harvester_config.0.network_info",
			},
			Description: "Network model",
			Deprecated:  "Use network_info instead",
		},
		"network_info": {
			Type:     schema.TypeString,
			Optional: true,
			AtLeastOneOf: []string{
				"harvester_config.0.network_name",
			},
			ConflictsWith: []string{
				"harvester_config.0.network_name",
				"harvester_config.0.network_model",
			},
			Description: "A JSON string specifying info for the networks e.g. `{\"interfaces\":[{\"networkName\":\"harvester-public/vlan1\"},{\"networkName\":\"harvester-public/vlan2\"}]}`",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "UserData content of cloud-init, base64 is supported. If the image does not contain the qemu-guest-agent package, you must install and start qemu-guest-agent using userdata",
		},
		"network_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "NetworkData content of cloud-init, base64 is supported",
		},
		"enable_efi": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable EFI mode",
		},
		"enable_secure_boot": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable secure boot, only available when enable_efi is true",
		},
	}

	return s
}
