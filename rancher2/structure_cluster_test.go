package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRegistrationTokenConf      *managementClient.ClusterRegistrationToken
	testClusterRegistrationToken2Conf     *managementClient.ClusterRegistrationToken
	testClusterRegistrationTokenInterface []interface{}
	testClusterGenerateKubeConfigOutput   *managementClient.GenerateKubeConfigOutput
	testClusterConfAKS                    *Cluster
	testClusterInterfaceAKS               map[string]interface{}
	testClusterConfEKS                    *Cluster
	testClusterInterfaceEKS               map[string]interface{}
	testClusterConfGKE                    *Cluster
	testClusterInterfaceGKE               map[string]interface{}
	testClusterConfRKE                    *Cluster
	testClusterInterfaceRKE               map[string]interface{}
)

func init() {
	testClusterRegistrationTokenConf = &managementClient.ClusterRegistrationToken{
		ClusterID:          "cluster_test",
		Name:               clusterRegistrationTokenName,
		Command:            "command",
		InsecureCommand:    "insecure_command",
		ManifestURL:        "manifest",
		NodeCommand:        "node_command",
		WindowsNodeCommand: "win_node_command",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testClusterRegistrationToken2Conf = &managementClient.ClusterRegistrationToken{
		ClusterID: "cluster_test",
		Name:      clusterRegistrationTokenName,
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testClusterRegistrationTokenInterface = []interface{}{
		map[string]interface{}{
			"id":                   "id",
			"cluster_id":           "cluster_test",
			"name":                 clusterRegistrationTokenName,
			"command":              "command",
			"insecure_command":     "insecure_command",
			"manifest_url":         "manifest",
			"node_command":         "node_command",
			"windows_node_command": "win_node_command",
			"annotations": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"labels": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
		},
	}
	testClusterGenerateKubeConfigOutput = &managementClient.GenerateKubeConfigOutput{
		Config: "kube_config",
	}
	testClusterConfAKS = &Cluster{
		AzureKubernetesServiceConfig: testClusterAKSConfigConf,
	}
	testClusterConfAKS.Name = "test"
	testClusterConfAKS.Description = "description"
	testClusterConfAKS.Driver = clusterDriverAKS
	testClusterInterfaceAKS = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"description":                "description",
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"kube_config":                "kube_config",
		"driver":                     clusterDriverAKS,
		"aks_config":                 testClusterAKSConfigInterface,
	}
	testClusterConfEKS = &Cluster{
		AmazonElasticContainerServiceConfig: testClusterEKSConfigConf,
	}
	testClusterConfEKS.Name = "test"
	testClusterConfEKS.Description = "description"
	testClusterConfEKS.Driver = clusterDriverEKS
	testClusterInterfaceEKS = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"description":                "description",
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"kube_config":                "kube_config",
		"driver":                     clusterDriverEKS,
		"eks_config":                 testClusterEKSConfigInterface,
	}
	testClusterConfGKE = &Cluster{
		GoogleKubernetesEngineConfig: testClusterGKEConfigConf,
	}
	testClusterConfGKE.Name = "test"
	testClusterConfGKE.Description = "description"
	testClusterConfGKE.Driver = clusterDriverGKE
	testClusterInterfaceGKE = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"description":                "description",
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"kube_config":                "kube_config",
		"driver":                     clusterDriverGKE,
		"gke_config":                 testClusterGKEConfigInterface,
	}
	testClusterConfRKE = &Cluster{}
	testClusterConfRKE.Name = "test"
	testClusterConfRKE.Description = "description"
	testClusterConfRKE.RancherKubernetesEngineConfig = testClusterRKEConfigConf
	testClusterConfRKE.Driver = clusterDriverRKE
	testClusterInterfaceRKE = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"description":                "description",
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"kube_config":                "kube_config",
		"driver":                     clusterDriverRKE,
		"rke_config":                 testClusterRKEConfigInterface,
	}
}

func TestFlattenClusterRegistationToken(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterRegistrationToken
		ExpectedOutput []interface{}
	}{
		{
			testClusterRegistrationTokenConf,
			testClusterRegistrationTokenInterface,
		},
	}

	for _, tc := range cases {
		tc.Input.ID = "id"
		output, err := flattenClusterRegistationToken(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenCluster(t *testing.T) {

	cases := []struct {
		Input          *Cluster
		InputToken     *managementClient.ClusterRegistrationToken
		InputKube      *managementClient.GenerateKubeConfigOutput
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterConfAKS,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceAKS,
		},
		{
			testClusterConfEKS,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceEKS,
		},
		{
			testClusterConfGKE,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceGKE,
		},
		{
			testClusterConfRKE,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceRKE,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterFields(), map[string]interface{}{})
		tc.InputToken.ID = "id"
		err := flattenCluster(output, tc.Input, tc.InputToken, tc.InputKube)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)

		}
		if tc.ExpectedOutput["driver"] == clusterDriverRKE {
			expectedOutput["rke_config"], _ = flattenClusterRKEConfig(tc.Input.RancherKubernetesEngineConfig, []interface{}{})
		}
		expectedOutput["id"] = "id"
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandClusterRegistationToken(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ClusterRegistrationToken
	}{
		{
			testClusterRegistrationTokenInterface,
			testClusterRegistrationToken2Conf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRegistationToken(tc.Input, tc.ExpectedOutput.ClusterID)
		tc.ExpectedOutput.ID = "id"
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCluster(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *Cluster
	}{
		{
			testClusterInterfaceAKS,
			testClusterConfAKS,
		},
		{
			testClusterInterfaceEKS,
			testClusterConfEKS,
		},
		{
			testClusterInterfaceGKE,
			testClusterConfGKE,
		},
		{
			testClusterInterfaceRKE,
			testClusterConfRKE,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterFields(), tc.Input)
		output, err := expandCluster(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
