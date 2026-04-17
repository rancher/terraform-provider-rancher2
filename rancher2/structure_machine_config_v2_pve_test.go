package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testMachineConfigV2PveConf = &MachineConfigV2Pve{
		machineConfigV2Pve: machineConfigV2Pve{
			PveUrl:            "https://pve.example.com:8006",
			PveTokenId:        "root@pam!rancher",
			PveTokenSecret:    "secret-uuid",
			PveInsecureTls:    false,
			PveResourcePool:   "rancher",
			PveTemplateId:     100,
			PveIsoDevice:      "scsi1",
			PveNetworkIface:   "net0",
			PveSshUser:        "ubuntu",
			PveSshPort:        22,
			PveProcessorSocks: "2",
			PveProcessorCores: "4",
			PveMemory:         "4096",
		},
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
	result := flattenMachineConfigV2Pve(testMachineConfigV2PveConf)
	assert.Equal(t, testMachineConfigV2PveInterface, result)
}

func TestExpandMachineConfigV2Pve(t *testing.T) {
	source := &MachineConfigV2{}
	result := expandMachineConfigV2Pve(testMachineConfigV2PveInterface, source)
	assert.Equal(t, testMachineConfigV2PveConf.PveUrl, result.PveUrl)
	assert.Equal(t, testMachineConfigV2PveConf.PveTokenId, result.PveTokenId)
	assert.Equal(t, testMachineConfigV2PveConf.PveTokenSecret, result.PveTokenSecret)
	assert.Equal(t, testMachineConfigV2PveConf.PveInsecureTls, result.PveInsecureTls)
	assert.Equal(t, testMachineConfigV2PveConf.PveResourcePool, result.PveResourcePool)
	assert.Equal(t, testMachineConfigV2PveConf.PveTemplateId, result.PveTemplateId)
	assert.Equal(t, testMachineConfigV2PveConf.PveIsoDevice, result.PveIsoDevice)
	assert.Equal(t, testMachineConfigV2PveConf.PveNetworkIface, result.PveNetworkIface)
	assert.Equal(t, testMachineConfigV2PveConf.PveSshUser, result.PveSshUser)
	assert.Equal(t, testMachineConfigV2PveConf.PveSshPort, result.PveSshPort)
	assert.Equal(t, testMachineConfigV2PveConf.PveProcessorSocks, result.PveProcessorSocks)
	assert.Equal(t, testMachineConfigV2PveConf.PveProcessorCores, result.PveProcessorCores)
	assert.Equal(t, testMachineConfigV2PveConf.PveMemory, result.PveMemory)
}
