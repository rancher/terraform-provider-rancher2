package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2IonoscloudKind         = "IonoscloudConfig"
	machineConfigV2IonoscloudAPIVersion   = "rke-machine-config.cattle.io/v1"
	machineConfigV2IonoscloudAPIType      = "rke-machine-config.cattle.io.ionoscloudconfig"
	machineConfigV2IonoscloudClusterIDsep = "."
)

// Types

type machineConfigV2Ionoscloud struct {
	metav1.TypeMeta        `json:",inline"`
	metav1.ObjectMeta      `json:"metadata,omitempty"`
	Endpoint               string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Username               string   `json:"username,omitempty" yaml:"username,omitempty"`
	Password               string   `json:"password,omitempty" yaml:"password,omitempty"`
	Token                  string   `json:"token,omitempty" yaml:"token,omitempty"`
	ServerCores            int      `json:"serverCores,string,omitempty" yaml:"serverCores,omitempty"`
	ServerRam              int      `json:"serverRam,string,omitempty" yaml:"serverRam,omitempty"`
	ServerCpuFamily        string   `json:"serverCpuFamily,omitempty" yaml:"serverCpuFamily,omitempty"`
	ServerAvailabilityZone string   `json:"serverAvailabilityZone,omitempty" yaml:"serverAvailabilityZone,omitempty"`
	DiskSize               int      `json:"diskSize,string,omitempty" yaml:"diskSize,omitempty"`
	DiskType               string   `json:"diskType,omitempty" yaml:"diskType,omitempty"`
	ServerType             string   `json:"serverType,omitempty" yaml:"serverType,omitempty"`
	Template               string   `json:"template,omitempty" yaml:"template,omitempty"`
	Image                  string   `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePassword          string   `json:"imagePassword,omitempty" yaml:"imagePassword,omitempty"`
	Location               string   `json:"location,omitempty" yaml:"location,omitempty"`
	DatacenterId           string   `json:"datacenterId,omitempty" yaml:"datacenterId,omitempty"`
	DatacenterName         string   `json:"datacenterName,omitempty" yaml:"datacenterName,omitempty"`
	LanId                  string   `json:"lanId,omitempty" yaml:"lanId,omitempty"`
	NicDhcp                bool     `json:"nicDhcp,omitempty" yaml:"nicDhcp,omitempty"`
	NicIps                 []string `json:"nicIps,omitempty" yaml:"nicIps,omitempty"`
	LanName                string   `json:"lanName,omitempty" yaml:"lanName,omitempty"`
	VolumeAvailabilityZone string   `json:"volumeAvailabilityZone,omitempty" yaml:"volumeAvailabilityZone,omitempty"`
	CloudInit              string   `json:"cloudInit,omitempty" yaml:"cloudInit,omitempty"`
	CloudInitB64           string   `json:"cloudInitB64,omitempty" yaml:"cloudInitB64,omitempty"`
	SshInCloudInit         bool     `json:"sshInCloudInit,omitempty" yaml:"sshInCloudInit,omitempty"`
	SshUser                string   `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	WaitForIpChange        bool     `json:"waitForIpChange,omitempty" yaml:"waitForIpChange,omitempty"`
	WaitForIpChangeTimeout int      `json:"waitForIpChangeTimeout,string,omitempty" yaml:"waitForIpChangeTimeout,omitempty"`
	NatId                  string   `json:"natId,omitempty" yaml:"natId,omitempty"`
	NatName                string   `json:"natName,omitempty" yaml:"natName,omitempty"`
	NatPublicIps           []string `json:"natPublicIps,omitempty" yaml:"natPublicIps,omitempty"`
	NatFlowlogs            []string `json:"natFlowlogs,omitempty" yaml:"natFlowlogs,omitempty"`
	NatRules               []string `json:"natRules,omitempty" yaml:"natRules,omitempty"`
	SkipDefaultNatRules    bool     `json:"skipDefaultNatRules,omitempty" yaml:"skipDefaultNatRules,omitempty"`
	NatLansToGateways      string   `json:"natLansToGateways,omitempty" yaml:"natLansToGateways,omitempty"`
	PrivateLan             bool     `json:"privateLan,omitempty" yaml:"privateLan,omitempty"`
	AdditionalLans         []string `json:"additionalLans,omitempty" yaml:"additionalLans,omitempty"`
	CreateNat              bool     `json:"createNat,omitempty" yaml:"createNat,omitempty"`
	AppendRkeUserdata      bool     `json:"appendRkeUserdata,omitempty" yaml:"appendRkeUserdata,omitempty"`
}

