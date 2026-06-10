package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2PveKind       = "PveConfig"
	machineConfigV2PveAPIVersion = "rke-machine-config.cattle.io/v1"
	machineConfigV2PveAPIType    = "rke-machine-config.cattle.io.pveconfig"
)

type machineConfigV2Pve struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	PveUrl            string `json:"url,omitempty" yaml:"url,omitempty"`
	PveTokenId        string `json:"tokenId,omitempty" yaml:"tokenId,omitempty"`
	PveTokenSecret    string `json:"tokenSecret,omitempty" yaml:"tokenSecret,omitempty"`
	PveInsecureTls    bool   `json:"insecureTls,omitempty" yaml:"insecureTls,omitempty"`
	PveResourcePool   string `json:"resourcePool,omitempty" yaml:"resourcePool,omitempty"`
	PveTemplateId     string `json:"template,omitempty" yaml:"template,omitempty"`
	PveIsoDevice      string `json:"isoDevice" yaml:"isoDevice"`
	PveNetworkIface   string `json:"networkInterface" yaml:"networkInterface"`
	PveSshUser        string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	PveSshPort        string `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	PveProcessorSocks string `json:"processorSockets,omitempty" yaml:"processorSockets,omitempty"`
	PveProcessorCores string `json:"processorCores,omitempty" yaml:"processorCores,omitempty"`
	PveMemory         string `json:"memory,omitempty" yaml:"memory,omitempty"`
	PveMemoryBalloon  string `json:"memoryBalloon,omitempty" yaml:"memoryBalloon,omitempty"`
	PveFullClone      bool   `json:"fullClone,omitempty" yaml:"fullClone,omitempty"`
	PveTags           string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type MachineConfigV2Pve struct {
	norman.Resource
	machineConfigV2Pve
}

// Flatteners

func flattenMachineConfigV2Pve(in *MachineConfigV2Pve) []interface{} {
	if in == nil {
		return nil
	}
	obj := make(map[string]interface{})

	if len(in.PveUrl) > 0 {
		obj["pve_url"] = in.PveUrl
	}
	if len(in.PveTokenId) > 0 {
		obj["pve_token_id"] = in.PveTokenId
	}
	if len(in.PveTokenSecret) > 0 {
		obj["pve_token_secret"] = in.PveTokenSecret
	}
	obj["pve_insecure_tls"] = in.PveInsecureTls
	obj["pve_resource_pool"] = in.PveResourcePool
	if len(in.PveTemplateId) > 0 {
		obj["pve_template_id"] = in.PveTemplateId
	}
	if len(in.PveIsoDevice) > 0 {
		obj["pve_iso_device"] = in.PveIsoDevice
	}
	if len(in.PveNetworkIface) > 0 {
		obj["pve_network_interface"] = in.PveNetworkIface
	}
	if len(in.PveSshUser) > 0 {
		obj["pve_ssh_user"] = in.PveSshUser
	}
	if len(in.PveSshPort) > 0 {
		obj["pve_ssh_port"] = in.PveSshPort
	}
	if len(in.PveProcessorSocks) > 0 {
		obj["pve_processor_sockets"] = in.PveProcessorSocks
	}
	if len(in.PveProcessorCores) > 0 {
		obj["pve_processor_cores"] = in.PveProcessorCores
	}
	if len(in.PveMemory) > 0 {
		obj["pve_memory"] = in.PveMemory
	}
	if len(in.PveMemoryBalloon) > 0 {
		obj["pve_memory_balloon"] = in.PveMemoryBalloon
	}
	obj["pve_full_clone"] = in.PveFullClone
	if len(in.PveTags) > 0 {
		obj["pve_tags"] = in.PveTags
	}

	return []interface{}{obj}
}

// Expanders

func expandMachineConfigV2Pve(p []interface{}, source *MachineConfigV2) *MachineConfigV2Pve {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &MachineConfigV2Pve{}
	if len(source.ID) > 0 {
		obj.ID = source.ID
	}
	in := p[0].(map[string]interface{})

	obj.TypeMeta.Kind = machineConfigV2PveKind
	obj.TypeMeta.APIVersion = machineConfigV2PveAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	if v, ok := in["pve_url"].(string); ok && len(v) > 0 {
		obj.PveUrl = v
	}
	if v, ok := in["pve_token_id"].(string); ok && len(v) > 0 {
		obj.PveTokenId = v
	}
	if v, ok := in["pve_token_secret"].(string); ok && len(v) > 0 {
		obj.PveTokenSecret = v
	}
	if v, ok := in["pve_insecure_tls"].(bool); ok {
		obj.PveInsecureTls = v
	}
	if v, ok := in["pve_resource_pool"].(string); ok {
		obj.PveResourcePool = v
	}
	if v, ok := in["pve_template_id"].(string); ok && len(v) > 0 {
		obj.PveTemplateId = v
	}
	if v, ok := in["pve_iso_device"].(string); ok {
		obj.PveIsoDevice = v
	}
	if v, ok := in["pve_network_interface"].(string); ok {
		obj.PveNetworkIface = v
	}
	if v, ok := in["pve_ssh_user"].(string); ok && len(v) > 0 {
		obj.PveSshUser = v
	}
	if v, ok := in["pve_ssh_port"].(string); ok && len(v) > 0 {
		obj.PveSshPort = v
	}
	if v, ok := in["pve_processor_sockets"].(string); ok && len(v) > 0 {
		obj.PveProcessorSocks = v
	}
	if v, ok := in["pve_processor_cores"].(string); ok && len(v) > 0 {
		obj.PveProcessorCores = v
	}
	if v, ok := in["pve_memory"].(string); ok && len(v) > 0 {
		obj.PveMemory = v
	}
	if v, ok := in["pve_memory_balloon"].(string); ok && len(v) > 0 {
		obj.PveMemoryBalloon = v
	}
	if v, ok := in["pve_full_clone"].(bool); ok {
		obj.PveFullClone = v
	}
	if v, ok := in["pve_tags"].(string); ok && len(v) > 0 {
		obj.PveTags = v
	}

	return obj
}
