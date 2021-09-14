package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2OpenstackKind         = "OpenstackConfig"
	machineConfigV2OpenstackAPIVersion   = "rke-machine-config.cattle.io/v1"
	machineConfigV2OpenstackAPIType      = "rke-machine-config.cattle.io.openstackconfig"
	machineConfigV2OpenstackClusterIDsep = "."
)

//Types

type machineConfigV2Openstack struct {
	metav1.TypeMeta             `json:",inline"`
	metav1.ObjectMeta           `json:"metadata,omitempty"`
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
	TenantDomainID              string `json:"tenantDomainId,omitempty" yaml:"tenantDomainId,omitempty"`
	TenantDomainName            string `json:"tenantDomainName,omitempty" yaml:"tenantDomainName,omitempty"`
	UserDataFile                string `json:"userDataFile,omitempty" yaml:"userDataFile,omitempty"`
	Username                    string `json:"username,omitempty" yaml:"username,omitempty"`
	UserDomainID                string `json:"userDomainId,omitempty" yaml:"userDomainId,omitempty"`
	UserDomainName              string `json:"userDomainName,omitempty" yaml:"userDomainName,omitempty"`
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

type MachineConfigV2Openstack struct {
	norman.Resource
	machineConfigV2Openstack
}

// Flatteners

func flattenMachineConfigV2Openstack(in *MachineConfigV2Openstack) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	obj["active_timeout"] = in.ActiveTimeout
	obj["auth_url"] = in.AuthURL
	obj["availability_zone"] = in.AvailabilityZone
	obj["cacert"] = in.CaCert
	obj["config_drive"] = in.ConfigDrive
	obj["domain_id"] = in.DomainID
	obj["domain_name"] = in.DomainName
	obj["endpoint_type"] = in.EndpointType
	obj["flavor_id"] = in.FlavorID
	obj["flavor_name"] = in.FlavorName
	obj["floating_ip_pool"] = in.FloatingIPPool
	obj["image_id"] = in.ImageID
	obj["image_name"] = in.ImageName
	obj["insecure"] = in.Insecure
	obj["ip_version"] = in.IPVersion
	obj["keypair_name"] = in.KeypairName
	obj["net_id"] = in.NetID
	obj["net_name"] = in.NetName
	obj["nova_network"] = in.NovaNetwork
	obj["password"] = in.Password
	obj["private_key_file"] = in.PrivateKeyFile
	obj["region"] = in.Region
	obj["sec_groups"] = in.SecGroups
	obj["ssh_port"] = in.SSHPort
	obj["ssh_user"] = in.SSHUser
	obj["tenant_id"] = in.TenantID
	obj["tenant_name"] = in.TenantName
	obj["tenant_domain_id"] = in.TenantDomainID
	obj["tenant_domain_name"] = in.TenantDomainName
	obj["user_data_file"] = in.UserDataFile
	obj["username"] = in.Username
	obj["user_domain_id"] = in.UserDomainID
	obj["user_domain_name"] = in.UserDomainName
	obj["application_credential_id"] = in.ApplicationCredentialID
	obj["application_credential_name"] = in.ApplicationCredentialName
	obj["application_credential_secret"] = in.ApplicationCredentialSecret
	obj["boot_from_volume"] = in.BootFromVolume
	obj["volume_size"] = in.VolumeSize
	obj["volume_type"] = in.VolumeType
	obj["volume_id"] = in.VolumeID
	obj["volume_name"] = in.VolumeName
	obj["volume_device_path"] = in.VolumeDevicePath

	return []interface{}{obj}
}

// Expanders

