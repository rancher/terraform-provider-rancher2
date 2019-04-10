package rancher2

import (
	"reflect"
	"testing"

	"github.com/rancher/norman/types"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testCloudCredentialGenericConf      *genericCredentialConfig
	testCloudCredentialGenericInterface []interface{}
)

func init() {
	testCloudCredentialGenericConf = &genericCredentialConfig{
		driverID:   "rackspaceDriverID",
		driverName: "rackspace",
		config: map[string]interface{}{
			"apiKey": "apiKey",
		},
	}
	testCloudCredentialGenericInterface = []interface{}{
		map[string]interface{}{
			"driver": "rackspace",
			"config": map[string]interface{}{
				"apiKey": "apiKey",
			},
		},
	}
}

func TestFlattenCloudCredentialGeneric(t *testing.T) {

	cases := []struct {
		Input          *genericCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialGenericConf,
			testCloudCredentialGenericInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialGeneric(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialGeneric(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *genericCredentialConfig
	}{
		{
			testCloudCredentialGenericInterface,
			testCloudCredentialGenericConf,
		},
	}

	for _, tc := range cases {
		output, err := expandCloudCredentialGeneric(tc.Input, &dummyNodeDriverFinder{
			driverID:   tc.ExpectedOutput.driverID,
			driverName: tc.ExpectedOutput.driverName,
		})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

type dummyNodeDriverFinder struct {
	driverID   string
	driverName string
}

func (f *dummyNodeDriverFinder) ByID(id string) (*managementClient.NodeDriver, error) {
	return &managementClient.NodeDriver{
		Resource: types.Resource{
			ID: f.driverID,
		},
		Name: f.driverName,
	}, nil
}
