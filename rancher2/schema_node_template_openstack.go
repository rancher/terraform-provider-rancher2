package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	openstackConfigDriver = "openstack"
)

//Types

type openstackConfig struct {
	ActiveTimeout               string `json:"activeTimeout,omitempty" yaml:"activeTimeout,omitempty"`
	AuthURL                     string `json:"authUrl,omitempty" yaml:"authUrl,omitempty"`
	AvailabilityZone            string `json:"availabilityZone,omitempty" yaml:"availabilityZone,omitempty"`
	CaCert                      string `json:"cacert,omitempty" yaml:"cacert,omitempty"`
	ConfigDrive                 bool   `json:"configDrive,omitempty" yaml:"configDrive,omitempty"`
	DomainID                    string `json:"domainId,omitempty" yaml:"domainId,omitempty"`
	DomainName                  string `json:"domainName,omitempty" yaml:"domainName,omitempty"`
	EndpointType                string `json:"endpointType,omitempty" yaml:"endpointType,omitempty"`
	FlavorID                    string `json:"flavorId,omitempty" yaml:"flavorId,omitempty"`
	FlavorName                  string `json:"flavorName,omitempty" yaml:"flavorName,omitempty"`
	FloatingIPPool              string `json:"floatingipPool,omitempty" yaml:"floatingipPool,omitempty"`
	ImageID                     string `json:"imageId,omitempty" yaml:"imageId,omitempty"`
	ImageName                   string `json:"imageName,omitempty" yaml:"imageName,omitempty"`
	Insecure                    bool   `json:"insecure,omitempty" yaml:"insecure,omitempty"`
	IPVersion                   string `json:"ipVersion,omitempty" yaml:"ipVersion,omitempty"`
	KeypairName                 string `json:"keypairName,omitempty" yaml:"keypairName,omitempty"`
	NetID                       string `json:"netId,omitempty" yaml:"netId,omitempty"`
	NetName                     string `json:"netName,omitempty" yaml:"netName,omitempty"`
	NovaNetwork                 bool   `json:"novaNetwork,omitempty" yaml:"novaNetwork,omitempty"`
	Password                    string `json:"password,omitempty" yaml:"password,omitempty"`
	PrivateKeyFile              string `json:"privateKeyFile,omitempty" yaml:"privateKeyFile,omitempty"`
	Region                      string `json:"region,omitempty" yaml:"region,omitempty"`
	SecGroups                   string `json:"secGroups,omitempty" yaml:"secGroups,omitempty"`
	SSHPort                     string `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser                     string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	TenantID                    string `json:"tenantId,omitempty" yaml:"tenantId,omitempty"`
	TenantName                  string `json:"tenantName,omitempty" yaml:"tenantName,omitempty"`
	UserDataFile                string `json:"userDataFile,omitempty" yaml:"userDataFile,omitempty"`
	Username                    string `json:"username,omitempty" yaml:"username,omitempty"`
	ApplicationCredentialID     string `json:"applicationCredentialId,omitempty" yaml:"applicationCredentialId,omitempty"`
	ApplicationCredentialName   string `json:"applicationCredentialName,omitempty" yaml:"applicationCredentialName,omitempty"`
	ApplicationCredentialSecret string `json:"applicationCredentialSecret,omitempty" yaml:"applicationCredentialSecret,omitempty"`
	BootFromVolume              bool   `json:"bootFromVolume,omitempty" yaml:"bootFromVolume,omitempty"`
	VolumeType                  string `json:"volumeType,omitempty" yaml:"volumeType,omitempty"`
	VolumeSize                  string `json:"volumeSize,omitempty" yaml:"volumeSize,omitempty"`
	VolumeID                    string `json:"volumeId,omitempty" yaml:"volumeId,omitempty"`
	VolumeName                  string `json:"volumeName,omitempty" yaml:"volumeName,omitempty"`
	VolumeDevicePath            string `json:"volumeDevicePath,omitempty" yaml:"volumeDevicePath,omitempty"`
}

//Schemas

func openstackConfigFields() map[string]*schema.Schema {
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