func expandMachineConfigV2Openstack(p []interface{}, source *MachineConfigV2) *MachineConfigV2Openstack {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &MachineConfigV2Openstack{}

	if len(source.ID) > 0 {
		obj.ID = source.ID
	}
	in := p[0].(map[string]interface{})

	obj.TypeMeta.Kind = machineConfigV2OpenstackKind
	obj.TypeMeta.APIVersion = machineConfigV2OpenstackAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	if v, ok := in["active_timeout"].(string); ok && len(v) > 0 {
		obj.ActiveTimeout = v
	}
	if v, ok := in["auth_url"].(string); ok && len(v) > 0 {
		obj.AuthURL = v
	}
	if v, ok := in["availability_zone"].(string); ok && len(v) > 0 {
		obj.AvailabilityZone = v
	}
	if v, ok := in["cacert"].(string); ok && len(v) > 0 {
		obj.CaCert = v
	}
	if v, ok := in["config_drive"].(bool); ok {
		obj.ConfigDrive = v
	}
	if v, ok := in["domain_id"].(string); ok && len(v) > 0 {
		obj.DomainID = v
	}
	if v, ok := in["domain_name"].(string); ok && len(v) > 0 {
		obj.DomainName = v
	}
	if v, ok := in["endpoint_type"].(string); ok && len(v) > 0 {
		obj.EndpointType = v
	}
	if v, ok := in["flavor_id"].(string); ok && len(v) > 0 {
		obj.FlavorID = v
	}
	if v, ok := in["flavor_name"].(string); ok && len(v) > 0 {
		obj.FlavorName = v
	}
	if v, ok := in["floating_ip_pool"].(string); ok && len(v) > 0 {
		obj.FloatingIPPool = v
	}
	if v, ok := in["ip_version"].(string); ok && len(v) > 0 {
		obj.IPVersion = v
	}
	if v, ok := in["image_id"].(string); ok && len(v) > 0 {
		obj.ImageID = v
	}
	if v, ok := in["image_name"].(string); ok && len(v) > 0 {
		obj.ImageName = v
	}
	if v, ok := in["insecure"].(bool); ok {
		obj.Insecure = v
	}
	if v, ok := in["ip_version"].(string); ok && len(v) > 0 {
		obj.IPVersion = v
	}
	if v, ok := in["keypair_name"].(string); ok && len(v) > 0 {
		obj.KeypairName = v
	}
	if v, ok := in["net_id"].(string); ok && len(v) > 0 {
		obj.NetID = v
	}
	if v, ok := in["net_name"].(string); ok && len(v) > 0 {
		obj.NetName = v
	}
	if v, ok := in["nova_network"].(bool); ok {
		obj.NovaNetwork = v
	}
	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in["private_key_file"].(string); ok && len(v) > 0 {
		obj.PrivateKeyFile = v
	}
	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}
	if v, ok := in["sec_groups"].(string); ok && len(v) > 0 {
		obj.SecGroups = v
	}
	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}
	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}
	if v, ok := in["tenant_id"].(string); ok && len(v) > 0 {
		obj.TenantID = v
	}
	if v, ok := in["tenant_name"].(string); ok && len(v) > 0 {
		obj.TenantName = v
	}
	if v, ok := in["tenant_domain_id"].(string); ok && len(v) > 0 {
		obj.TenantDomainID = v
	}
	if v, ok := in["tenant_domain_name"].(string); ok && len(v) > 0 {
		obj.TenantDomainName = v
	}
	if v, ok := in["user_data_file"].(string); ok && len(v) > 0 {
		obj.UserDataFile = v
	}
	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}
	if v, ok := in["user_domain_id"].(string); ok && len(v) > 0 {
		obj.UserDomainID = v
	}
	if v, ok := in["user_domain_name"].(string); ok && len(v) > 0 {
		obj.UserDomainName = v
	}
	if v, ok := in["application_credential_id"].(string); ok && len(v) > 0 {
		obj.ApplicationCredentialID = v
	}
	if v, ok := in["application_credential_name"].(string); ok && len(v) > 0 {
		obj.ApplicationCredentialName = v
	}
	if v, ok := in["application_credential_secret"].(string); ok && len(v) > 0 {
		obj.ApplicationCredentialSecret = v
	}
	if v, ok := in["boot_from_volume"].(bool); ok {
		obj.BootFromVolume = v
	}
	if v, ok := in["volume_size"].(string); ok && len(v) > 0 {
		obj.VolumeSize = v
	}
	if v, ok := in["volume_type"].(string); ok && len(v) > 0 {
		obj.VolumeType = v
	}
	if v, ok := in["volume_id"].(string); ok && len(v) > 0 {
		obj.VolumeID = v
	}
	if v, ok := in["volume_name"].(string); ok && len(v) > 0 {
		obj.VolumeName = v
	}
	if v, ok := in["volume_device_path"].(string); ok && len(v) > 0 {
		obj.VolumeDevicePath = v
	}

	return obj
}
