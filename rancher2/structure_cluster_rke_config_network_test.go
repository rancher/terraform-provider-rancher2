package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRKEConfigNetworkCalicoConf       *managementClient.CalicoNetworkProvider
	testClusterRKEConfigNetworkCalicoInterface  []interface{}
	testClusterRKEConfigNetworkCanalConf        *managementClient.CanalNetworkProvider
	testClusterRKEConfigNetworkCanalInterface   []interface{}
	testClusterRKEConfigNetworkFlannelConf      *managementClient.FlannelNetworkProvider
	testClusterRKEConfigNetworkFlannelInterface []interface{}
	testClusterRKEConfigNetworkWeaveConf        *managementClient.WeaveNetworkProvider
	testClusterRKEConfigNetworkWeaveInterface   []interface{}
	testClusterRKEConfigNetworkConfCalico       *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceCalico  []interface{}
	testClusterRKEConfigNetworkConfCanal        *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceCanal   []interface{}
	testClusterRKEConfigNetworkConfFlannel      *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceFlannel []interface{}
	testClusterRKEConfigNetworkConfWeave        *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceWeave   []interface{}
)

func init() {
	testClusterRKEConfigNetworkCalicoConf = &managementClient.CalicoNetworkProvider{
		CloudProvider: "aws",
	}
	testClusterRKEConfigNetworkCalicoInterface = []interface{}{
		map[string]interface{}{
			"cloud_provider": "aws",
		},
	}
	testClusterRKEConfigNetworkCanalConf = &managementClient.CanalNetworkProvider{
		Iface: "eth0",
	}
	testClusterRKEConfigNetworkCanalInterface = []interface{}{
		map[string]interface{}{
			"iface": "eth0",
		},
	}
	testClusterRKEConfigNetworkFlannelConf = &managementClient.FlannelNetworkProvider{
		Iface: "eth0",
	}
	testClusterRKEConfigNetworkFlannelInterface = []interface{}{
		map[string]interface{}{
			"iface": "eth0",
		},
	}
	testClusterRKEConfigNetworkWeaveConf = &managementClient.WeaveNetworkProvider{
		Password: "password",
	}
	testClusterRKEConfigNetworkWeaveInterface = []interface{}{
		map[string]interface{}{
			"password": "password",
		},
	}
	testClusterRKEConfigNetworkConfCalico = &managementClient.NetworkConfig{
		CalicoNetworkProvider: testClusterRKEConfigNetworkCalicoConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: networkPluginCalicoName,
	}
	testClusterRKEConfigNetworkInterfaceCalico = []interface{}{
		map[string]interface{}{
			"calico_network_provider": testClusterRKEConfigNetworkCalicoInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": networkPluginCalicoName,
		},
	}
	testClusterRKEConfigNetworkConfCanal = &managementClient.NetworkConfig{
		CanalNetworkProvider: testClusterRKEConfigNetworkCanalConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: networkPluginCanalName,
	}
	testClusterRKEConfigNetworkInterfaceCanal = []interface{}{
		map[string]interface{}{
			"canal_network_provider": testClusterRKEConfigNetworkCanalInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": networkPluginCanalName,
		},
	}
	testClusterRKEConfigNetworkConfFlannel = &managementClient.NetworkConfig{
		FlannelNetworkProvider: testClusterRKEConfigNetworkFlannelConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: networkPluginFlannelName,
	}
	testClusterRKEConfigNetworkInterfaceFlannel = []interface{}{
		map[string]interface{}{
			"flannel_network_provider": testClusterRKEConfigNetworkFlannelInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": networkPluginFlannelName,
		},
	}
	testClusterRKEConfigNetworkConfWeave = &managementClient.NetworkConfig{
		WeaveNetworkProvider: testClusterRKEConfigNetworkWeaveConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: networkPluginWeaveName,
	}
	testClusterRKEConfigNetworkInterfaceWeave = []interface{}{
		map[string]interface{}{
			"weave_network_provider": testClusterRKEConfigNetworkWeaveInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": networkPluginWeaveName,
		},
	}
}

func TestFlattenClusterRKEConfigNetworkCalico(t *testing.T) {

	cases := []struct {
		Input          *managementClient.CalicoNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkCalicoConf,
			testClusterRKEConfigNetworkCalicoInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkCalico(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetworkCanal(t *testing.T) {

	cases := []struct {
		Input          *managementClient.CanalNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkCanalConf,
			testClusterRKEConfigNetworkCanalInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkCanal(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetworkFlannel(t *testing.T) {

	cases := []struct {
		Input          *managementClient.FlannelNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkFlannelConf,
			testClusterRKEConfigNetworkFlannelInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkFlannel(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetworkWeave(t *testing.T) {

	cases := []struct {
		Input          *managementClient.WeaveNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkWeaveConf,
			testClusterRKEConfigNetworkWeaveInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkWeave(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetwork(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NetworkConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkConfCalico,
			testClusterRKEConfigNetworkInterfaceCalico,
		},
		{
			testClusterRKEConfigNetworkConfCanal,
			testClusterRKEConfigNetworkInterfaceCanal,
		},
		{
			testClusterRKEConfigNetworkConfFlannel,
			testClusterRKEConfigNetworkInterfaceFlannel,
		},
		{
			testClusterRKEConfigNetworkConfWeave,
			testClusterRKEConfigNetworkInterfaceWeave,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetwork(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkCalico(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.CalicoNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkCalicoInterface,
			testClusterRKEConfigNetworkCalicoConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkCalico(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkCanal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.CanalNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkCanalInterface,
			testClusterRKEConfigNetworkCanalConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkCanal(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkFlannel(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.FlannelNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkFlannelInterface,
			testClusterRKEConfigNetworkFlannelConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkFlannel(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkWeave(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.WeaveNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkWeaveInterface,
			testClusterRKEConfigNetworkWeaveConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkWeave(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetwork(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.NetworkConfig
	}{
		{
			testClusterRKEConfigNetworkInterfaceCalico,
			testClusterRKEConfigNetworkConfCalico,
		},
		{
			testClusterRKEConfigNetworkInterfaceCanal,
			testClusterRKEConfigNetworkConfCanal,
		},
		{
			testClusterRKEConfigNetworkInterfaceFlannel,
			testClusterRKEConfigNetworkConfFlannel,
		},
		{
			testClusterRKEConfigNetworkInterfaceWeave,
			testClusterRKEConfigNetworkConfWeave,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetwork(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
