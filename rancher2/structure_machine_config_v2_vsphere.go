package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2VmwarevsphereKind         = "VmwarevsphereConfig"
	machineConfigV2VmwarevsphereAPIVersion   = "rke-machine-config.cattle.io/v1"
	machineConfigV2VmwarevsphereAPIType      = "rke-machine-config.cattle.io.vmwarevsphereconfig"
	machineConfigV2VmwarevsphereClusterIDsep = "."
)

//Types

type machineConfigV2Vmwarevsphere struct {
	metav1.TypeMeta        `json:",inline"`
	metav1.ObjectMeta      `json:"metadata,omitempty"`
	Boot2dockerURL         string   `json:"boot2dockerUrl,omitempty" yaml:"boot2dockerUrl,omitempty"`
	Cfgparam               []string `json:"cfgparam,omitempty" yaml:"cfgparam,omitempty"`
	CloneFrom              string   `json:"cloneFrom,omitempty" yaml:"cloneFrom,omitempty"`
	CloudConfig            string   `json:"cloudConfig,omitempty" yaml:"cloudConfig,omitempty"`
	Cloudinit              string   `json:"cloudinit,omitempty" yaml:"cloudinit,omitempty"`
	ContentLibrary         string   `json:"contentLibrary,omitempty" yaml:"contentLibrary,omitempty"`
	CPUCount               string   `json:"cpuCount,omitempty" yaml:"cpuCount,omitempty"`
	CreationType           string   `json:"creationType,omitempty" yaml:"creationType,omitempty"`
	CustomAttributes       []string `json:"customAttribute,omitempty" yaml:"customAttribute,omitempty"`
	Datacenter             string   `json:"datacenter,omitempty" yaml:"datacenter,omitempty"`
	Datastore              string   `json:"datastore,omitempty" yaml:"datastore,omitempty"`
	DatastoreCluster       string   `json:"datastoreCluster,omitempty" yaml:"datastoreCluster,omitempty"`
	DiskSize               string   `json:"diskSize,omitempty" yaml:"diskSize,omitempty"`
	Folder                 string   `json:"folder,omitempty" yaml:"folder,omitempty"`
	Hostsystem             string   `json:"hostsystem,omitempty" yaml:"hostsystem,omitempty"`
	MemorySize             string   `json:"memorySize,omitempty" yaml:"memorySize,omitempty"`
	Network                []string `json:"network,omitempty" yaml:"network,omitempty"`
	Password               string   `json:"password,omitempty" yaml:"password,omitempty"`
	Pool                   string   `json:"pool,omitempty" yaml:"pool,omitempty"`
	SSHPassword            string   `json:"sshPassword,omitempty" yaml:"sshPassword,omitempty"`
	SSHPort                string   `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser                string   `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	SSHUserGroup           string   `json:"sshUserGroup,omitempty" yaml:"sshUserGroup,omitempty"`
	Tags                   []string `json:"tag,omitempty" yaml:"tag,omitempty"`
	Username               string   `json:"username,omitempty" yaml:"username,omitempty"`
	VappIpallocationpolicy string   `json:"vappIpallocationpolicy,omitempty" yaml:"vappIpallocationpolicy,omitempty"`
	VappIpprotocol         string   `json:"vappIpprotocol,omitempty" yaml:"vappIpprotocol,omitempty"`
	VappProperty           []string `json:"vappProperty,omitempty" yaml:"vappProperty,omitempty"`
	VappTransport          string   `json:"vappTransport,omitempty" yaml:"vappTransport,omitempty"`
	Vcenter                string   `json:"vcenter,omitempty" yaml:"vcenter,omitempty"`
	VcenterPort            string   `json:"vcenterPort,omitempty" yaml:"vcenterPort,omitempty"`
}

type MachineConfigV2Vmwarevsphere struct {
	norman.Resource
	machineConfigV2Vmwarevsphere
}

// Flatteners

func flattenMachineConfigV2Vmwarevsphere(in *MachineConfigV2Vmwarevsphere) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.Boot2dockerURL) > 0 {
		obj["boot2docker_url"] = in.Boot2dockerURL
	}
	if len(in.Cfgparam) > 0 {
		obj["cfgparam"] = toArrayInterface(in.Cfgparam)
	}
	if len(in.CloneFrom) > 0 {
		obj["clone_from"] = in.CloneFrom
	}
	if len(in.CloudConfig) > 0 {
		obj["cloud_config"] = in.CloudConfig
	}
	if len(in.Cloudinit) > 0 {
		obj["cloudinit"] = in.Cloudinit
	}
	if len(in.ContentLibrary) > 0 {
		obj["content_library"] = in.ContentLibrary
	}
	if len(in.CPUCount) > 0 {
		obj["cpu_count"] = in.CPUCount
	}
	if len(in.CreationType) > 0 {
		obj["creation_type"] = in.CreationType
	}
	if len(in.CustomAttributes) > 0 {
		obj["custom_attributes"] = toArrayInterface(in.CustomAttributes)
	}
	if len(in.Datacenter) > 0 {
		obj["datacenter"] = in.Datacenter
	}
	if len(in.Datastore) > 0 {
		obj["datastore"] = in.Datastore
	}
	if len(in.DatastoreCluster) > 0 {
		obj["datastore_cluster"] = in.DatastoreCluster
	}
	if len(in.DiskSize) > 0 {
		obj["disk_size"] = in.DiskSize
	}
	if len(in.Folder) > 0 {
		obj["folder"] = in.Folder
	}
	if len(in.Hostsystem) > 0 {
		obj["hostsystem"] = in.Hostsystem
	}
	if len(in.MemorySize) > 0 {
		obj["memory_size"] = in.MemorySize
	}
	if len(in.Network) > 0 {
		obj["network"] = toArrayInterface(in.Network)
	}
	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}
	if len(in.Pool) > 0 {
		obj["pool"] = in.Pool
	}
	if len(in.SSHPassword) > 0 {
		obj["ssh_password"] = in.SSHPassword
	}
	if len(in.SSHPort) > 0 {
		obj["ssh_port"] = in.SSHPort
	}
	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}
	if len(in.SSHUserGroup) > 0 {
		obj["ssh_user_group"] = in.SSHUserGroup
	}
	if len(in.Tags) > 0 {
		obj["tags"] = toArrayInterface(in.Tags)
	}
	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}
	if len(in.VappIpallocationpolicy) > 0 {
		obj["vapp_ip_allocation_policy"] = in.VappIpallocationpolicy
	}
	if len(in.VappIpprotocol) > 0 {
		obj["vapp_ip_protocol"] = in.VappIpprotocol
	}
	if len(in.VappProperty) > 0 {
		obj["vapp_property"] = toArrayInterface(in.VappProperty)
	}
	if len(in.VappTransport) > 0 {
		obj["vapp_transport"] = in.VappTransport
	}
	if len(in.Vcenter) > 0 {
		obj["vcenter"] = in.Vcenter
	}
	if len(in.VcenterPort) > 0 {
		obj["vcenter_port"] = in.VcenterPort
	}

	return []interface{}{obj}
}

