package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testClusterV2RKEConfigSystemConfigLabelSelectorExpressionConf      []metav1.LabelSelectorRequirement
	testClusterV2RKEConfigSystemConfigLabelSelectorExpressionInterface []interface{}
	testClusterV2RKEConfigSystemConfigLabelSelectorConf                *metav1.LabelSelector
	testClusterV2RKEConfigSystemConfigLabelSelectorInterface           []interface{}
	testClusterV2RKEConfigSystemConfigConf                             []rkev1.RKESystemConfig
	testClusterV2RKEConfigSystemConfigInterface                        []interface{}
	testClusterV2RKEConfigMachineSelectorFilesConf                     []rkev1.RKEProvisioningFiles
	testClusterV2RKEConfigMachineSelectorFilesInterface                []interface{}
)

func init() {
	testClusterV2RKEConfigSystemConfigLabelSelectorExpressionConf = []metav1.LabelSelectorRequirement{
		{
			Key:      "key",
			Operator: "operator",
			Values:   []string{"value1", "value2"},
		},
	}

	testClusterV2RKEConfigSystemConfigLabelSelectorExpressionInterface = []interface{}{
		map[string]interface{}{
			"key":      "key",
			"operator": "operator",
			"values":   []interface{}{"value1", "value2"},
		},
	}
	testClusterV2RKEConfigSystemConfigLabelSelectorConf = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"label_one": "one",
			"label_two": "two",
		},
		MatchExpressions: testClusterV2RKEConfigSystemConfigLabelSelectorExpressionConf,
	}
	testClusterV2RKEConfigSystemConfigLabelSelectorInterface = []interface{}{
		map[string]interface{}{
			"match_labels": map[string]interface{}{
				"label_one": "one",
				"label_two": "two",
			},
			"match_expressions": testClusterV2RKEConfigSystemConfigLabelSelectorExpressionInterface,
		},
	}
	testClusterV2RKEConfigSystemConfigConf = []rkev1.RKESystemConfig{
		{
			MachineLabelSelector: testClusterV2RKEConfigSystemConfigLabelSelectorConf,
			Config: rkev1.GenericMap{
				Data: map[string]interface{}{
					"config_one": "one",
					"config_two": "two",
				},
			},
		},
	}
	testClusterV2RKEConfigSystemConfigInterface = []interface{}{
		map[string]interface{}{
			"machine_label_selector": testClusterV2RKEConfigSystemConfigLabelSelectorInterface,
			"config":                 "config_one: one\nconfig_two: two\n",
		},
	}

	testClusterV2RKEConfigMachineSelectorFilesConf = []rkev1.RKEProvisioningFiles{
		{
			MachineLabelSelector: testClusterV2RKEConfigSystemConfigLabelSelectorConf,
			FileSources: []rkev1.ProvisioningFileSource{
				{
					Secret: rkev1.K8sObjectFileSource{
						Name:               "test-config-secret",
						DefaultPermissions: "0644",
						Items: []rkev1.KeyToPath{
							{
								Key:         "a",
								Path:        "/etc/rancher/rke2/test.yaml",
								Permissions: "600",
								Hash:        "abcdefg",
								Dynamic:     true,
							},
						},
					},
					ConfigMap: rkev1.K8sObjectFileSource{
						Name:               "test-config-configmap",
						DefaultPermissions: "0644",
						Items: []rkev1.KeyToPath{
							{
								Key:         "a",
								Path:        "/etc/rancher/rke2/test.yaml",
								Permissions: "600",
								Hash:        "abcdefg",
								Dynamic:     true,
							},
						},
					},
				},
			},
		},
	}

	testClusterV2RKEConfigMachineSelectorFilesInterface = []interface{}{
		map[string]interface{}{
			"machine_label_selector": testClusterV2RKEConfigSystemConfigLabelSelectorInterface,
			"file_sources": []interface{}{
				map[string]interface{}{
					"secret": []interface{}{
						map[string]interface{}{
							"name":                "test-config-secret",
							"default_permissions": "0644",
							"items": []interface{}{
								map[string]interface{}{
									"key":         "a",
									"path":        "/etc/rancher/rke2/test.yaml",
									"permissions": "600",
									"hash":        "abcdefg",
									"dynamic":     true,
								},
							},
						},
					},
					"configmap": []interface{}{
						map[string]interface{}{
							"name":                "test-config-configmap",
							"default_permissions": "0644",
							"items": []interface{}{
								map[string]interface{}{
									"key":         "a",
									"path":        "/etc/rancher/rke2/test.yaml",
									"permissions": "600",
									"hash":        "abcdefg",
									"dynamic":     true,
								},
							},
						},
					},
				},
			},
		},
	}

}

func TestFlattenClusterV2RKEConfigSystemConfigLabelSelectorExpression(t *testing.T) {

	cases := []struct {
		Input          []metav1.LabelSelectorRequirement
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigSystemConfigLabelSelectorExpressionConf,
			testClusterV2RKEConfigSystemConfigLabelSelectorExpressionInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigSystemConfigLabelSelectorExpression(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenClusterV2RKEConfigSystemConfigLabelSelector(t *testing.T) {

	cases := []struct {
		Input          *metav1.LabelSelector
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigSystemConfigLabelSelectorConf,
			testClusterV2RKEConfigSystemConfigLabelSelectorInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigSystemConfigLabelSelector(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenClusterV2RKEConfigSystemConfig(t *testing.T) {

	cases := []struct {
		Input          []rkev1.RKESystemConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigSystemConfigConf,
			testClusterV2RKEConfigSystemConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigSystemConfig(tc.Input)

		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterV2RKEConfigSystemConfigLabelSelectorExpression(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []metav1.LabelSelectorRequirement
	}{
		{
			testClusterV2RKEConfigSystemConfigLabelSelectorExpressionInterface,
			testClusterV2RKEConfigSystemConfigLabelSelectorExpressionConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigSystemConfigLabelSelectorExpression(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandClusterV2RKEConfigSystemConfigLabelSelector(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *metav1.LabelSelector
	}{
		{
			testClusterV2RKEConfigSystemConfigLabelSelectorInterface,
			testClusterV2RKEConfigSystemConfigLabelSelectorConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigSystemConfigLabelSelector(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandClusterV2RKEConfigSystemConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []rkev1.RKESystemConfig
	}{
		{
			testClusterV2RKEConfigSystemConfigInterface,
			testClusterV2RKEConfigSystemConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigSystemConfig(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
