package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterEKSImportConf      *managementClient.EKSClusterConfigSpec
	testClusterEKSImportInterface []interface{}
)

func init() {
	testClusterEKSImportConf = &managementClient.EKSClusterConfigSpec{
		AmazonCredentialSecret: "test",
		DisplayName:            "eksimport",
		Region:                 "test",
		Imported:               true,
	}
	testClusterEKSImportInterface = []interface{}{
		map[string]interface{}{
			"cloud_credential": "test",
			"name":             "eksimport",
			"region":           "test",
		},
	}
}

func TestFlattenClusterEKSImport(t *testing.T) {

	cases := []struct {
		Input          *managementClient.EKSClusterConfigSpec
		ExpectedOutput []interface{}
	}{
		{
			testClusterEKSImportConf,
			testClusterEKSImportInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterEKSImport(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterEKSImport(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.EKSClusterConfigSpec
	}{
		{
			testClusterEKSImportInterface,
			testClusterEKSImportConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterEKSImport(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
