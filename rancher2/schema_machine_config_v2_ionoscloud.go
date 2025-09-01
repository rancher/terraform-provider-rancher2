package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ionoscloudConfigDriver = "ionoscloud"
)

//Schemas

func machineConfigV2IonoscloudFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"nat_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Ionos Cloud NAT Gateway name. Note that setting this will NOT implicitly create a NAT, this flag will only be read if need be",
		},
		"nat_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Ionos Cloud existing and configured NAT Gateway",
		},
		"create_nat": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If set, will create a default NAT. Requires private LAN",
		},
		"nat_public_ips": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ionos Cloud NAT Gateway public IPs",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"nat_flowlogs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ionos Cloud NAT Gateway Flowlogs",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"nat_rules": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ionos Cloud NAT Gateway Rules",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"skip_default_nat_rules": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Should the driver skip creating default nat rules if creating a NAT, creating only the specified rules",
		},
		"nat_lans_to_gateways": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Ionos Cloud NAT map of LANs to a slice of their Gateway IPs. Example: \"1=10.0.0.1,10.0.0.2:2=10.0.0.10\"",
		},
		"private_lan": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Should the created LAN be private? Does nothing if LAN ID is provided",
		},
		"additional_lans": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Names of existing IONOS Lans to connect the machine to. Names that are not found are ignored",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"wait_for_ip_change": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Should the driver wait for the NIC IP to be set by external sources?",
		},
		"wait_for_ip_change_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     600,
			Description: "Timeout used when waiting for NIC IP changes",
		},
		"nic_dhcp": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Should the created NIC have DHCP set to true or false?",
		},
		"nic_ips": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ionos Cloud NIC IPs",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "https://api.ionos.com/cloudapi/v6",
			Description: "Ionos Cloud API Endpoint",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Ionos Cloud Username",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Ionos Cloud Password",
		},
		"token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Ionos Cloud Token",
		},
		"server_cores": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     2,
			Description: "Ionos Cloud Server Cores (2, 3, 4, 5, 6, etc.)",
		},
		"server_ram": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     2048,
			Description: "Ionos Cloud Server Ram in MB(1024, 2048, 3072, 4096, etc.)",
		},
		"disk_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     50,
			Description: "Ionos Cloud Volume Disk-Size in GB(10, 50, 100, 200, 400)",
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ubuntu:latest",
			Description: "Ionos Cloud Image Id or Alias (ubuntu:latest, debian:latest, etc.)",
		},
		"image_password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Ionos Cloud Image Password to be able to access the server from DCD platform",
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "us/las",
			Description: "Ionos Cloud Location",
		},
		"disk_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "HDD",
			Description: "Ionos Cloud Volume Disk-Type (HDD, SSD)",
		},
		"server_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ENTERPRISE",
			Description: "Ionos Cloud Server Type(ENTERPRISE or CUBE). CUBE servers are only available in certain locations.",
		},
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Basic Cube XS",
			Description: "Ionos Cloud CUBE Template, only used for CUBE servers.",
		},
		"server_cpu_family": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Ionos Cloud Server CPU families (INTEL_XEON, INTEL_SKYLAKE, INTEL_ICELAKE, AMD_EPYC, INTEL_SIERRAFOREST)",
		},
		"datacenter_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Ionos Cloud Virtual Data Center Id",
		},
		"datacenter_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-machine-data-center",
			Description: "Ionos Cloud Virtual Data Center Name",
		},
		"lan_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Ionos Cloud LAN Id",
		},
		"lan_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-machine-lan",
			Description: "Ionos Cloud LAN Name",
		},
		"volume_availability_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "AUTO",
			Description: "Ionos Cloud Volume Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)",
		},
		"server_availability_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "AUTO",
			Description: "Ionos Cloud Server Availability Zone (AUTO, ZONE_1, ZONE_2)",
		},
		"cloud_init": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cloud-init configuration for the volume as a multi-line string",
		},
		"cloud_init_b64": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cloud-init configuration for the volume as a base64 encoded string",
		},
		"append_rke_userdata": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Should the driver append the rke user-data to the user-data sent to the ionos server",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "root",
			Description: "The name of the user the driver will use for ssh",
		},
		"ssh_in_cloud_init": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Should the driver only add the SSH info in the user data? (required for custom images)",
		},
	}

	return s
}
