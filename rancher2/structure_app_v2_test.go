package rancher2

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rancher/pkg/api/steve/catalog/types"
	v1 "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

var (
	testAppV2Conf                   *AppV2
	testAppV2ChartInfo              *types.ChartInfo
	testAppV2ChartInstallConf       []types.ChartInstall
	testAppV2ChartInstallActionConf *types.ChartInstallAction
	testAppV2ChartUpgradeConf       []types.ChartUpgrade
	testAppV2ChartUpgradeActionConf *types.ChartUpgradeAction
	testAppV2InstallInterface       map[string]interface{}
	testAppV2Interface              map[string]interface{}
)

func init() {
	testAppV2Conf = &AppV2{}

	testAppV2ChartInfo = &types.ChartInfo{
		Chart: v3.MapStringInterface{
			"name":    "chart_name",
			"version": "chart_version",
		},
	}

	testAppV2Conf.TypeMeta.Kind = appV2Kind
	testAppV2Conf.TypeMeta.APIVersion = appV2APIGroup + "/" + appV2APIVersion

	testAppV2Conf.ObjectMeta.Name = "name"
	testAppV2Conf.ObjectMeta.Namespace = "namespace"
	testAppV2Conf.ObjectMeta.Annotations = map[string]string{
		"annotation1": "one",
		"annotation2": "two",
	}
	testAppV2Conf.ObjectMeta.Labels = map[string]string{
		"label1": "one",
		"label2": "two",
	}
	testAppV2Conf.Spec.Name = "name"
	testAppV2Conf.Spec.Namespace = "namespace"
	testAppV2Conf.Spec.Chart = &v1.Chart{
		Metadata: &v1.Metadata{
			Name:    "chart_name",
			Version: "chart_version",
		},
	}
	testAppV2Conf.Spec.Values = v3.MapStringInterface{
		"value1": "one",
		"value2": "two",
		"global": map[string]interface{}{
			"systemDefaultRegistry": "",
			"cattle": map[string]interface{}{
				"clusterId":             "cluster_id",
				"clusterName":           "",
				"systemDefaultRegistry": "",
			},
		},
	}

	testAppV2ChartInstallConf = []types.ChartInstall{
		{
			ChartName:   "chart_name",
			Version:     "chart_version",
			ReleaseName: "name",
			Values:      testAppV2Conf.Spec.Values,
			Annotations: map[string]string{
				"annotation1": "one",
				"annotation2": "two",
			},
		},
	}

	testAppV2ChartInstallActionConf = &types.ChartInstallAction{
		Timeout:                  nil,
		Wait:                     true,
		DisableHooks:             true,
		DisableOpenAPIValidation: true,
		Namespace:                "namespace",
		Charts:                   testAppV2ChartInstallConf,
		ProjectID:                "project_id",
	}

	testAppV2ChartUpgradeConf = []types.ChartUpgrade{
		{
			ChartName:   "chart_name",
			Version:     "chart_version",
			ReleaseName: "name",
			Force:       true,
			Values:      testAppV2Conf.Spec.Values,
			Annotations: map[string]string{
				"annotation1": "one",
				"annotation2": "two",
			},
		},
	}

	testAppV2ChartUpgradeActionConf = &types.ChartUpgradeAction{
		Timeout:                  nil,
		Wait:                     true,
		DisableHooks:             true,
		DisableOpenAPIValidation: true,
		Force:                    true,
		CleanupOnFail:            true,
		Namespace:                "namespace",
		Charts:                   testAppV2ChartUpgradeConf,
	}

	valuesStr, err := interfaceToGhodssyaml(testAppV2Conf.Spec.Values)
	if err != nil {
		log.Fatalf("[ERROR] initializing: %#v", err)
	}

	testAppV2InstallInterface = map[string]interface{}{
		"cluster_id":                  "cluster_id",
		"project_id":                  "project_id",
		"name":                        "name",
		"namespace":                   "namespace",
		"chart_name":                  "chart_name",
		"chart_version":               "chart_version",
		"cleanup_on_fail":             true,
		"disable_hooks":               true,
		"disable_open_api_validation": true,
		"force_upgrade":               true,
		"annotations": map[string]interface{}{
			"annotation1": "one",
			"annotation2": "two",
		},
		"values": valuesStr,
		"wait":   true,
	}

	testAppV2Interface = map[string]interface{}{
		"name":          "name",
		"namespace":     "namespace",
		"chart_name":    "chart_name",
		"chart_version": "chart_version",
		"annotations": map[string]interface{}{
			"annotation1": "one",
			"annotation2": "two",
		},
		"labels": map[string]interface{}{
			"label1": "one",
			"label2": "two",
		},
		"values": valuesStr,
	}
}

func TestFlattenAppV2(t *testing.T) {

	cases := []struct {
		Input          *AppV2
		ExpectedOutput map[string]interface{}
	}{
		{
			testAppV2Conf,
			testAppV2Interface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, appV2Fields(), tc.ExpectedOutput)
		err := flattenAppV2(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandChartInstallV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ChartInfo      *types.ChartInfo
		ExpectedOutput []types.ChartInstall
	}{
		{
			testAppV2InstallInterface,
			testAppV2ChartInfo,
			testAppV2ChartInstallConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, appV2Fields(), tc.Input)
		_, output, err := expandChartInstallV2(inputResourceData, tc.ChartInfo)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandChartInstallActionV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ChartInfo      *types.ChartInfo
		ExpectedOutput *types.ChartInstallAction
	}{
		{
			testAppV2InstallInterface,
			testAppV2ChartInfo,
			testAppV2ChartInstallActionConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, appV2Fields(), tc.Input)
		output, err := expandChartInstallActionV2(inputResourceData, tc.ChartInfo)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		tc.ExpectedOutput.Timeout = output.Timeout
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandChartUpgradeV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ChartInfo      *types.ChartInfo
		ExpectedOutput []types.ChartUpgrade
	}{
		{
			testAppV2InstallInterface,
			testAppV2ChartInfo,
			testAppV2ChartUpgradeConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, appV2Fields(), tc.Input)
		_, output, err := expandChartUpgradeV2(inputResourceData, tc.ChartInfo)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandChartUpgradeActionV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ChartInfo      *types.ChartInfo
		ExpectedOutput *types.ChartUpgradeAction
	}{
		{
			testAppV2InstallInterface,
			testAppV2ChartInfo,
			testAppV2ChartUpgradeActionConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, appV2Fields(), tc.Input)
		output, err := expandChartUpgradeActionV2(inputResourceData, tc.ChartInfo)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		tc.ExpectedOutput.Timeout = output.Timeout
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
