package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"
)

func TestTCloudPublicClusterV2MachinePoolReferences(t *testing.T) {
	for _, workers := range []int{0, 1, 2} {
		t.Run(fmt.Sprintf("%d-nodes", workers+1), func(t *testing.T) {
			machineData := schema.TestResourceDataRaw(t, machineConfigV2Fields(), map[string]interface{}{
				"generate_name": "tcloud-pool", "fleet_namespace": "fleet-test",
				"tcloud_public_config": []interface{}{map[string]interface{}{
					"network_scope": "shared", "skip_default_sg": true,
					"vpc_id": "prepared-vpc", "subnet_id": "prepared-subnet", "sec_groups": "prepared-group",
				}},
			})
			machine := expandMachineConfigV2(machineData)
			machine.ID = "fleet-test:nc-tcloud-pool-example"
			machine.Name = "nc-tcloud-pool-example"
			require.NoError(t, flattenMachineConfigV2(machineData, machine))
			require.Equal(t, "shared", machine.OpentelekomcloudConfig.NetworkScope)
			require.True(t, machine.OpentelekomcloudConfig.SkipDefaultSG)
			require.Equal(t, "prepared-group", machine.OpentelekomcloudConfig.SecGroups)

			reference := []interface{}{map[string]interface{}{
				"kind": machineData.Get("kind"), "name": machineData.Get("name"),
				"api_version": "rke-machine-config.cattle.io/v1",
			}}
			annotations := map[string]interface{}{
				"ui.rancher/provider":                              "opentelekomcloud",
				"infrastructure.otc.t-systems.com/cluster-network": "tcloud-network",
				"infrastructure.otc.t-systems.com/network-policy":  "Observe",
			}
			values := map[string]interface{}{
				"name": "tcloud-cluster", "fleet_namespace": machineData.Get("fleet_namespace"),
				"kubernetes_version": "v1.34.1+rke2r1", "annotations": annotations,
				"rke_config": []interface{}{map[string]interface{}{
					"machine_pools": []interface{}{
						map[string]interface{}{
							"name": "server", "quantity": 1,
							"control_plane_role": true, "etcd_role": true, "worker_role": true,
							"cloud_credential_secret_name": "cattle-global-data:cc-server",
							"machine_config":               reference,
						},
						map[string]interface{}{
							"name": "workers", "quantity": workers,
							"control_plane_role": false, "etcd_role": false, "worker_role": true,
							"cloud_credential_secret_name": "cattle-global-data:cc-workers",
							"machine_config":               reference,
						},
					},
				}},
			}
			_, errs := resourceRancher2ClusterV2().Validate(terraform.NewResourceConfigRaw(values))
			require.Empty(t, errs)
			d := schema.TestResourceDataRaw(t, clusterV2Fields(), values)
			cluster, err := expandClusterV2(d)
			require.NoError(t, err)
			require.Equal(t, machine.Namespace, cluster.Namespace)
			require.Equal(t, toMapString(annotations), cluster.Annotations)
			require.NotNil(t, cluster.Spec.RKEConfig)
			pools := cluster.Spec.RKEConfig.MachinePools
			require.Len(t, pools, 2)
			for _, pool := range pools {
				require.NotNil(t, pool.NodeConfig)
				require.Equal(t, "OpentelekomcloudConfig", pool.NodeConfig.Kind)
				require.Equal(t, machine.Name, pool.NodeConfig.Name)
				require.Equal(t, machine.APIVersion, pool.NodeConfig.APIVersion)
			}
			require.Equal(t, "cattle-global-data:cc-server", pools[0].CloudCredentialSecretName)
			require.True(t, pools[0].ControlPlaneRole)
			require.True(t, pools[0].EtcdRole)
			require.True(t, pools[0].WorkerRole)
			require.Equal(t, int32(1), *pools[0].Quantity)
			require.Equal(t, "cattle-global-data:cc-workers", pools[1].CloudCredentialSecretName)
			require.False(t, pools[1].ControlPlaneRole)
			require.False(t, pools[1].EtcdRole)
			require.True(t, pools[1].WorkerRole)
			require.Equal(t, int32(workers), *pools[1].Quantity)
			require.Equal(t, pools, expandClusterV2RKEConfigMachinePools(flattenClusterV2RKEConfigMachinePools(pools)))
			require.NoError(t, flattenClusterV2(d, cluster))
			require.Equal(t, annotations, d.Get("annotations"))
		})
	}
}
