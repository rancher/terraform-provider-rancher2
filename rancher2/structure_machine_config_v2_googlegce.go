package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2GoogleGCEKind       = "GoogleConfig"
	machineConfigV2GoogleGCEAPIVersion = "rke-machine-config.cattle.io/v1"
	machineConfigV2GoogleGCEAPIType    = "rke-machine-config.cattle.io.googleconfig"
)

type machineConfigV2GoogleGCE struct {
	metav1.TypeMeta            `json:",inline"`
	metav1.ObjectMeta          `json:"metadata,omitempty"`
	Address                    string   `json:"address,omitempty" yaml:"address,omitempty"`
	AuthEncodedJson            string   `json:"authEncodedJson,omitempty" yaml:"authEncodedJson,omitempty"`
	DiskSize                   string   `json:"diskSize,omitempty" yaml:"diskSize,omitempty"`
	DiskType                   string   `json:"diskType,omitempty" yaml:"diskType,omitempty"`
	ExternalFirewallRulePrefix string   `json:"externalFirewallRulePrefix,omitempty" yaml:"externalFirewallRulePrefix,omitempty"`
	InternalFirewallRulePrefix string   `json:"internalFirewallRulePrefix,omitempty" yaml:"internalFirewallRulePrefix,omitempty"`
	Labels                     string   `json:"labels,omitempty" yaml:"labels,omitempty"`
	MachineImage               string   `json:"machineImage,omitempty" yaml:"machineImage,omitempty"`
	MachineType                string   `json:"machineType,omitempty" yaml:"machineType,omitempty"`
	Network                    string   `json:"network,omitempty" yaml:"network,omitempty"`
	OpenPort                   []string `json:"openPort,omitempty" yaml:"openPort,omitempty"`
	Preemptable                bool     `json:"preemptable,omitempty" yaml:"preemptable,omitempty"`
	Project                    string   `json:"project,omitempty" yaml:"project,omitempty"`
	Scopes                     string   `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	SubNetwork                 string   `json:"subNetwork,omitempty" yaml:"subNetwork,omitempty"`
	Tags                       string   `json:"tags,omitempty" yaml:"tags,omitempty"`
	UseExisting                bool     `json:"useExisting,omitempty" yaml:"useExisting,omitempty"`
	UseInternalIP              bool     `json:"useInternalIP,omitempty" yaml:"useInternalIP,omitempty"`
	UseInternalIPOnly          bool     `json:"useInternalIPOnly,omitempty" yaml:"useInternalIPOnly,omitempty"`
	Userdata                   string   `json:"userdata,omitempty" yaml:"userdata,omitempty"`
	Zone                       string   `json:"zone,omitempty" yaml:"zone,omitempty"`
}

type MachineConfigV2GoogleGCE struct {
	norman.Resource
	machineConfigV2GoogleGCE
}

func flattenMachineConfigV2GoogleGCE(in *MachineConfigV2GoogleGCE) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.Address) > 0 {
		obj["address"] = in.Address
	}

	if len(in.AuthEncodedJson) > 0 {
		obj["auth_encoded_json"] = in.AuthEncodedJson
	}

	if len(in.DiskSize) > 0 {
		obj["disk_size"] = in.DiskSize
	}

	if len(in.DiskType) > 0 {
		obj["disk_type"] = in.DiskType
	}

	if len(in.ExternalFirewallRulePrefix) > 0 {
		obj["external_firewall_rule_prefix"] = in.ExternalFirewallRulePrefix
	}

	if len(in.InternalFirewallRulePrefix) > 0 {
		obj["internal_firewall_rule_prefix"] = in.InternalFirewallRulePrefix
	}

	if len(in.Labels) > 0 {
		obj["labels"] = in.Labels
	}

	if len(in.MachineImage) > 0 {
		obj["machine_image"] = in.MachineImage
	}

	if len(in.MachineType) > 0 {
		obj["machine_type"] = in.MachineType
	}

	if len(in.Network) > 0 {
		obj["network"] = in.Network
	}

	if len(in.OpenPort) > 0 {
		obj["open_port"] = in.OpenPort
	}

	if in.Preemptable {
		obj["preemptable"] = in.Preemptable
	}

	if len(in.Project) > 0 {
		obj["project"] = in.Project
	}

	if len(in.Scopes) > 0 {
		obj["scopes"] = in.Scopes
	}

	if len(in.SubNetwork) > 0 {
		obj["sub_network"] = in.SubNetwork
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if in.UseExisting {
		obj["use_existing"] = in.UseExisting
	}

	if in.UseInternalIP {
		obj["use_internal_ip"] = in.UseInternalIP
	}

	if in.UseInternalIPOnly {
		obj["use_internal_ip_only"] = in.UseInternalIPOnly
	}

	if len(in.Userdata) > 0 {
		obj["userdata"] = in.Userdata
	}

	if len(in.Zone) > 0 {
		obj["zone"] = in.Zone
	}

	return []interface{}{obj}
}

func expandMachineConfigV2GoogleGCE(p []interface{}, source *MachineConfigV2) *MachineConfigV2GoogleGCE {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &MachineConfigV2GoogleGCE{}
	if len(source.ID) > 0 {
		obj.ID = source.ID
	}

	in := p[0].(map[string]interface{})

	obj.TypeMeta.Kind = machineConfigV2GoogleGCEKind
	obj.TypeMeta.APIVersion = machineConfigV2GoogleGCEAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	if v, ok := in["address"].(string); ok && len(v) > 0 {
		obj.Address = v
	}

	if v, ok := in["auth_encoded_json"].(string); ok && len(v) > 0 {
		obj.AuthEncodedJson = v
	}

	if v, ok := in["disk_size"].(string); ok && len(v) > 0 {
		obj.DiskSize = v
	}

	if v, ok := in["disk_type"].(string); ok && len(v) > 0 {
		obj.DiskType = v
	}

	if v, ok := in["external_firewall_rule_prefix"].(string); ok && len(v) > 0 {
		obj.ExternalFirewallRulePrefix = v
	}

	if v, ok := in["internal_firewall_rule_prefix"].(string); ok && len(v) > 0 {
		obj.InternalFirewallRulePrefix = v
	}

	if v, ok := in["labels"].(string); ok && len(v) > 0 {
		obj.Labels = v
	}

	if v, ok := in["machine_image"].(string); ok && len(v) > 0 {
		obj.MachineImage = v
	}

	if v, ok := in["machine_type"].(string); ok && len(v) > 0 {
		obj.MachineType = v
	}

	if v, ok := in["network"].(string); ok && len(v) > 0 {
		obj.Network = v
	}

	if v, ok := in["open_port"].([]interface{}); ok && len(v) > 0 {
		obj.OpenPort = toArrayString(v)
	}

	if v, ok := in["preemptable"].(bool); ok {
		obj.Preemptable = v
	}

	if v, ok := in["project"].(string); ok && len(v) > 0 {
		obj.Project = v
	}

	if v, ok := in["scopes"].(string); ok && len(v) > 0 {
		obj.Scopes = v
	}

	if v, ok := in["sub_network"].(string); ok && len(v) > 0 {
		obj.SubNetwork = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["use_existing"].(bool); ok {
		obj.UseExisting = v
	}

	if v, ok := in["use_internal_ip"].(bool); ok {
		obj.UseInternalIP = v
	}

	if v, ok := in["use_internal_ip_only"].(bool); ok {
		obj.UseInternalIPOnly = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.Userdata = v
	}

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		obj.Zone = v
	}

	return obj
}
