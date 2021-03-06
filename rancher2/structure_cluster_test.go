package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterAnswersConf                         *managementClient.Answer
	testClusterAnswersInterface                    []interface{}
	testClusterQuestionsConf                       []managementClient.Question
	testClusterQuestionsInterface                  []interface{}
	testLocalClusterAuthEndpointConf               *managementClient.LocalClusterAuthEndpoint
	testLocalClusterAuthEndpointInterface          []interface{}
	testClusterRegistrationTokenConf               *managementClient.ClusterRegistrationToken
	testClusterRegistrationToken2Conf              *managementClient.ClusterRegistrationToken
	testClusterRegistrationTokenInterface          []interface{}
	testClusterGenerateKubeConfigOutput            *managementClient.GenerateKubeConfigOutput
	testClusterScheduledClusterScanConfigConf      *managementClient.ScheduledClusterScanConfig
	testClusterScheduledClusterScanConfigInterface []interface{}
	testClusterScheduledClusterScanConf            *managementClient.ScheduledClusterScan
	testClusterScheduledClusterScanInterface       []interface{}
	testClusterConfAKS                             *Cluster
	testClusterInterfaceAKS                        map[string]interface{}
	testClusterConfEKS                             *Cluster
	testClusterInterfaceEKS                        map[string]interface{}
	testClusterConfEKSV2                           *Cluster
	testClusterInterfaceEKSV2                      map[string]interface{}
	testClusterConfGKE                             *Cluster
	testClusterInterfaceGKE                        map[string]interface{}
	testClusterConfOKE                             *Cluster
	testClusterInterfaceOKE                        map[string]interface{}
	testClusterConfRKE                             *Cluster
	testClusterInterfaceRKE                        map[string]interface{}
	testClusterConfTemplate                        *Cluster
	testClusterInterfaceTemplate                   map[string]interface{}
)

