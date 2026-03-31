package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenMachineConfigV2Linode(t *testing.T) {
	input := &MachineConfigV2Linode{
		machineConfigV2Linode: machineConfigV2Linode{
			AuthorizedUsers:           "user1,user2",
			CreatePrivateIP:           false,
			DockerPort:                "2377",
			Image:                     "linode/ubuntu22.04",
			InstanceType:              "g6-standard-2",
			Label:                     "example",
			Region:                    "us-east",
			RootPass:                  "secret",
			SSHPort:                   "2222",
			SSHUser:                   "ubuntu",
			StackScript:               "user/stack",
			StackscriptData:           "{\"foo\":\"bar\"}",
			UserData:                  "#cloud-config\npackages:\n- htop\n",
			SwapSize:                  "256",
			Tags:                      "tag1,tag2",
			Token:                     "token",
			UAPrefix:                  "tf-provider/1.0",
			UseInterfaces:             true,
			VPCSubnetID:               "456",
			VPCPrivateIP:              "10.0.0.10",
			PublicInterfaceFirewallID: "789",
			VPCInterfaceFirewallID:    "321",
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			"authorized_users":             "user1,user2",
			"create_private_ip":            false,
			"docker_port":                  "2377",
			"image":                        "linode/ubuntu22.04",
			"instance_type":                "g6-standard-2",
			"label":                        "example",
			"region":                       "us-east",
			"root_pass":                    "secret",
			"ssh_port":                     "2222",
			"ssh_user":                     "ubuntu",
			"stackscript":                  "user/stack",
			"stackscript_data":             "{\"foo\":\"bar\"}",
			"user_data":                    "#cloud-config\npackages:\n- htop\n",
			"swap_size":                    "256",
			"tags":                         "tag1,tag2",
			"token":                        "token",
			"use_interfaces":               true,
			"vpc_subnet_id":                "456",
			"vpc_private_ip":               "10.0.0.10",
			"public_interface_firewall_id": "789",
			"vpc_interface_firewall_id":    "321",
			"ua_prefix":                    "tf-provider/1.0",
		},
	}

	assert.Equal(t, expected, flattenMachineConfigV2Linode(input))
}

func TestExpandMachineConfigV2Linode(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"authorized_users":             "user1,user2",
			"create_private_ip":            false,
			"docker_port":                  "2377",
			"image":                        "linode/ubuntu22.04",
			"instance_type":                "g6-standard-2",
			"label":                        "example",
			"region":                       "us-east",
			"root_pass":                    "secret",
			"ssh_port":                     "2222",
			"ssh_user":                     "ubuntu",
			"stackscript":                  "user/stack",
			"stackscript_data":             "{\"foo\":\"bar\"}",
			"user_data":                    "#cloud-config\npackages:\n- htop\n",
			"swap_size":                    "256",
			"tags":                         "tag1,tag2",
			"token":                        "token",
			"use_interfaces":               true,
			"vpc_subnet_id":                "456",
			"vpc_private_ip":               "10.0.0.10",
			"public_interface_firewall_id": "789",
			"vpc_interface_firewall_id":    "321",
			"ua_prefix":                    "tf-provider/1.0",
		},
	}

	source := &MachineConfigV2{}
	got := expandMachineConfigV2Linode(input, source)

	expected := &MachineConfigV2Linode{
		machineConfigV2Linode: machineConfigV2Linode{
			AuthorizedUsers:           "user1,user2",
			CreatePrivateIP:           false,
			DockerPort:                "2377",
			Image:                     "linode/ubuntu22.04",
			InstanceType:              "g6-standard-2",
			Label:                     "example",
			Region:                    "us-east",
			RootPass:                  "secret",
			SSHPort:                   "2222",
			SSHUser:                   "ubuntu",
			StackScript:               "user/stack",
			StackscriptData:           "{\"foo\":\"bar\"}",
			UserData:                  "#cloud-config\npackages:\n- htop\n",
			SwapSize:                  "256",
			Tags:                      "tag1,tag2",
			Token:                     "token",
			UAPrefix:                  "tf-provider/1.0",
			UseInterfaces:             true,
			VPCSubnetID:               "456",
			VPCPrivateIP:              "10.0.0.10",
			PublicInterfaceFirewallID: "789",
			VPCInterfaceFirewallID:    "321",
		},
	}
	expected.TypeMeta.Kind = machineConfigV2LinodeKind
	expected.TypeMeta.APIVersion = machineConfigV2LinodeAPIVersion

	assert.Equal(t, expected, got)
	assert.Equal(t, expected.TypeMeta, source.TypeMeta)
}
