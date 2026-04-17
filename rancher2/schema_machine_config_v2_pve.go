package rancher2

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func machineConfigV2PveFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pve_url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Proxmox VE URL (e.g. 'https://<PROXMOX VE ADDRESS>:8006')",
		},
		"pve_token_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Proxmox VE API Token ID (including username and realm, e.g. 'root@pam!rancher')",
		},
		"pve_token_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Proxmox VE API Token secret",
		},
		"pve_insecure_tls": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disables Proxmox VE TLS certificate verification",
		},
		"pve_resource_pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Proxmox VE Resource Pool name",
		},
		"pve_template_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the Proxmox VE template",
		},
		"pve_iso_device": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Bus/Device of the CD/DVD Drive to mount cloud-init ISO to (e.g. 'scsi1')",
		},
		"pve_network_interface": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Bus/Device of the network interface to read machine's IP address from (e.g. 'net0')",
		},
		"pve_ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ubuntu",
			Description: "Username for the SSH user created via cloud-init",
		},
		"pve_ssh_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     22,
			Description: "Port to use when connecting to the machine via SSH",
		},
		"pve_processor_sockets": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Number of processor sockets to configure for the machine",
		},
		"pve_processor_cores": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Number of processor cores to configure for the machine",
		},
		"pve_memory": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Amount of memory in MiB to configure for the machine",
		},
	}
}
