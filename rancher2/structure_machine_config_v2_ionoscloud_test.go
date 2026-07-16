package rancher2

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testMachineConfigV2IonoscloudConf = &MachineConfigV2Ionoscloud{
		machineConfigV2Ionoscloud: machineConfigV2Ionoscloud{
			TypeMeta: metav1.TypeMeta{
				Kind:       machineConfigV2IonoscloudKind,
				APIVersion: machineConfigV2IonoscloudAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ic-test",
				Namespace: "fleet-default",
			},
			Endpoint:               "https://api.ionos.com/cloudapi/v6",
			Username:               "user@example.com",
			Password:               "secret",
			ServerCores:            4,
			ServerRam:              8192,
			ServerCpuFamily:        "INTEL_SKYLAKE",
			ServerAvailabilityZone: "AUTO",
			DiskSize:               50,
			DiskType:               "SSD",
			ServerType:             "ENTERPRISE",
			Image:                  "ubuntu:20.04",
			ImagePassword:          "imgpass",
			Location:               "de/txl",
			DatacenterName:         "my-dc",
			NicDhcp:                true,
			NicIps:                 []string{"10.0.0.5"},
			LanName:                "my-lan",
			VolumeAvailabilityZone: "AUTO",
			CloudInit:              "#cloud-config",
			SshInCloudInit:         true,
			SshUser:                "root",
			WaitForIpChange:        true,
			WaitForIpChangeTimeout: 300,
			NatName:                "my-nat",
			NatPublicIps:           []string{"1.2.3.4"},
			NatFlowlogs:            []string{"flow1"},
			NatRules:               []string{"rule1"},
			NatLansToGateways:      "1:10.0.0.1",
			PrivateLan:             true,
			AdditionalLans:         []string{"2"},
			CreateNat:              true,
			AppendRkeCloudInit:     true,
		},
	}
	testMachineConfigV2IonoscloudInterface = []any{
		map[string]any{
			"endpoint":                   "https://api.ionos.com/cloudapi/v6",
			"username":                   "user@example.com",
			"password":                   "secret",
			"server_cores":               4,
			"server_ram":                 8192,
			"server_cpu_family":          "INTEL_SKYLAKE",
			"server_availability_zone":   "AUTO",
			"disk_size":                  50,
			"disk_type":                  "SSD",
			"server_type":                "ENTERPRISE",
			"image":                      "ubuntu:20.04",
			"image_password":             "imgpass",
			"location":                   "de/txl",
			"datacenter_name":            "my-dc",
			"nic_dhcp":                   true,
			"nic_ips":                    []any{"10.0.0.5"},
			"lan_name":                   "my-lan",
			"volume_availability_zone":   "AUTO",
			"cloud_init":                 "#cloud-config",
			"ssh_in_cloud_init":          true,
			"ssh_user":                   "root",
			"wait_for_ip_change":         true,
			"wait_for_ip_change_timeout": 300,
			"nat_name":                   "my-nat",
			"nat_public_ips":             []any{"1.2.3.4"},
			"nat_flowlogs":               []any{"flow1"},
			"nat_rules":                  []any{"rule1"},
			"skip_default_nat_rules":     false,
			"nat_lans_to_gateways":       "1:10.0.0.1",
			"private_lan":                true,
			"additional_lans":            []any{"2"},
			"create_nat":                 true,
			"append_rke_cloud_init":      true,
		},
	}
)

func testIonoscloudSource() *MachineConfigV2 {
	return &MachineConfigV2{
		machineConfigV2: machineConfigV2{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ic-test",
				Namespace: "fleet-default",
			},
		},
	}
}

func TestFlattenMachineConfigV2Ionoscloud(t *testing.T) {
	output := flattenMachineConfigV2Ionoscloud(testMachineConfigV2IonoscloudConf)
	assert.Equal(t, testMachineConfigV2IonoscloudInterface, output, "Unexpected output from flattener.")
}

func TestExpandMachineConfigV2Ionoscloud(t *testing.T) {
	output := expandMachineConfigV2Ionoscloud(testMachineConfigV2IonoscloudInterface, testIonoscloudSource())
	assert.Equal(t, testMachineConfigV2IonoscloudConf.machineConfigV2Ionoscloud, output.machineConfigV2Ionoscloud, "Unexpected output from expander.")
}

// The wire (JSON) field names must match the ionoscloud rke-machine-config CRD
// (cores/ram/cpuFamily — NOT serverCores/serverRam), otherwise Rancher silently
// drops the values and they never round-trip (state reads back 0 -> permadiff).
func TestMachineConfigV2IonoscloudJSONFieldNames(t *testing.T) {
	b, err := json.Marshal(testMachineConfigV2IonoscloudConf.machineConfigV2Ionoscloud)
	assert.NoError(t, err)
	j := string(b)
	assert.Contains(t, j, `"cores":"4"`, "server cores must marshal as CRD field `cores`")
	assert.Contains(t, j, `"ram":"8192"`, "server ram must marshal as CRD field `ram`")
	assert.Contains(t, j, `"cpuFamily":`, "cpu family must marshal as CRD field `cpuFamily`")
	assert.NotContains(t, j, "serverCores", "stale wrong field name must be gone")
	assert.NotContains(t, j, "serverRam", "stale wrong field name must be gone")
}

// When an *_id is provided, the matching *_name must be dropped so a stale
// name is never sent alongside an id (regression guard for the id/name switch).
func TestExpandMachineConfigV2IonoscloudPrefersIDs(t *testing.T) {
	input := []any{
		map[string]any{
			"datacenter_id":   "dc-uuid",
			"datacenter_name": "my-dc",
			"lan_id":          "lan-uuid",
			"lan_name":        "my-lan",
			"nat_id":          "nat-uuid",
			"nat_name":        "my-nat",
		},
	}
	output := expandMachineConfigV2Ionoscloud(input, testIonoscloudSource())
	assert.Equal(t, "dc-uuid", output.DatacenterId)
	assert.Empty(t, output.DatacenterName, "datacenter_name must be dropped when datacenter_id is set")
	assert.Equal(t, "lan-uuid", output.LanId)
	assert.Empty(t, output.LanName, "lan_name must be dropped when lan_id is set")
	assert.Equal(t, "nat-uuid", output.NatId)
	assert.Empty(t, output.NatName, "nat_name must be dropped when nat_id is set")
}
