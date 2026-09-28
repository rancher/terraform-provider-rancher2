package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func machineConfigV2OpentelekomcloudFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description: "T-Cloud Public access key. Prefer a cloud credential.",
		},
		"auth_url": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public identity service endpoint.",
		},
		"availability_zone": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public availability zone.",
		},
		"bandwidth_size": {
			Type: schema.TypeString, Optional: true, Default: "100",
			Description: "Elastic IP bandwidth size, encoded as a string.",
		},
		"bandwidth_type": {
			Type: schema.TypeString, Optional: true, Default: "PER",
			Description: "Elastic IP bandwidth share type.",
		},
		"cacert": {
			Type: schema.TypeString, Optional: true,
			Description: "CA bundle option exposed by the driver. No local file is uploaded.",
		},
		"cloud": {
			Type: schema.TypeString, Optional: true,
			Description: "Named cloud in clouds.yaml on the machine provisioner.",
		},
		"domain_id": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public domain ID.",
		},
		"domain_name": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public domain name.",
		},
		"eip": {
			Type: schema.TypeString, Optional: true,
			Description: "Existing elastic IP address.",
		},
		"eip_type": {
			Type: schema.TypeString, Optional: true, Default: "5_bgp",
			Description: "Elastic IP type.",
		},
		"endpoint_type": {
			Type: schema.TypeString, Optional: true, Default: "public",
			Description: "T-Cloud Public endpoint interface type.",
		},
		"flavor_id": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public instance flavor ID.",
		},
		"flavor_name": {
			Type: schema.TypeString, Optional: true, Default: "s3.xlarge.2",
			Description: "T-Cloud Public instance flavor name.",
		},
		"image_id": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public image ID.",
		},
		"image_name": {
			Type: schema.TypeString, Optional: true, Default: "Standard_Ubuntu_24.04_amd64_uefi_latest",
			Description: "T-Cloud Public image name.",
		},
		"ip_version": {
			Type: schema.TypeString, Optional: true, Default: "4",
			ValidateFunc: validation.StringInSlice([]string{"4", "6"}, false),
			Description:  "IP address version, encoded as a string.",
		},
		"keypair_name": {
			Type: schema.TypeString, Optional: true,
			Description: "Existing SSH key pair name.",
		},
		"network_scope": {
			Type: schema.TypeString, Optional: true, Default: "machine",
			ValidateFunc: validation.StringInSlice([]string{"machine", "shared"}, false),
			Description:  "Network ownership scope. Shared requires an existing VPC, subnet and security groups.",
		},
		"password": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description: "T-Cloud Public password. Prefer a cloud credential.",
		},
		"private_key_file": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description: "SSH private key path on the machine provisioner, or PEM content supported by the driver.",
		},
		"project_id": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public project ID.",
		},
		"project_name": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public project name.",
		},
		"region": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public region.",
		},
		"root_volume_size": {
			Type: schema.TypeString, Optional: true, Default: "40",
			Description: "Root volume size in GiB, encoded as a string.",
		},
		"root_volume_type": {
			Type: schema.TypeString, Optional: true, Default: "SSD",
			Description: "T-Cloud Public root volume type.",
		},
		"sec_groups": {
			Type: schema.TypeString, Optional: true,
			Description: "Comma-separated existing security group names.",
		},
		"secret_key": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description: "T-Cloud Public secret key. Prefer a cloud credential.",
		},
		"server_group": {
			Type: schema.TypeString, Optional: true,
			Description: "Server group name.",
		},
		"server_group_id": {
			Type: schema.TypeString, Optional: true,
			Description: "Server group ID.",
		},
		"skip_default_sg": {
			Type: schema.TypeBool, Optional: true, Default: false,
			Description: "Do not create the driver's default security group.",
		},
		"skip_eip": {
			Type: schema.TypeBool, Optional: true, Default: false,
			Description: "Do not create an elastic IP; use the machine's private IP.",
		},
		"ssh_allow_cidr": {
			Type: schema.TypeString, Optional: true,
			Description: "CIDR allowed to reach SSH in the driver's default security group.",
		},
		"ssh_port": {
			Type: schema.TypeString, Optional: true, Default: "22",
			Description: "SSH port, encoded as a string.",
		},
		"ssh_user": {
			Type: schema.TypeString, Optional: true, Default: "ubuntu",
			Description: "Machine SSH username.",
		},
		"subnet_id": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public subnet ID.",
		},
		"subnet_name": {
			Type: schema.TypeString, Optional: true, Default: "subnet-docker-machine",
			Description: "T-Cloud Public subnet name.",
		},
		"tags": {
			Type: schema.TypeString, Optional: true,
			Description: "Comma-separated instance tags in key.value form.",
		},
		"token": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description: "T-Cloud Public authentication token.",
		},
		"user_data_file": {
			Type: schema.TypeString, Optional: true,
			Description: "User-data file path on the machine provisioner, not the Terraform host.",
		},
		"user_data_raw": {
			Type: schema.TypeString, Optional: true, Sensitive: true,
			Description: "Inline user-data content passed to the driver.",
		},
		"username": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public username. Prefer a cloud credential.",
		},
		"vpc_id": {
			Type: schema.TypeString, Optional: true,
			Description: "T-Cloud Public VPC ID.",
		},
		"vpc_name": {
			Type: schema.TypeString, Optional: true, Default: "vpc-docker-machine",
			Description: "T-Cloud Public VPC name.",
		},
	}
}
