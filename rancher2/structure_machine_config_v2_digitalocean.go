package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2DigitaloceanKind         = "DigitaloceanConfig"
	machineConfigV2DigitaloceanAPIVersion   = "rke-machine-config.cattle.io/v1"
	machineConfigV2DigitaloceanAPIType      = "rke-machine-config.cattle.io.digitaloceanconfig"
	machineConfigV2DigitaloceanClusterIDsep = "."
)

//Types

type machineConfigV2Digitalocean struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AccessToken       string `json:"accessToken,omitempty" yaml:"accessToken,omitempty"`
	Backups           bool   `json:"backups,omitempty" yaml:"backups,omitempty"`
	Image             string `json:"image,omitempty" yaml:"image,omitempty"`
	IPV6              bool   `json:"ipv6,omitempty" yaml:"ipv6,omitempty"`
	Monitoring        bool   `json:"monitoring,omitempty" yaml:"monitoring,omitempty"`
	PrivateNetworking bool   `json:"privateNetworking,omitempty" yaml:"privateNetworking,omitempty"`
	Region            string `json:"region,omitempty" yaml:"region,omitempty"`
	Size              string `json:"size,omitempty" yaml:"size,omitempty"`
	SSHKeyContents    string `json:"sshKeyContents,omitempty" yaml:"sshKeyContents,omitempty"`
	SSHKeyFingerprint string `json:"sshKeyFingerprint,omitempty" yaml:"sshKeyFingerprint,omitempty"`
	SSHPort           string `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser           string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	Tags              string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Userdata          string `json:"userdata,omitempty" yaml:"userdata,omitempty"`
}

type MachineConfigV2Digitalocean struct {
	norman.Resource
	machineConfigV2Digitalocean
}

// Flatteners

func flattenMachineConfigV2Digitalocean(in *MachineConfigV2Digitalocean) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.AccessToken) > 0 {
		obj["access_token"] = in.AccessToken
	}

	obj["backups"] = in.Backups

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["ipv6"] = in.IPV6
	obj["monitoring"] = in.Monitoring
	obj["private_networking"] = in.PrivateNetworking

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.Size) > 0 {
		obj["size"] = in.Size
	}

	if len(in.SSHKeyContents) > 0 {
		obj["ssh_key_contents"] = in.SSHKeyContents
	}

	if len(in.SSHKeyFingerprint) > 0 {
		obj["ssh_key_fingerprint"] = in.SSHKeyFingerprint
	}

	if len(in.SSHPort) > 0 {
		obj["ssh_port"] = in.SSHPort
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if len(in.Userdata) > 0 {
		obj["userdata"] = in.Userdata
	}

	return []interface{}{obj}
}

// Expanders

func expandMachineConfigV2Digitalocean(p []interface{}, source *MachineConfigV2) *MachineConfigV2Digitalocean {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &MachineConfigV2Digitalocean{}

	if len(source.ID) > 0 {
		obj.ID = source.ID
	}
	in := p[0].(map[string]interface{})

	obj.TypeMeta.Kind = machineConfigV2DigitaloceanKind
	obj.TypeMeta.APIVersion = machineConfigV2DigitaloceanAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	if v, ok := in["access_token"].(string); ok && len(v) > 0 {
		obj.AccessToken = v
	}

	if v, ok := in["backups"].(bool); ok {
		obj.Backups = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["ipv6"].(bool); ok {
		obj.IPV6 = v
	}

	if v, ok := in["monitoring"].(bool); ok {
		obj.Monitoring = v
	}
	if v, ok := in["private_networking"].(bool); ok {
		obj.PrivateNetworking = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["size"].(string); ok && len(v) > 0 {
		obj.Size = v
	}

	if v, ok := in["ssh_key_contents"].(string); ok && len(v) > 0 {
		obj.SSHKeyContents = v
	}

	if v, ok := in["ssh_key_fingerprint"].(string); ok && len(v) > 0 {
		obj.SSHKeyFingerprint = v
	}

	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.Userdata = v
	}

	return obj
}