type MachineConfigV2Ionoscloud struct {
	norman.Resource
	machineConfigV2Ionoscloud
}

// Flatteners

func flattenMachineConfigV2Ionoscloud(in *MachineConfigV2Ionoscloud) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.Endpoint) > 0 {
		obj["endpoint"] = in.Endpoint
	}

	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}

	obj["server_cores"] = in.ServerCores
	obj["server_ram"] = in.ServerRam

	if len(in.ServerCpuFamily) > 0 {
		obj["server_cpu_family"] = in.ServerCpuFamily
	}

	if len(in.ServerAvailabilityZone) > 0 {
		obj["server_availability_zone"] = in.ServerAvailabilityZone
	}

	obj["disk_size"] = in.DiskSize

	if len(in.DiskType) > 0 {
		obj["disk_type"] = in.DiskType
	}

	if len(in.ServerType) > 0 {
		obj["server_type"] = in.ServerType
	}

	if len(in.Template) > 0 {
		obj["template"] = in.Template
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.ImagePassword) > 0 {
		obj["image_password"] = in.ImagePassword
	}

	if len(in.Location) > 0 {
		obj["location"] = in.Location
	}

	if len(in.DatacenterId) > 0 {
		obj["datacenter_id"] = in.DatacenterId
	}

	if len(in.DatacenterName) > 0 {
		obj["datacenter_name"] = in.DatacenterName
	}

	if len(in.LanId) > 0 {
		obj["lan_id"] = in.LanId
	}

	obj["nic_dhcp"] = in.NicDhcp

	if len(in.NicIps) > 0 {
		obj["nic_ips"] = toArrayInterface(in.NicIps)
	}

	if len(in.LanName) > 0 {
		obj["lan_name"] = in.LanName
	}

	if len(in.VolumeAvailabilityZone) > 0 {
		obj["volume_availability_zone"] = in.VolumeAvailabilityZone
	}

	if len(in.CloudInit) > 0 {
		obj["cloud_init"] = in.CloudInit
	}

	if len(in.CloudInitB64) > 0 {
		obj["cloud_init_b64"] = in.CloudInitB64
	}

	obj["ssh_in_cloud_init"] = in.SshInCloudInit

	if len(in.SshUser) > 0 {
		obj["ssh_user"] = in.SshUser
	}

	obj["wait_for_ip_change"] = in.WaitForIpChange
	obj["wait_for_ip_change_timeout"] = in.WaitForIpChangeTimeout

	if len(in.NatId) > 0 {
		obj["nat_id"] = in.NatId
	}

	if len(in.NatName) > 0 {
		obj["nat_name"] = in.NatName
	}

	if len(in.NatPublicIps) > 0 {
		obj["nat_public_ips"] = toArrayInterface(in.NatPublicIps)
	}

	if len(in.NatFlowlogs) > 0 {
		obj["nat_flowlogs"] = toArrayInterface(in.NatFlowlogs)
	}

	if len(in.NatRules) > 0 {
		obj["nat_rules"] = toArrayInterface(in.NatRules)
	}

	obj["skip_default_nat_rules"] = in.SkipDefaultNatRules

	if len(in.NatLansToGateways) > 0 {
		obj["nat_lans_to_gateways"] = in.NatLansToGateways
	}

	obj["private_lan"] = in.PrivateLan

	if len(in.AdditionalLans) > 0 {
		obj["additional_lans"] = toArrayInterface(in.AdditionalLans)
	}

	obj["create_nat"] = in.CreateNat
	obj["append_rke_userdata"] = in.AppendRkeUserdata

	return []interface{}{obj}
}

// Expanders

