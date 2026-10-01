package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"
)

func TestTCloudPublicMachineConfigSchema(t *testing.T) {
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"generate_name": "tcloud-pool",
		"tcloud_public_config": []interface{}{map[string]interface{}{
			"network_scope": "shared", "skip_default_sg": true,
			"vpc_id": "test-vpc", "subnet_id": "test-subnet", "sec_groups": "test-security-group",
		}},
	})
	_, errs := resourceRancher2MachineConfigV2().Validate(config)
	require.Empty(t, errs)

	fields := machineConfigV2Fields()
	block := fields["tcloud_public_config"]
	require.Equal(t, schema.TypeList, block.Type)
	require.True(t, block.Optional)
	require.Equal(t, 1, block.MaxItems)
	for _, other := range allMachineDriverConfigFields {
		if other != "tcloud_public_config" {
			require.Contains(t, block.ConflictsWith, other)
			require.Contains(t, fields[other].ConflictsWith, "tcloud_public_config")
		}
	}
	for _, key := range []string{"access_key", "secret_key", "password", "token", "private_key_file", "user_data_raw"} {
		require.True(t, machineConfigV2OpentelekomcloudFields()[key].Sensitive, key)
	}
	require.Nil(t, resourceRancher2MachineConfigV2().Importer)
}

func TestTCloudPublicMachineConfigDefaults(t *testing.T) {
	d := schema.TestResourceDataRaw(t, machineConfigV2Fields(), map[string]interface{}{
		"generate_name": "tcloud-pool", "tcloud_public_config": []interface{}{map[string]interface{}{}},
	})
	expected := map[string]interface{}{
		"bandwidth_size": "100", "bandwidth_type": "PER", "eip_type": "5_bgp",
		"endpoint_type": "public", "flavor_name": "s3.xlarge.2",
		"image_name": "Standard_Ubuntu_24.04_amd64_uefi_latest",
		"ip_version": "4", "network_scope": "machine",
		"root_volume_size": "40", "root_volume_type": "SSD",
		"ssh_port": "22", "ssh_user": "ubuntu",
		"subnet_name": "subnet-docker-machine", "vpc_name": "vpc-docker-machine",
		"skip_default_sg": false, "skip_eip": false,
	}
	for key := range machineConfigV2OpentelekomcloudFields() {
		if _, ok := expected[key]; !ok {
			expected[key] = ""
		}
	}
	actual := d.Get("tcloud_public_config").([]interface{})[0]
	require.Equal(t, expected, actual)
	require.Equal(t, []interface{}{expected}, flattenMachineConfigV2Opentelekomcloud(expandMachineConfigV2(d).OpentelekomcloudConfig))
}

func TestTCloudPublicMachineConfigInvalidConfiguration(t *testing.T) {
	for name, values := range map[string]map[string]interface{}{
		"network scope": {"network_scope": "unknown"},
		"IP version":    {"ip_version": "5"},
	} {
		t.Run(name, func(t *testing.T) {
			_, errs := resourceRancher2MachineConfigV2().Validate(terraform.NewResourceConfigRaw(map[string]interface{}{
				"generate_name": "tcloud-pool", "tcloud_public_config": []interface{}{values},
			}))
			require.NotEmpty(t, errs)
		})
	}
	_, errs := resourceRancher2MachineConfigV2().Validate(terraform.NewResourceConfigRaw(map[string]interface{}{
		"generate_name":        "tcloud-pool",
		"tcloud_public_config": []interface{}{map[string]interface{}{}},
		"amazonec2_config":     []interface{}{map[string]interface{}{}},
	}))
	require.NotEmpty(t, errs)
}
