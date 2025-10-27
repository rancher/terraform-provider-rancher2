package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func machineConfigV2OpenstackFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"auth_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Required: true,
		},
		"region": {
			Type:     schema.TypeString,
			Required: true,
		},
		"username": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"active_timeout": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "200",
		},
		"cacert": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"config_drive": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"domain_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"domain_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"endpoint_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"flavor_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"flavor_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"floating_ip_pool": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"insecure": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"ip_version": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "4",
		},
		"keypair_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"net_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"net_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"nova_network": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"private_key_file": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"sec_groups": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ssh_port": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "22",
		},
		"ssh_user": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "root",
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tenant_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tenant_domain_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tenant_domain_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_domain_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_domain_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_data_file": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"application_credential_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"application_credential_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"application_credential_secret": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"boot_from_volume": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"volume_size": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"volume_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"volume_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"volume_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"volume_device_path": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}
