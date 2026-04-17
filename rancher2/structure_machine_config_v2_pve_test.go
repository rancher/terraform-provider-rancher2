package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testMachineConfigV2PveConf = map[string]interface{}{
		"pveUrl":              "https://pve.example.com:8006",
		"pveTokenId":          "root@pam!rancher",
		"pveTokenSecret":      "secret-uuid",
		"pveInsecureTls":      false,
		"pveResourcePool":     "rancher",
		"pveTemplateId":       float64(100),
		"pveIsoDevice":        "scsi1",
		"pveNetworkInterface": "net0",
		"pveSshUser":          "ubuntu",
		"pveSshPort":          float64(22),
		"pveProcessorSockets": "2",
		"pveProcessorCores":   "4",
		"pveMemory":           "4096",
	}
	testMachineConfigV2PveInterface = []interface{}{
		map[string]interface{}{
			"pve_url":               "https://pve.example.com:8006",
			"pve_token_id":          "root@pam!rancher",
			"pve_token_secret":      "secret-uuid",
			"pve_insecure_tls":      false,
			"pve_resource_pool":     "rancher",
			"pve_template_id":       100,
			"pve_iso_device":        "scsi1",
			"pve_network_interface": "net0",
			"pve_ssh_user":          "ubuntu",
			"pve_ssh_port":          22,
			"pve_processor_sockets": "2",
			"pve_processor_cores":   "4",
			"pve_memory":            "4096",
		},
	}
)

func TestFlattenMachineConfigV2Pve(t *testing.T) {
	result := flattenMachineConfigV2Pve(testMachineConfigV2PveConf, []interface{}{})
	assert.Equal(t, testMachineConfigV2PveInterface, result)
}

func TestExpandMachineConfigV2Pve(t *testing.T) {
	result := expandMachineConfigV2Pve(testMachineConfigV2PveInterface)
	assert.Equal(t, testMachineConfigV2PveConf, result)
}