func testCluster() {
	testClusterAnswersConf = &managementClient.Answer{
		ClusterID: "cluster_id",
		ProjectID: "project_id",
		Values: map[string]string{
			"string.value1": "one",
			"string.value2": "two",
		},
	}
	testClusterAnswersInterface = []interface{}{
		map[string]interface{}{
			"cluster_id": "cluster_id",
			"project_id": "project_id",
			"values": map[string]interface{}{
				"string.value1": "one",
				"string.value2": "two",
			},
		},
	}
	testClusterQuestionsConf = []managementClient.Question{
		{
			Default:  "default",
			Required: true,
			Type:     "string",
			Variable: "string.value1",
		},
		{
			Default:  "default",
			Required: true,
			Type:     "string",
			Variable: "string.value2",
		},
		{
			Default:  "default",
			Required: true,
			Type:     "password",
			Variable: "password.var",
		},
	}
	testClusterQuestionsInterface = []interface{}{
		map[string]interface{}{
			"default":  "default",
			"required": true,
			"type":     "string",
			"variable": "string.value1",
		},
		map[string]interface{}{
			"default":  "default",
			"required": true,
			"type":     "string",
			"variable": "string.value2",
		},
		map[string]interface{}{
			"default":  "default",
			"required": true,
			"type":     "password",
			"variable": "password.var",
		},
	}
	testLocalClusterAuthEndpointConf = &managementClient.LocalClusterAuthEndpoint{
		CACerts: "cacerts",
		Enabled: true,
		FQDN:    "fqdn",
	}
	testLocalClusterAuthEndpointInterface = []interface{}{
		map[string]interface{}{
			"ca_certs": "cacerts",
			"enabled":  true,
			"fqdn":     "fqdn",
		},
	}
	testClusterRegistrationTokenConf = &managementClient.ClusterRegistrationToken{
		ClusterID:          "cluster_test",
		Name:               clusterRegistrationTokenName,
		Command:            "command",
		InsecureCommand:    "insecure_command",
		ManifestURL:        "manifest",
		NodeCommand:        "node_command",
		Token:              "token",
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
			"token":                "token",
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
	testClusterScheduledClusterScanConfigConf = &managementClient.ScheduledClusterScanConfig{
		CronSchedule: "cron_schedule",
		Retention:    5,
	}
	testClusterScheduledClusterScanConfigInterface = []interface{}{
		map[string]interface{}{
			"cron_schedule": "cron_schedule",
			"retention":     5,
		},
	}
	testClusterScheduledClusterScanConf = &managementClient.ScheduledClusterScan{
		Enabled:        true,
		ScanConfig:     testClusterScanConfigConf,
		ScheduleConfig: testClusterScheduledClusterScanConfigConf,
	}
	testClusterScheduledClusterScanInterface = []interface{}{
		map[string]interface{}{
			"enabled":         true,
			"scan_config":     testClusterScanConfigInterface,
			"schedule_config": testClusterScheduledClusterScanConfigInterface,
		},
	}
	testClusterConfAKS = &Cluster{
		AzureKubernetesServiceConfig: testClusterAKSConfigConf,
	}
	testClusterConfAKS.Name = "test"
	testClusterConfAKS.Description = "description"
	testClusterConfAKS.Driver = clusterDriverAKS
	testClusterConfAKS.DefaultPodSecurityPolicyTemplateID = "restricted"
	testClusterConfAKS.EnableClusterMonitoring = true
	testClusterConfAKS.EnableNetworkPolicy = newTrue()
	testClusterConfAKS.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceAKS = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_policy_template_id": "restricted",
		"enable_cluster_monitoring":               true,
		"enable_network_policy":                   true,
		"kube_config":                             "kube_config",
		"driver":                                  clusterDriverAKS,
		"aks_config":                              testClusterAKSConfigInterface,
		"system_project_id":                       "system_project_id",
	}
	testClusterConfEKS = &Cluster{
		AmazonElasticContainerServiceConfig: testClusterEKSConfigConf,
	}
	testClusterConfEKS.Name = "test"
	testClusterConfEKS.Description = "description"
	testClusterConfEKS.Driver = clusterDriverEKS
	testClusterConfEKS.DefaultPodSecurityPolicyTemplateID = "restricted"
	testClusterConfEKS.EnableClusterMonitoring = true
	testClusterConfEKS.EnableNetworkPolicy = newTrue()
	testClusterConfEKS.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceEKS = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_policy_template_id": "restricted",
		"enable_cluster_monitoring":               true,
		"enable_network_policy":                   true,
		"kube_config":                             "kube_config",
		"driver":                                  clusterDriverEKS,
		"eks_config":                              testClusterEKSConfigInterface,
		"system_project_id":                       "system_project_id",
	}
	testClusterConfEKSV2 = &Cluster{}
	testClusterConfEKSV2.EKSConfig = testClusterEKSConfigV2Conf
	testClusterConfEKSV2.Name = "test"
	testClusterConfEKSV2.Description = "description"
	testClusterConfEKSV2.Driver = clusterDriverEKSV2
	testClusterConfEKSV2.DefaultPodSecurityPolicyTemplateID = "restricted"
	testClusterConfEKSV2.EnableClusterMonitoring = true
	testClusterConfEKSV2.EnableNetworkPolicy = newTrue()
	testClusterConfEKSV2.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceEKSV2 = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_policy_template_id": "restricted",
		"enable_cluster_monitoring":               true,
		"enable_network_policy":                   true,
		"kube_config":                             "kube_config",
		"driver":                                  clusterDriverEKSV2,
		"eks_config_v2":                           testClusterEKSConfigV2Interface,
		"system_project_id":                       "system_project_id",
	}
	testClusterConfGKE = &Cluster{
		GoogleKubernetesEngineConfig: testClusterGKEConfigConf,
	}
	testClusterConfGKE.Name = "test"
	testClusterConfGKE.Description = "description"
	testClusterConfGKE.Driver = clusterDriverGKE
	testClusterConfGKE.DefaultPodSecurityPolicyTemplateID = "restricted"
	testClusterConfGKE.EnableClusterMonitoring = true
	testClusterConfGKE.EnableNetworkPolicy = newTrue()
	testClusterConfGKE.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceGKE = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_policy_template_id": "restricted",
		"enable_cluster_monitoring":               true,
		"enable_network_policy":                   true,
		"kube_config":                             "kube_config",
		"driver":                                  clusterDriverGKE,
		"gke_config":                              testClusterGKEConfigInterface,
		"system_project_id":                       "system_project_id",
	}
	testClusterConfOKE = &Cluster{
		OracleKubernetesEngineConfig: testClusterOKEConfigConf,
	}
	testClusterConfOKE.Name = "test"
	testClusterConfOKE.Description = "description"
	testClusterConfOKE.Driver = clusterOKEKind
	testClusterConfOKE.DefaultPodSecurityPolicyTemplateID = "restricted"
	testClusterConfOKE.EnableClusterMonitoring = true
	testClusterConfOKE.EnableNetworkPolicy = newTrue()
	testClusterConfOKE.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceOKE = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_policy_template_id": "restricted",
		"enable_cluster_monitoring":               true,
		"enable_network_policy":                   true,
		"kube_config":                             "kube_config",
		"driver":                                  clusterOKEKind,
		"oke_config":                              testClusterOKEConfigInterface,
		"system_project_id":                       "system_project_id",
	}
	testClusterConfRKE = &Cluster{}
	testClusterConfRKE.Name = "test"
	testClusterConfRKE.Description = "description"
	testClusterConfRKE.RancherKubernetesEngineConfig = testClusterRKEConfigConf
	testClusterConfRKE.Driver = clusterDriverRKE
	testClusterConfRKE.DefaultPodSecurityPolicyTemplateID = "restricted"
	testClusterConfRKE.EnableClusterMonitoring = true
	testClusterConfRKE.EnableNetworkPolicy = newTrue()
	testClusterConfRKE.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterConfRKE.ScheduledClusterScan = testClusterScheduledClusterScanConf
	testClusterInterfaceRKE = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_policy_template_id": "restricted",
		"enable_cluster_monitoring":               true,
		"enable_network_policy":                   true,
		"kube_config":                             "kube_config",
		"driver":                                  clusterDriverRKE,
		"rke_config":                              testClusterRKEConfigInterface,
		"scheduled_cluster_scan":                  testClusterScheduledClusterScanInterface,
		"system_project_id":                       "system_project_id",
		"windows_prefered_cluster":                false,
	}
	testClusterConfTemplate = &Cluster{}
	testClusterConfTemplate.Name = "test"
	testClusterConfTemplate.Description = "description"
	testClusterConfTemplate.ClusterTemplateAnswers = testClusterAnswersConf
	testClusterConfTemplate.ClusterTemplateID = "cluster_template_id"
	testClusterConfTemplate.ClusterTemplateQuestions = testClusterQuestionsConf
	testClusterConfTemplate.ClusterTemplateRevisionID = "cluster_template_revision_id"
	testClusterConfTemplate.Driver = clusterDriverRKE
	testClusterConfTemplate.DefaultPodSecurityPolicyTemplateID = "restricted"
	testClusterConfTemplate.EnableClusterAlerting = true
	testClusterConfTemplate.EnableClusterMonitoring = true
	testClusterConfTemplate.EnableNetworkPolicy = newTrue()
	testClusterConfTemplate.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterConfTemplate.ScheduledClusterScan = testClusterScheduledClusterScanConf
	testClusterInterfaceTemplate = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_policy_template_id": "restricted",
		"enable_cluster_alerting":                 true,
		"enable_cluster_monitoring":               true,
		"enable_network_policy":                   true,
		"kube_config":                             "kube_config",
		"driver":                                  clusterDriverRKE,
		"cluster_template_answers":                testClusterAnswersInterface,
		"cluster_template_id":                     "cluster_template_id",
		"cluster_template_questions":              testClusterQuestionsInterface,
		"cluster_template_revision_id":            "cluster_template_revision_id",
		"rke_config":                              []interface{}{},
		"scheduled_cluster_scan":                  testClusterScheduledClusterScanInterface,
		"system_project_id":                       "system_project_id",
		"windows_prefered_cluster":                false,
	}
}

