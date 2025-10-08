package rancher2

import (
	"testing"

	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterV2RKEConfigConf      *provisionv1.RKEConfig
	testClusterV2RKEConfigInterface []interface{}
)

func init() {
	testClusterV2RKEConfigConf = &provisionv1.RKEConfig{
		MachinePools: testClusterV2RKEConfigMachinePoolsConf,
	}
	testClusterV2RKEConfigConf.AdditionalManifest = "additional_manifest"
	testClusterV2RKEConfigConf.DataDirectories = testClusterV2RKEConfigDataDirectoriesConf
	testClusterV2RKEConfigConf.UpgradeStrategy = testClusterV2RKEConfigUpgradeStrategyConf
	testClusterV2RKEConfigConf.ChartValues = rkev1.GenericMap{
		Data: map[string]interface{}{
			"chart_one": "one",
			"chart_two": "two",
		},
	}
	testClusterV2RKEConfigConf.MachinePoolDefaults = testClusterV2RKEConfigMachinePoolDefaultsConf
	testClusterV2RKEConfigConf.MachineGlobalConfig = rkev1.GenericMap{
		Data: map[string]interface{}{
			"config_one": "one",
			"config_two": "two",
		},
	}
	testClusterV2RKEConfigConf.MachineSelectorConfig = testClusterV2RKEConfigSystemConfigConf
	testClusterV2RKEConfigConf.MachineSelectorFiles = testClusterV2RKEConfigMachineSelectorFilesConf
	testClusterV2RKEConfigConf.Registries = testClusterV2RKEConfigRegistryConf
	testClusterV2RKEConfigConf.ETCD = testClusterV2RKEConfigETCDConf
	testClusterV2RKEConfigConf.RotateCertificates = testClusterV2RKEConfigRotateCertificatesConf
	testClusterV2RKEConfigConf.ETCDSnapshotCreate = testClusterV2RKEConfigETCDSnapshotCreateConf
	testClusterV2RKEConfigConf.ETCDSnapshotRestore = testClusterV2RKEConfigETCDSnapshotRestoreConf

	testClusterV2RKEConfigInterface = []interface{}{
		map[string]interface{}{
			"additional_manifest":     "additional_manifest",
			"data_directories":        testClusterV2RKEConfigDataDirectoriesInterface,
			"upgrade_strategy":        testClusterV2RKEConfigUpgradeStrategyInterface,
			"chart_values":            "chart_one: one\nchart_two: two\n",
			"machine_global_config":   "config_one: one\nconfig_two: two\n",
			"machine_pools":           testClusterV2RKEConfigMachinePoolsInterface,
			"machine_pool_defaults":   testClusterV2RKEConfigMachinePoolDefaultsInterface,
			"machine_selector_config": testClusterV2RKEConfigSystemConfigInterface,
			"machine_selector_files":  testClusterV2RKEConfigMachineSelectorFilesInterface,
			"registries":              testClusterV2RKEConfigRegistryInterface,
			"etcd":                    testClusterV2RKEConfigETCDInterface,
			"rotate_certificates":     testClusterV2RKEConfigRotateCertificatesInterface,
			"etcd_snapshot_create":    testClusterV2RKEConfigETCDSnapshotCreateInterface,
			"etcd_snapshot_restore":   testClusterV2RKEConfigETCDSnapshotRestoreInterface,
		},
	}
}

func TestFlattenClusterV2RKEConfig(t *testing.T) {

	cases := []struct {
		Input          *provisionv1.RKEConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigConf,
			testClusterV2RKEConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfig(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterV2RKEConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *provisionv1.RKEConfig
	}{
		{
			testClusterV2RKEConfigInterface,
			testClusterV2RKEConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfig(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
