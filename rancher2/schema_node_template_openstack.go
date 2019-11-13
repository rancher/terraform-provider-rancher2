package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	openstackConfigDriver = "openstack"
)

//Types

type openstackConfig struct {
	ActiveTimeout    string `json:"activeTimeout,omitempty" yaml:"activeTimeout,omitempty"`
	AuthURL          string `json:"authUrl,omitempty" yaml:"authUrl,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty" yaml:"availabilityZone,omitempty"`
	CaCert           string `json:"cacert,omitempty" yaml:"cacert,omitempty"`
	ConfigDrive      bool   `json:"configDrive,omitempty" yaml:"configDrive,omitempty"`
	DomainID         string `json:"domainId,omitempty" yaml:"domainId,omitempty"`
	DomainName       string `json:"domainName,omitempty" yaml:"domainName,omitempty"`
	EndpointType     string `json:"endpointType,omitempty" yaml:"endpointType,omitempty"`
	FlavorID         string `json:"flavorId,omitempty" yaml:"flavorId,omitempty"`
	FlavorName       string `json:"flavorName,omitempty" yaml:"flavorName,omitempty"`
	FloatingIPPool   string `json:"floatingipPool,omitempty" yaml:"floatingipPool,omitempty"`
	ImageID          string `json:"imageId,omitempty" yaml:"imageId,omitempty"`
	ImageName        string `json:"imageName,omitempty" yaml:"imageName,omitempty"`
	Insecure         bool   `json:"insecure,omitempty" yaml:"insecure,omitempty"`
	IPVersion        string `json:"ipVersion,omitempty" yaml:"ipVersion,omitempty"`
	KeypairName      string `json:"keypairName,omitempty" yaml:"keypairName,omitempty"`
	NetID            string `json:"netId,omitempty" yaml:"netId,omitempty"`
	NetName          string `json:"netName,omitempty" yaml:"netName,omitempty"`
	NovaNetwork      bool   `json:"novaNetwork,omitempty" yaml:"novaNetwork,omitempty"`
	Password         string `json:"password,omitempty" yaml:"password,omitempty"`
	PrivateKeyFile   string `json:"privateKeyFile,omitempty" yaml:"privateKeyFile,omitempty"`
	Region           string `json:"region,omitempty" yaml:"region,omitempty"`
	SecGroups        string `json:"secGroups,omitempty" yaml:"secGroups,omitempty"`
	SSHPort          string `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser          string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	TenantID         string `json:"tenantId,omitempty" yaml:"tenantId,omitempty"`
	TenantName       string `json:"tenantName,omitempty" yaml:"tenantName,omitempty"`
	UserDataFile     string `json:"userDataFile,omitempty" yaml:"userDataFile,omitempty"`
	Username         string `json:"username,omitempty" yaml:"username,omitempty"`
}

//Schemas

func openstackConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"auth_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"availability_zone": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"region": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"active_timeout": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "200",
		},
		"cacert": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"config_drive": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"domain_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"domain_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"endpoint_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"flavor_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"flavor_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"floating_ip_pool": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"image_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"insecure": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"ip_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "4",
		},
		"keypair_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"net_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"net_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"nova_network": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"password": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"private_key_file": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"sec_groups": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"ssh_port": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "22",
		},
		"ssh_user": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "root",
		},
		"tenant_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"tenant_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_data_file": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}