func TestFlattenClusterRegistationToken(t *testing.T) {
	testCluster()
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
			testClusterConfEKSV2,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceEKSV2,
		},
		{
			testClusterConfGKE,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceGKE,
		},
		{
			testClusterConfOKE,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceOKE,
		},
		{
			testClusterConfRKE,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceRKE,
		},
		{
			testClusterConfTemplate,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceTemplate,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterFields(), map[string]interface{}{})
		tc.InputToken.ID = "id"
		err := flattenCluster(output, tc.Input, tc.InputToken, tc.InputKube, tc.ExpectedOutput["default_project_id"].(string), tc.ExpectedOutput["system_project_id"].(string), nil)
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
			testClusterInterfaceEKSV2,
			testClusterConfEKSV2,
		},
		{
			testClusterInterfaceGKE,
			testClusterConfGKE,
		},
		{
			testClusterInterfaceOKE,
			testClusterConfOKE,
		},
		{
			testClusterInterfaceRKE,
			testClusterConfRKE,
		},
		{
			testClusterInterfaceTemplate,
			testClusterConfTemplate,
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

func TestFlattenClusterWithPreservedClusterTemplateAnswers(t *testing.T) {

	testClusterInterfaceTemplate["cluster_template_answers"] = []interface{}{
		map[string]interface{}{
			"cluster_id": "cluster_id",
			"project_id": "project_id",
			"values": map[string]interface{}{
				"string.value1": "one",
				"string.value2": "two",
				"password.var":  "password",
			},
		},
	}

	cases := []struct {
		Input          *Cluster
		InputToken     *managementClient.ClusterRegistrationToken
		InputKube      *managementClient.GenerateKubeConfigOutput
		ExpectedOutput map[string]interface{}
	}{
		{

			testClusterConfTemplate,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceTemplate,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterFields(), map[string]interface{}{
			"cluster_template_answers": []interface{}{
				map[string]interface{}{
					"cluster_id": "cluster_id",
					"project_id": "project_id",
					"values": map[string]interface{}{
						"password.var": "password",
					},
				},
			},
		})
		tc.InputToken.ID = "id"
		err := flattenCluster(output, tc.Input, tc.InputToken, tc.InputKube, tc.ExpectedOutput["default_project_id"].(string), tc.ExpectedOutput["system_project_id"].(string), nil)
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

func TestReadPreservedClusterTemplateAnswers(t *testing.T) {

	inputResourceData := schema.TestResourceDataRaw(t, clusterFields(), map[string]interface{}{
		"cluster_template_answers": []interface{}{
			map[string]interface{}{
				"cluster_id": "cluster_id",
				"project_id": "project_id",
				"values": map[string]interface{}{
					"password.var":  "password",
					"string.value1": "one",
				},
			},
		},
		"cluster_template_questions": []interface{}{
			map[string]interface{}{
				"default":  "default",
				"required": true,
				"type":     "string",
				"variable": "string.value1",
			},
			map[string]interface{}{
				"default":  "default",
				"required": true,
				"type":     "password",
				"variable": "password.var",
			},
		},
	})

	expectedOutput := map[string]string{
		"password.var": "password",
	}

	result := readPreservedClusterTemplateAnswers(inputResourceData)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Fatalf("Unexpected result from preserved answers.\nExpected: %#v\nGiven:    %#v",
			expectedOutput, result)
	}
}

func TestFlattenClusterNodes(t *testing.T) {

	testClusterNodes := []managementClient.Node{
		{
			Resource: types.Resource{
				ID: "id",
			},
			Annotations: map[string]string{
				"node_one": "one",
				"node_two": "two",
			},
			ClusterID: "cluster_id",
			Capacity: map[string]string{
				"cpu":    "4",
				"memory": "8156056Ki",
				"pods":   "110",
			},
			ControlPlane:      true,
			Etcd:              true,
			ExternalIPAddress: "172.18.0.5",
			Hostname:          "hostname",
			IPAddress:         "172.18.0.5",
			Info: &managementClient.NodeInfo{
				CPU: &managementClient.CPUInfo{
					Count: 4,
				},
				Kubernetes: &managementClient.KubernetesInfo{
					KubeProxyVersion: "v1.19.7",
					KubeletVersion:   "v1.19.7",
				},
				Memory: &managementClient.MemoryInfo{
					MemTotalKiB: 8156056,
				},
				OS: &managementClient.OSInfo{
					DockerVersion:   "containerd://1.4.3",
					KernelVersion:   "4.19.121",
					OperatingSystem: "Unknown",
				},
			},
			Labels: map[string]string{
				"option1": "value1",
				"option2": "value2",
			},
			Name:              "name",
			NodeName:          "node_name",
			NodePoolID:        "node_pool_id",
			NodeTemplateID:    "node_template_id",
			ProviderId:        "provider_id",
			RequestedHostname: "requested_hostname",
			SshUser:           "ssh_user",
			Worker:            true,
		},
	}

	testClusterNodesInterface := []interface{}{
		map[string]interface{}{
			"id": "id",
			"annotations": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"capacity": map[string]interface{}{
				"cpu":    "4",
				"memory": "8156056Ki",
				"pods":   "110",
			},
			"cluster_id":          "cluster_id",
			"external_ip_address": "172.18.0.5",
			"hostname":            "hostname",
			"ip_address":          "172.18.0.5",
			"system_info": map[string]string{
				"kube_proxy_version":        "v1.19.7",
				"kubelet_version":           "v1.19.7",
				"container_runtime_version": "containerd://1.4.3",
				"kernel_version":            "4.19.121",
				"operating_system":          "Unknown",
			},
			"labels": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"name":               "node_name",
			"node_pool_id":       "node_pool_id",
			"node_template_id":   "node_template_id",
			"provider_id":        "provider_id",
			"roles":              []string{"control_plane", "etcd", "worker"},
			"requested_hostname": "requested_hostname",
			"ssh_user":           "ssh_user",
		},
	}

	cases := []struct {
		Input          []managementClient.Node
		ExpectedOutput []interface{}
	}{
		{
			Input:          testClusterNodes,
			ExpectedOutput: testClusterNodesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterNodes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}

}
