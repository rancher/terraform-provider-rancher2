package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2OpentelekomcloudKind       = "OpentelekomcloudConfig"
	machineConfigV2OpentelekomcloudAPIVersion = "rke-machine-config.cattle.io/v1"
	machineConfigV2OpentelekomcloudAPIType    = "rke-machine-config.cattle.io.opentelekomcloudconfig"
)

type machineConfigV2Opentelekomcloud struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Empty fields must be sent so updates can clear previously configured values.
	AccessKey        string `json:"accessKey"`
	AuthURL          string `json:"authUrl"`
	AvailabilityZone string `json:"availabilityZone"`
	BandwidthSize    string `json:"bandwidthSize"`
	BandwidthType    string `json:"bandwidthType"`
	CACert           string `json:"cacert"`
	Cloud            string `json:"cloud"`
	DomainID         string `json:"domainId"`
	DomainName       string `json:"domainName"`
	EIP              string `json:"eip"`
	EIPType          string `json:"eipType"`
	EndpointType     string `json:"endpointType"`
	FlavorID         string `json:"flavorId"`
	FlavorName       string `json:"flavorName"`
	ImageID          string `json:"imageId"`
	ImageName        string `json:"imageName"`
	IPVersion        string `json:"ipVersion"`
	KeypairName      string `json:"keypairName"`
	NetworkScope     string `json:"networkScope"`
	Password         string `json:"password"`
	PrivateKeyFile   string `json:"privateKeyFile"`
	ProjectID        string `json:"projectId"`
	ProjectName      string `json:"projectName"`
	Region           string `json:"region"`
	RootVolumeSize   string `json:"rootVolumeSize"`
	RootVolumeType   string `json:"rootVolumeType"`
	SecGroups        string `json:"secGroups"`
	SecretKey        string `json:"secretKey"`
	ServerGroup      string `json:"serverGroup"`
	ServerGroupID    string `json:"serverGroupId"`
	SkipDefaultSG    bool   `json:"skipDefaultSg"`
	SkipEIP          bool   `json:"skipEip"`
	SSHAllowCIDR     string `json:"sshAllowCidr"`
	SSHPort          string `json:"sshPort"`
	SSHUser          string `json:"sshUser"`
	SubnetID         string `json:"subnetId"`
	SubnetName       string `json:"subnetName"`
	Tags             string `json:"tags"`
	Token            string `json:"token"`
	UserDataFile     string `json:"userDataFile"`
	UserDataRaw      string `json:"userDataRaw"`
	Username         string `json:"username"`
	VPCID            string `json:"vpcId"`
	VPCName          string `json:"vpcName"`
}

type MachineConfigV2Opentelekomcloud struct {
	norman.Resource
	machineConfigV2Opentelekomcloud
}

func flattenMachineConfigV2Opentelekomcloud(in *MachineConfigV2Opentelekomcloud) []interface{} {
	if in == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"access_key": in.AccessKey, "auth_url": in.AuthURL,
		"availability_zone": in.AvailabilityZone,
		"bandwidth_size":    in.BandwidthSize, "bandwidth_type": in.BandwidthType,
		"cacert": in.CACert, "cloud": in.Cloud,
		"domain_id": in.DomainID, "domain_name": in.DomainName,
		"eip": in.EIP, "eip_type": in.EIPType, "endpoint_type": in.EndpointType,
		"flavor_id": in.FlavorID, "flavor_name": in.FlavorName,
		"image_id": in.ImageID, "image_name": in.ImageName, "ip_version": in.IPVersion,
		"keypair_name": in.KeypairName, "network_scope": in.NetworkScope,
		"password": in.Password, "private_key_file": in.PrivateKeyFile,
		"project_id": in.ProjectID, "project_name": in.ProjectName, "region": in.Region,
		"root_volume_size": in.RootVolumeSize, "root_volume_type": in.RootVolumeType,
		"sec_groups": in.SecGroups, "secret_key": in.SecretKey,
		"server_group": in.ServerGroup, "server_group_id": in.ServerGroupID,
		"skip_default_sg": in.SkipDefaultSG, "skip_eip": in.SkipEIP,
		"ssh_allow_cidr": in.SSHAllowCIDR, "ssh_port": in.SSHPort, "ssh_user": in.SSHUser,
		"subnet_id": in.SubnetID, "subnet_name": in.SubnetName,
		"tags": in.Tags, "token": in.Token,
		"user_data_file": in.UserDataFile, "user_data_raw": in.UserDataRaw, "username": in.Username,
		"vpc_id": in.VPCID, "vpc_name": in.VPCName,
	}}
}

