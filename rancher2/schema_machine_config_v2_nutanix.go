package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	machineConfigV2NutanixBootTypes = []string{"legacy", "uefi"}
)

func machineConfigV2NutanixFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Nutanix management endpoint IP address/FQDN",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Nutanix management username",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Nutanix management password",
		},
		"port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "9440",
			Description: "Nutanix management endpoint port",
		},
		"insecure": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Allow insecure SSL requests",
		},
		"cluster": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Nutanix cluster to install VM on",
		},
		"vm_mem": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "2048",
			Description: "Memory in MB of the VM to be created",
		},
		"vm_cpus": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "2",
			Description: "Number of VCPUs of the VM to be created",
		},
		"vm_cores": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "1",
			Description: "Number of cores per VCPU of the VM to be created",
		},
		"vm_cpu_passthrough": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable passthrough of host CPU features to the VM",
		},
		"vm_network": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Network names or UUIDs to attach to the VM",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"vm_image": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the VM disk image to clone from",
		},
		"vm_image_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "0",
			Description: "Increase the size of the template image (GiB)",
		},
		"vm_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Categories to apply to the VM",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"storage_container": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "UUID of the storage container for additional disk",
		},
		"disk_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "0",
			Description: "Size of the additional disk (GiB)",
		},
		"cloud_init": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cloud-init configuration",
		},
		"vm_serial_port": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Attach a serial port to the VM",
		},
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the project to assign the VM",
		},
		"boot_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "legacy",
			ValidateFunc: validation.StringInSlice(machineConfigV2NutanixBootTypes, true),
			Description:  "Boot type of the VM. Supported values: legacy, uefi",
		},
		"timeout": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "300",
			Description: "Timeout for Nutanix operations in seconds",
		},
		"vm_gpu": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "GPU devices to attach to the VM",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