// Expanders

func expandMachineConfigV2Vmwarevsphere(p []interface{}, source *MachineConfigV2) *MachineConfigV2Vmwarevsphere {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &MachineConfigV2Vmwarevsphere{}

	if len(source.ID) > 0 {
		obj.ID = source.ID
	}
	in := p[0].(map[string]interface{})

	obj.TypeMeta.Kind = machineConfigV2VmwarevsphereKind
	obj.TypeMeta.APIVersion = machineConfigV2VmwarevsphereAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	if v, ok := in["boot2docker_url"].(string); ok && len(v) > 0 {
		obj.Boot2dockerURL = v
	}
	if v, ok := in["cfgparam"].([]interface{}); ok && len(v) > 0 {
		obj.Cfgparam = toArrayString(v)
	}
	if v, ok := in["clone_from"].(string); ok && len(v) > 0 {
		obj.CloneFrom = v
	}
	if v, ok := in["cloud_config"].(string); ok && len(v) > 0 {
		obj.CloudConfig = v
	}
	if v, ok := in["cloudinit"].(string); ok && len(v) > 0 {
		obj.Cloudinit = v
	}
	if v, ok := in["content_library"].(string); ok && len(v) > 0 {
		obj.ContentLibrary = v
	}
	if v, ok := in["cpu_count"].(string); ok && len(v) > 0 {
		obj.CPUCount = v
	}
	if v, ok := in["creation_type"].(string); ok && len(v) > 0 {
		obj.CreationType = v
	}
	if v, ok := in["custom_attributes"].([]interface{}); ok && len(v) > 0 {
		obj.CustomAttributes = toArrayString(v)
	}
	if v, ok := in["datacenter"].(string); ok && len(v) > 0 {
		obj.Datacenter = v
	}
	if v, ok := in["datastore"].(string); ok && len(v) > 0 {
		obj.Datastore = v
	}
	if v, ok := in["datastore_cluster"].(string); ok && len(v) > 0 {
		obj.DatastoreCluster = v
	}
	if v, ok := in["disk_size"].(string); ok && len(v) > 0 {
		obj.DiskSize = v
	}
	if v, ok := in["folder"].(string); ok && len(v) > 0 {
		obj.Folder = v
	}
	if v, ok := in["hostsystem"].(string); ok && len(v) > 0 {
		obj.Hostsystem = v
	}
	if v, ok := in["memory_size"].(string); ok && len(v) > 0 {
		obj.MemorySize = v
	}
	if v, ok := in["network"].([]interface{}); ok && len(v) > 0 {
		obj.Network = toArrayString(v)
	}
	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in["pool"].(string); ok && len(v) > 0 {
		obj.Pool = v
	}
	if v, ok := in["ssh_password"].(string); ok && len(v) > 0 {
		obj.SSHPassword = v
	}
	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}
	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}
	if v, ok := in["ssh_user_group"].(string); ok && len(v) > 0 {
		obj.SSHUserGroup = v
	}
	if v, ok := in["tags"].([]interface{}); ok && len(v) > 0 {
		obj.Tags = toArrayString(v)
	}
	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}
	if v, ok := in["vapp_ip_allocation_policy"].(string); ok && len(v) > 0 {
		obj.VappIpallocationpolicy = v
	}
	if v, ok := in["vapp_ip_protocol"].(string); ok && len(v) > 0 {
		obj.VappIpprotocol = v
	}
	if v, ok := in["vapp_property"].([]interface{}); ok && len(v) > 0 {
		obj.VappProperty = toArrayString(v)
	}
	if v, ok := in["vapp_transport"].(string); ok && len(v) > 0 {
		obj.VappTransport = v
	}
	if v, ok := in["vcenter"].(string); ok && len(v) > 0 {
		obj.Vcenter = v
	}
	if v, ok := in["vcenter_port"].(string); ok && len(v) > 0 {
		obj.VcenterPort = v
	}

	return obj
}