func expandMachineConfigV2Opentelekomcloud(p []interface{}, source *MachineConfigV2) *MachineConfigV2Opentelekomcloud {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})
	out := &MachineConfigV2Opentelekomcloud{}
	out.Resource = source.Resource
	out.TypeMeta = metav1.TypeMeta{
		Kind: machineConfigV2OpentelekomcloudKind, APIVersion: machineConfigV2OpentelekomcloudAPIVersion,
	}
	source.TypeMeta = out.TypeMeta
	out.ObjectMeta = source.ObjectMeta
	out.AccessKey, _ = in["access_key"].(string)
	out.AuthURL, _ = in["auth_url"].(string)
	out.AvailabilityZone, _ = in["availability_zone"].(string)
	out.BandwidthSize, _ = in["bandwidth_size"].(string)
	out.BandwidthType, _ = in["bandwidth_type"].(string)
	out.CACert, _ = in["cacert"].(string)
	out.Cloud, _ = in["cloud"].(string)
	out.DomainID, _ = in["domain_id"].(string)
	out.DomainName, _ = in["domain_name"].(string)
	out.EIP, _ = in["eip"].(string)
	out.EIPType, _ = in["eip_type"].(string)
	out.EndpointType, _ = in["endpoint_type"].(string)
	out.FlavorID, _ = in["flavor_id"].(string)
	out.FlavorName, _ = in["flavor_name"].(string)
	out.ImageID, _ = in["image_id"].(string)
	out.ImageName, _ = in["image_name"].(string)
	out.IPVersion, _ = in["ip_version"].(string)
	out.KeypairName, _ = in["keypair_name"].(string)
	out.NetworkScope, _ = in["network_scope"].(string)
	out.Password, _ = in["password"].(string)
	out.PrivateKeyFile, _ = in["private_key_file"].(string)
	out.ProjectID, _ = in["project_id"].(string)
	out.ProjectName, _ = in["project_name"].(string)
	out.Region, _ = in["region"].(string)
	out.RootVolumeSize, _ = in["root_volume_size"].(string)
	out.RootVolumeType, _ = in["root_volume_type"].(string)
	out.SecGroups, _ = in["sec_groups"].(string)
	out.SecretKey, _ = in["secret_key"].(string)
	out.ServerGroup, _ = in["server_group"].(string)
	out.ServerGroupID, _ = in["server_group_id"].(string)
	out.SkipDefaultSG, _ = in["skip_default_sg"].(bool)
	out.SkipEIP, _ = in["skip_eip"].(bool)
	out.SSHAllowCIDR, _ = in["ssh_allow_cidr"].(string)
	out.SSHPort, _ = in["ssh_port"].(string)
	out.SSHUser, _ = in["ssh_user"].(string)
	out.SubnetID, _ = in["subnet_id"].(string)
	out.SubnetName, _ = in["subnet_name"].(string)
	out.Tags, _ = in["tags"].(string)
	out.Token, _ = in["token"].(string)
	out.UserDataFile, _ = in["user_data_file"].(string)
	out.UserDataRaw, _ = in["user_data_raw"].(string)
	out.Username, _ = in["username"].(string)
	out.VPCID, _ = in["vpc_id"].(string)
	out.VPCName, _ = in["vpc_name"].(string)
	return out
}