func expandMachineConfigV2Ionoscloud(p []interface{}, source *MachineConfigV2) *MachineConfigV2Ionoscloud {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &MachineConfigV2Ionoscloud{}

	if len(source.ID) > 0 {
		obj.ID = source.ID
	}
	in := p[0].(map[string]interface{})

	obj.TypeMeta.Kind = machineConfigV2IonoscloudKind
	obj.TypeMeta.APIVersion = machineConfigV2IonoscloudAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	var (
		isServerTypeCube     bool
		isExistingDatacenter bool
		isExistingLan        bool
		isExistingNat        bool
	)
	if v, ok := in["server_type"].(string); ok && len(v) > 0 {
		obj.ServerType = v
		isServerTypeCube = v == "CUBE"
	}
	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	if v, ok := in["server_cores"].(int); ok {
		obj.ServerCores = v
	} else if isServerTypeCube {
		obj.ServerCores = 0
	}

	if v, ok := in["server_ram"].(int); ok && !isServerTypeCube {
		obj.ServerRam = v
	} else if isExistingDatacenter {
		obj.ServerRam = 0
	}

	if v, ok := in["server_cpu_family"].(string); ok && len(v) > 0 && !isServerTypeCube {
		obj.ServerCpuFamily = v
	}

	if v, ok := in["server_availability_zone"].(string); ok && len(v) > 0 {
		obj.ServerAvailabilityZone = v
	}

	if v, ok := in["disk_size"].(int); ok {
		obj.DiskSize = v
	} else if isServerTypeCube {
		obj.DiskSize = 0
	}

	if v, ok := in["disk_type"].(string); ok && len(v) > 0 {
		obj.DiskType = v
	}

	if v, ok := in["template"].(string); ok && len(v) > 0 {
		obj.Template = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["image_password"].(string); ok && len(v) > 0 {
		obj.ImagePassword = v
	}

	if v, ok := in["location"].(string); ok && len(v) > 0 {
		obj.Location = v
	}

	if v, ok := in["datacenter_id"].(string); ok && len(v) > 0 {
		obj.DatacenterId = v
	}

	if v, ok := in["datacenter_name"].(string); ok && len(v) > 0 {
		obj.DatacenterName = v
	} else if isExistingDatacenter {
		obj.DatacenterName = ""
	}

	if v, ok := in["lan_id"].(string); ok && len(v) > 0 {
		obj.LanId = v
	}

	if v, ok := in["nic_dhcp"].(bool); ok {
		obj.NicDhcp = v
	}

	if v, ok := in["nic_ips"].([]interface{}); ok && len(v) > 0 {
		obj.NicIps = toArrayString(v)
	}

	if v, ok := in["lan_name"].(string); ok && len(v) > 0 {
		obj.LanName = v
	} else if isExistingLan {
		obj.DatacenterName = ""
	}

	if v, ok := in["volume_availability_zone"].(string); ok && len(v) > 0 {
		obj.VolumeAvailabilityZone = v
	}

	if v, ok := in["cloud_init"].(string); ok && len(v) > 0 {
		obj.CloudInit = v
	}

	if v, ok := in["cloud_init_b64"].(string); ok && len(v) > 0 {
		obj.CloudInitB64 = v
	}

	if v, ok := in["ssh_in_cloud_init"].(bool); ok {
		obj.SshInCloudInit = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SshUser = v
	}

	if v, ok := in["wait_for_ip_change"].(bool); ok {
		obj.WaitForIpChange = v
	}

	if v, ok := in["wait_for_ip_change_timeout"].(int); ok {
		obj.WaitForIpChangeTimeout = v
	}

	if v, ok := in["nat_id"].(string); ok && len(v) > 0 {
		obj.NatId = v
	}

	if v, ok := in["nat_name"].(string); ok && len(v) > 0 {
		obj.NatName = v
	} else if isExistingNat {
		obj.NatName = ""
	}

	if v, ok := in["nat_public_ips"].([]interface{}); ok && len(v) > 0 {
		obj.NatPublicIps = toArrayString(v)
	}

	if v, ok := in["nat_flowlogs"].([]interface{}); ok && len(v) > 0 {
		obj.NatFlowlogs = toArrayString(v)
	}

	if v, ok := in["nat_rules"].([]interface{}); ok && len(v) > 0 {
		obj.NatRules = toArrayString(v)
	}

	if v, ok := in["skip_default_nat_rules"].(bool); ok {
		obj.SkipDefaultNatRules = v
	}

	if v, ok := in["nat_lans_to_gateways"].(string); ok && len(v) > 0 {
		obj.NatLansToGateways = v
	}

	if v, ok := in["private_lan"].(bool); ok {
		obj.PrivateLan = v
	}

	if v, ok := in["additional_lans"].([]interface{}); ok && len(v) > 0 {
		obj.AdditionalLans = toArrayString(v)
	}

	if v, ok := in["create_nat"].(bool); ok {
		obj.CreateNat = v
	}

	if v, ok := in["append_rke_userdata"].(bool); ok {
		obj.AppendRkeUserdata = v
	}

	return obj
}
