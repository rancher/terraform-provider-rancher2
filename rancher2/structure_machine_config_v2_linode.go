package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2LinodeKind         = "LinodeConfig"
	machineConfigV2LinodeAPIVersion   = "rke-machine-config.cattle.io/v1"
	machineConfigV2LinodeAPIType      = "rke-machine-config.cattle.io.linodeconfig"
	machineConfigV2LinodeClusterIDsep = "."
)

//Types

type machineConfigV2Linode struct {
	metav1.TypeMeta           `json:",inline"`
	metav1.ObjectMeta         `json:"metadata,omitempty"`
	AuthorizedUsers           string `json:"authorizedUsers,omitempty" yaml:"authorizedUsers,omitempty"`
	CreatePrivateIP           bool   `json:"createPrivateIp,omitempty" yaml:"createPrivateIp,omitempty"`
	DockerPort                string `json:"dockerPort,omitempty" yaml:"dockerPort,omitempty"`
	Image                     string `json:"image,omitempty" yaml:"image,omitempty"`
	InstanceType              string `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
	Label                     string `json:"label,omitempty" yaml:"label,omitempty"`
	Region                    string `json:"region,omitempty" yaml:"region,omitempty"`
	RootPass                  string `json:"rootPass,omitempty" yaml:"rootPass,omitempty"`
	SSHPort                   string `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser                   string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	StackScript               string `json:"stackscript,omitempty" yaml:"stackscript,omitempty"`
	StackscriptData           string `json:"stackscriptData,omitempty" yaml:"stackscriptData,omitempty"`
	SwapSize                  string `json:"swapSize,omitempty" yaml:"swapSize,omitempty"`
	Tags                      string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Token                     string `json:"token,omitempty" yaml:"token,omitempty"`
	UserData                  string `json:"userData,omitempty" yaml:"userData,omitempty"`
	UAPrefix                  string `json:"uaPrefix,omitempty" yaml:"uaPrefix,omitempty"`
	UseInterfaces             bool   `json:"useInterfaces,omitempty" yaml:"useInterfaces,omitempty"`
	VPCInterfaceFirewallID    string `json:"vpcInterfaceFirewallId,omitempty" yaml:"vpcInterfaceFirewallId,omitempty"`
	VPCPrivateIP              string `json:"vpcPrivateIp,omitempty" yaml:"vpcPrivateIp,omitempty"`
	VPCSubnetID               string `json:"vpcSubnetId,omitempty" yaml:"vpcSubnetId,omitempty"`
	PublicInterfaceFirewallID string `json:"publicInterfaceFirewallId,omitempty" yaml:"publicInterfaceFirewallId,omitempty"`
}

type MachineConfigV2Linode struct {
	norman.Resource
	machineConfigV2Linode
}

// Flatteners

func flattenMachineConfigV2Linode(in *MachineConfigV2Linode) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.AuthorizedUsers) > 0 {
		obj["authorized_users"] = in.AuthorizedUsers
	}

	obj["create_private_ip"] = in.CreatePrivateIP

	if len(in.DockerPort) > 0 {
		obj["docker_port"] = in.DockerPort
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.InstanceType) > 0 {
		obj["instance_type"] = in.InstanceType
	}

	if len(in.Label) > 0 {
		obj["label"] = in.Label
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.RootPass) > 0 {
		obj["root_pass"] = in.RootPass
	}

	if len(in.SSHPort) > 0 {
		obj["ssh_port"] = in.SSHPort
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.StackScript) > 0 {
		obj["stackscript"] = in.StackScript
	}

	if len(in.StackscriptData) > 0 {
		obj["stackscript_data"] = in.StackscriptData
	}

	if len(in.UserData) > 0 {
		obj["user_data"] = in.UserData
	}

	if len(in.SwapSize) > 0 {
		obj["swap_size"] = in.SwapSize
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}

	obj["use_interfaces"] = in.UseInterfaces

	if len(in.VPCSubnetID) > 0 {
		obj["vpc_subnet_id"] = in.VPCSubnetID
	}

	if len(in.VPCPrivateIP) > 0 {
		obj["vpc_private_ip"] = in.VPCPrivateIP
	}

	if len(in.PublicInterfaceFirewallID) > 0 {
		obj["public_interface_firewall_id"] = in.PublicInterfaceFirewallID
	}

	if len(in.VPCInterfaceFirewallID) > 0 {
		obj["vpc_interface_firewall_id"] = in.VPCInterfaceFirewallID
	}

	if len(in.UAPrefix) > 0 {
		obj["ua_prefix"] = in.UAPrefix
	}

	return []interface{}{obj}
}

// Expanders

func expandMachineConfigV2Linode(p []interface{}, source *MachineConfigV2) *MachineConfigV2Linode {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &MachineConfigV2Linode{}

	if len(source.ID) > 0 {
		obj.ID = source.ID
	}
	in := p[0].(map[string]interface{})

	obj.TypeMeta.Kind = machineConfigV2LinodeKind
	obj.TypeMeta.APIVersion = machineConfigV2LinodeAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	if v, ok := in["authorized_users"].(string); ok && len(v) > 0 {
		obj.AuthorizedUsers = v
	}

	if v, ok := in["create_private_ip"].(bool); ok {
		obj.CreatePrivateIP = v
	}

	if v, ok := in["docker_port"].(string); ok && len(v) > 0 {
		obj.DockerPort = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["instance_type"].(string); ok && len(v) > 0 {
		obj.InstanceType = v
	}

	if v, ok := in["label"].(string); ok && len(v) > 0 {
		obj.Label = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["root_pass"].(string); ok && len(v) > 0 {
		obj.RootPass = v
	}

	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["stackscript"].(string); ok && len(v) > 0 {
		obj.StackScript = v
	}

	if v, ok := in["stackscript_data"].(string); ok && len(v) > 0 {
		obj.StackscriptData = v
	}

	if v, ok := in["user_data"].(string); ok && len(v) > 0 {
		obj.UserData = v
	}

	if v, ok := in["swap_size"].(string); ok && len(v) > 0 {
		obj.SwapSize = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	if v, ok := in["use_interfaces"].(bool); ok {
		obj.UseInterfaces = v
	}

	if v, ok := in["vpc_subnet_id"].(string); ok && len(v) > 0 {
		obj.VPCSubnetID = v
	}

	if v, ok := in["vpc_private_ip"].(string); ok && len(v) > 0 {
		obj.VPCPrivateIP = v
	}

	if v, ok := in["public_interface_firewall_id"].(string); ok && len(v) > 0 {
		obj.PublicInterfaceFirewallID = v
	}

	if v, ok := in["vpc_interface_firewall_id"].(string); ok && len(v) > 0 {
		obj.VPCInterfaceFirewallID = v
	}

	if v, ok := in["ua_prefix"].(string); ok && len(v) > 0 {
		obj.UAPrefix = v
	}

	return obj
}
