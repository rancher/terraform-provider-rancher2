package rancher2

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"
)

func tcloudPublicMachineConfigWireFields() map[string]string {
	return map[string]string{
		"access_key": "accessKey", "auth_url": "authUrl", "availability_zone": "availabilityZone",
		"bandwidth_size": "bandwidthSize", "bandwidth_type": "bandwidthType",
		"cacert": "cacert", "cloud": "cloud", "domain_id": "domainId", "domain_name": "domainName",
		"eip": "eip", "eip_type": "eipType", "endpoint_type": "endpointType",
		"flavor_id": "flavorId", "flavor_name": "flavorName",
		"image_id": "imageId", "image_name": "imageName", "ip_version": "ipVersion",
		"keypair_name": "keypairName", "network_scope": "networkScope",
		"password": "password", "private_key_file": "privateKeyFile",
		"project_id": "projectId", "project_name": "projectName", "region": "region",
		"root_volume_size": "rootVolumeSize", "root_volume_type": "rootVolumeType",
		"sec_groups": "secGroups", "secret_key": "secretKey",
		"server_group": "serverGroup", "server_group_id": "serverGroupId",
		"skip_default_sg": "skipDefaultSg", "skip_eip": "skipEip",
		"ssh_allow_cidr": "sshAllowCidr", "ssh_port": "sshPort", "ssh_user": "sshUser",
		"subnet_id": "subnetId", "subnet_name": "subnetName",
		"tags": "tags", "token": "token", "user_data_file": "userDataFile",
		"user_data_raw": "userDataRaw", "username": "username", "vpc_id": "vpcId", "vpc_name": "vpcName",
	}
}

func TestTCloudPublicMachineConfigRoundTrip(t *testing.T) {
	wireFields := tcloudPublicMachineConfigWireFields()
	require.Len(t, machineConfigV2OpentelekomcloudFields(), len(wireFields))
	values := make(map[string]interface{}, len(wireFields))
	for key := range wireFields {
		values[key] = "test-" + key
	}
	values["skip_default_sg"], values["skip_eip"] = true, true
	values["network_scope"], values["ip_version"] = "shared", "4"
	values["bandwidth_size"], values["root_volume_size"], values["ssh_port"] = "120", "80", "2222"
	values["sec_groups"], values["tags"] = "group-a,group-b", "environment.test,team.platform"
	config := map[string]interface{}{
		"generate_name": "tcloud-pool", "fleet_namespace": "fleet-test",
		"tcloud_public_config": []interface{}{values},
		"labels":               map[string]interface{}{"environment": "test"},
		"annotations":          map[string]interface{}{"ui.rancher/provider": "opentelekomcloud"},
	}
	d := schema.TestResourceDataRaw(t, machineConfigV2Fields(), config)
	d.SetId("fleet-test:nc-tcloud-pool-abc")
	require.NoError(t, d.Set("name", "nc-tcloud-pool-abc"))
	require.NoError(t, d.Set("resource_version", "7"))
	obj := expandMachineConfigV2(d)
	require.Equal(t, "OpentelekomcloudConfig", obj.Kind)
	require.Equal(t, "rke-machine-config.cattle.io/v1", obj.APIVersion)
	require.Equal(t, obj.TypeMeta, obj.OpentelekomcloudConfig.TypeMeta)
	require.Equal(t, obj.ObjectMeta, obj.OpentelekomcloudConfig.ObjectMeta)
	require.Equal(t, "nc-tcloud-pool-", obj.OpentelekomcloudConfig.GenerateName)
	require.Equal(t, d.Id(), obj.OpentelekomcloudConfig.ID)

	encoded, err := json.Marshal(obj.OpentelekomcloudConfig)
	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(encoded, &body))
	for key, wireKey := range wireFields {
		require.Equal(t, values[key], body[wireKey], wireKey)
		delete(body, wireKey)
	}
	require.Equal(t, "OpentelekomcloudConfig", body["kind"])
	require.Equal(t, "rke-machine-config.cattle.io/v1", body["apiVersion"])
	delete(body, "kind")
	delete(body, "apiVersion")
	delete(body, "id")
	delete(body, "metadata")
	require.Nil(t, body["links"])
	require.Nil(t, body["actions"])
	delete(body, "links")
	delete(body, "actions")
	require.Empty(t, body, "driver options must be top-level fields with no wrapper or spec")

	decoded := &MachineConfigV2Opentelekomcloud{}
	require.NoError(t, json.Unmarshal(encoded, decoded))
	require.Equal(t, obj.OpentelekomcloudConfig, decoded)
	require.Equal(t, []interface{}{values}, flattenMachineConfigV2Opentelekomcloud(decoded))
	target := schema.TestResourceDataRaw(t, machineConfigV2Fields(), map[string]interface{}{"generate_name": "tcloud-pool"})
	require.NoError(t, flattenMachineConfigV2(target, obj))
	for _, field := range []string{"name", "fleet_namespace", "resource_version", "annotations", "labels", "tcloud_public_config"} {
		require.Equal(t, d.Get(field), target.Get(field), field)
	}
	require.Equal(t, d.Id(), target.Id())
	require.Equal(t, "OpentelekomcloudConfig", target.Get("kind"))
	diff, err := resourceRancher2MachineConfigV2().Diff(target.State(), terraform.NewResourceConfigRaw(config), nil)
	require.NoError(t, err)
	if diff != nil {
		require.True(t, diff.Empty(), "refresh should not introduce a new plan: %v", diff)
	}
}

func TestTCloudPublicMachineConfigEmptyValues(t *testing.T) {
	require.Nil(t, flattenMachineConfigV2Opentelekomcloud(nil))
	require.Nil(t, expandMachineConfigV2Opentelekomcloud(nil, &MachineConfigV2{}))
	require.Nil(t, expandMachineConfigV2Opentelekomcloud([]interface{}{nil}, &MachineConfigV2{}))
	obj := expandMachineConfigV2Opentelekomcloud([]interface{}{map[string]interface{}{
		"bandwidth_size": "0", "root_volume_size": "0", "ssh_port": "0",
	}}, &MachineConfigV2{})
	encoded, err := json.Marshal(obj)
	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(encoded, &body))
	flattened := flattenMachineConfigV2Opentelekomcloud(obj)[0].(map[string]interface{})
	for key, wireKey := range tcloudPublicMachineConfigWireFields() {
		var expected interface{} = ""
		switch key {
		case "skip_default_sg", "skip_eip":
			expected = false
		case "bandwidth_size", "root_volume_size", "ssh_port":
			expected = "0"
		}
		require.Contains(t, body, wireKey)
		require.Equal(t, expected, body[wireKey], wireKey)
		require.Equal(t, expected, flattened[key], key)
	}
}
