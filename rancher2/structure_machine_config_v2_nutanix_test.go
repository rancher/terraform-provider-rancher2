package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testMachineConfigV2NutanixConf = &MachineConfigV2Nutanix{
		machineConfigV2Nutanix: machineConfigV2Nutanix{
			TypeMeta: metav1.TypeMeta{
				Kind:       machineConfigV2NutanixKind,
				APIVersion: machineConfigV2NutanixAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nc-test",
				Namespace: "fleet-default",
			},
			Endpoint:         "pc.example.com",
			Username:         "X-ntnx-api-key",
			Password:         "secret",
			Port:             "9440",
			Insecure:         true,
			Cluster:          "cluster-a",
			VMMem:            "4096",
			VMCPUs:           "4",
			VMCores:          "2",
			VMCPUPassthrough: true,
			VMNetwork:        []string{"subnet-a", "subnet-b"},
			VMImage:          "ubuntu-cloudinit",
			VMImageSize:      "80",
			VMCategories:     []string{"app:tf", "env:dev"},
			StorageContainer: "container-uuid",
			DiskSize:         "40",
			CloudInit:        "#cloud-config",
			VMSerialPort:     true,
			Project:          "project-a",
			BootType:         "uefi",
			Timeout:          "600",
			VMGPU:            []string{"gpu-a"},
		},
	}
	testMachineConfigV2NutanixInterface = []any{
		map[string]any{
			"endpoint":           "pc.example.com",
			"username":           "X-ntnx-api-key",
			"password":           "secret",
			"port":               "9440",
			"insecure":           true,
			"cluster":            "cluster-a",
			"vm_mem":             "4096",
			"vm_cpus":            "4",
			"vm_cores":           "2",
			"vm_cpu_passthrough": true,
			"vm_network":         []any{"subnet-a", "subnet-b"},
			"vm_image":           "ubuntu-cloudinit",
			"vm_image_size":      "80",
			"vm_categories":      []any{"app:tf", "env:dev"},
			"storage_container":  "container-uuid",
			"disk_size":          "40",
			"cloud_init":         "#cloud-config",
			"vm_serial_port":     true,
			"project":            "project-a",
			"boot_type":          "uefi",
			"timeout":            "600",
			"vm_gpu":             []any{"gpu-a"},
		},
	}
)

func TestFlattenMachineConfigV2Nutanix(t *testing.T) {
	cases := []struct {
		Input          *MachineConfigV2Nutanix
		ExpectedOutput []any
	}{
		{
			Input:          testMachineConfigV2NutanixConf,
			ExpectedOutput: testMachineConfigV2NutanixInterface,
		},
	}

	for _, tc := range cases {
		output := flattenMachineConfigV2Nutanix(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandMachineConfigV2Nutanix(t *testing.T) {
	cases := []struct {
		Input          []any
		Source         *MachineConfigV2
		ExpectedOutput *MachineConfigV2Nutanix
	}{
		{
			Input: testMachineConfigV2NutanixInterface,
			Source: &MachineConfigV2{
				machineConfigV2: machineConfigV2{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nc-test",
						Namespace: "fleet-default",
					},
				},
			},
			ExpectedOutput: testMachineConfigV2NutanixConf,
		},
	}

	for _, tc := range cases {
		output := expandMachineConfigV2Nutanix(tc.Input, tc.Source)
		assert.Equal(t, tc.ExpectedOutput.machineConfigV2Nutanix, output.machineConfigV2Nutanix, "Unexpected output from expander.")
	}
}
