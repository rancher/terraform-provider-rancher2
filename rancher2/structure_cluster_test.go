package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterEnvVarsConf                           []managementClient.EnvVar
	testClusterEnvVarsInterface                      []interface{}
	testClusterAgentDeploymentCustomizationConf      *managementClient.AgentDeploymentCustomization
	testClusterAgentDeploymentCustomizationInterface []interface{}
	testFleetAgentDeploymentCustomizationConf        *managementClient.AgentDeploymentCustomization
	testFleetAgentDeploymentCustomizationInterface   []interface{}
	testClusterAnswersConf                           *managementClient.Answer
	testClusterAnswersInterface                      []interface{}
	testClusterQuestionsConf                         []managementClient.Question
	testClusterQuestionsInterface                    []interface{}
	testLocalClusterAuthEndpointConf                 *managementClient.LocalClusterAuthEndpoint
	testLocalClusterAuthEndpointInterface            []interface{}
	testClusterRegistrationTokenConf                 *managementClient.ClusterRegistrationToken
	testClusterRegistrationToken2Conf                *managementClient.ClusterRegistrationToken
	testClusterRegistrationTokenInterface            []interface{}
	testClusterGenerateKubeConfigOutput              *managementClient.GenerateKubeConfigOutput
	testClusterConfAKS                               *Cluster
	testClusterInterfaceAKS                          map[string]interface{}
	testClusterConfEKS                               *Cluster
	testClusterInterfaceEKS                          map[string]interface{}
	testClusterConfEKSV2                             *Cluster
	testClusterInterfaceEKSV2                        map[string]interface{}
	testClusterConfGKE                               *Cluster
	testClusterInterfaceGKE                          map[string]interface{}
	testClusterConfK3S                               *Cluster
	testClusterInterfaceK3S                          map[string]interface{}
	testClusterConfGKEV2                             *Cluster
	testClusterInterfaceGKEV2                        map[string]interface{}
	testClusterConfOKE                               *Cluster
	testClusterInterfaceOKE                          map[string]interface{}
	testClusterConfRKE                               *Cluster
	testClusterInterfaceRKE                          map[string]interface{}
	testClusterConfRKE2                              *Cluster
	testClusterInterfaceRKE2                         map[string]interface{}
	testClusterConfTemplate                          *Cluster
	testClusterInterfaceTemplate                     map[string]interface{}
)

func testCluster() {
	testClusterEnvVarsConf = []managementClient.EnvVar{
		{
			Name:  "name1",
			Value: "value1",
		},
		{
			Name:  "name2",
			Value: "value2",
		},
	}
	testClusterEnvVarsInterface = []interface{}{
		map[string]interface{}{
			"name":  "name1",
			"value": "value1",
		},
		map[string]interface{}{
			"name":  "name2",
			"value": "value2",
		},
	}

	// fleet agent customization
	testClusterAppendTolerations := []managementClient.Toleration{{
		Effect:   "NoSchedule",
		Key:      "tolerate/test",
		Operator: "Equal",
		Value:    "true",
	},
	}
	testClusterOverrideAffinity := &managementClient.Affinity{
		NodeAffinity: &managementClient.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &managementClient.NodeSelector{
				NodeSelectorTerms: []managementClient.NodeSelectorTerm{
					{
						MatchExpressions: []managementClient.NodeSelectorRequirement{
							{
								Key:      "not.this/nodepool",
								Operator: "NotIn",
								Values:   []string{"true"},
							},
						},
					},
				},
			},
		},
	}
	testClusterOverrideResourceRequirements := &managementClient.ResourceRequirements{
		Limits: map[string]string{
			"cpu":    "500",
			"memory": "500",
		},
		Requests: map[string]string{
			"cpu":    "500",
			"memory": "500",
		},
	}
	testFleetAgentDeploymentCustomizationConf = &managementClient.AgentDeploymentCustomization{
		AppendTolerations:            testClusterAppendTolerations,
		OverrideAffinity:             testClusterOverrideAffinity,
		OverrideResourceRequirements: testClusterOverrideResourceRequirements,
	}
	testFleetAgentDeploymentCustomizationInterface = []interface{}{
		map[string]interface{}{
			"append_tolerations": []interface{}{
				map[string]interface{}{
					"effect":   "NoSchedule",
					"key":      "tolerate/test",
					"operator": "Equal",
					"value":    "true",
				},
			},
			"override_affinity": "{\"nodeAffinity\":{\"requiredDuringSchedulingIgnoredDuringExecution\":{\"nodeSelectorTerms\":[{\"matchExpressions\":[{\"key\":\"not.this/nodepool\",\"operator\":\"NotIn\",\"values\":[\"true\"]}]}]}}}",
			"override_resource_requirements": []interface{}{
				map[string]interface{}{
					"cpu_limit":      "500",
					"cpu_request":    "500",
					"memory_limit":   "500",
					"memory_request": "500",
				},
			},
		},
	}

	// cluster agent customization
	testClusterAgentDeploymentCustomizationConf = &managementClient.AgentDeploymentCustomization{
		AppendTolerations:            testClusterAppendTolerations,
		OverrideAffinity:             testClusterOverrideAffinity,
		OverrideResourceRequirements: testClusterOverrideResourceRequirements,
		SchedulingCustomization: &managementClient.AgentSchedulingCustomization{
			// note: Priority Classes are intentionally omitted, as the generated v3 cluster object
			// uses an int64 which causes a panic when parsing the interface to cty.Value objects.
			// see: https://github.com/hashicorp/terraform-plugin-sdk/blob/2c03a32a9d1be63a12eb18aaf12d2c5270c42346/internal/configs/hcl2shim/values.go#L204-L228
			PodDisruptionBudget: &managementClient.PodDisruptionBudgetSpec{
				MinAvailable: "1",
			},
		},
	}

	testClusterAgentDeploymentCustomizationInterface = []interface{}{
		map[string]interface{}{
			"append_tolerations": []interface{}{
				map[string]interface{}{
					"effect":   "NoSchedule",
					"key":      "tolerate/test",
					"operator": "Equal",
					"value":    "true",
				},
			},
			"scheduling_customization": []interface{}{
				map[string]interface{}{
					"pod_disruption_budget": []interface{}{
						map[string]interface{}{
							"min_available": "1",
						},
					},
				},
			},
			"override_affinity": "{\"nodeAffinity\":{\"requiredDuringSchedulingIgnoredDuringExecution\":{\"nodeSelectorTerms\":[{\"matchExpressions\":[{\"key\":\"not.this/nodepool\",\"operator\":\"NotIn\",\"values\":[\"true\"]}]}]}}}",
			"override_resource_requirements": []interface{}{
				map[string]interface{}{
					"cpu_limit":      "500",
					"cpu_request":    "500",
					"memory_limit":   "500",
					"memory_request": "500",
				},
			},
		},
	}

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
		ClusterID:                  "cluster_test",
		Name:                       clusterRegistrationTokenName,
		Command:                    "command",
		InsecureCommand:            "insecure_command",
		InsecureNodeCommand:        "insecure_node_command",
		InsecureWindowsNodeCommand: "insecure_windows_node_command",
		ManifestURL:                "manifest",
		NodeCommand:                "node_command",
		Token:                      "token",
		WindowsNodeCommand:         "win_node_command",
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
			"id":                            "id",
			"cluster_id":                    "cluster_test",
			"name":                          clusterRegistrationTokenName,
			"command":                       "command",
			"insecure_command":              "insecure_command",
			"insecure_node_command":         "insecure_node_command",
			"insecure_windows_node_command": "insecure_windows_node_command",
			"manifest_url":                  "manifest",
			"node_command":                  "node_command",
			"token":                         "token",
			"windows_node_command":          "win_node_command",
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
	testClusterConfEKSV2 = &Cluster{}
	testClusterConfEKSV2.EKSConfig = testClusterEKSConfigV2Conf
	testClusterConfEKSV2.Name = "test"
	testClusterConfEKSV2.Description = "description"
	testClusterConfEKSV2.Driver = clusterDriverEKSV2
	testClusterConfEKSV2.AgentEnvVars = testClusterEnvVarsConf
	testClusterConfEKSV2.ClusterAgentDeploymentCustomization = testClusterAgentDeploymentCustomizationConf
	testClusterConfEKSV2.FleetAgentDeploymentCustomization = testFleetAgentDeploymentCustomizationConf
	testClusterConfEKSV2.DefaultPodSecurityAdmissionConfigurationTemplateName = "default_pod_security_admission_configuration_template_name"
	testClusterConfEKSV2.EnableNetworkPolicy = newTrue()
	testClusterConfEKSV2.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceEKSV2 = map[string]interface{}{
		"id":                                     "id",
		"name":                                   "test",
		"agent_env_vars":                         testClusterEnvVarsInterface,
		"cluster_agent_deployment_customization": testClusterAgentDeploymentCustomizationInterface,
		"fleet_agent_deployment_customization":   testFleetAgentDeploymentCustomizationInterface,
		"default_project_id":                     "default_project_id",
		"description":                            "description",
		"cluster_auth_endpoint":                  testLocalClusterAuthEndpointInterface,
		"cluster_registration_token":             testClusterRegistrationTokenInterface,
		"default_pod_security_admission_configuration_template_name": "default_pod_security_admission_configuration_template_name",
		"enable_network_policy": true,
		"kube_config":           "kube_config",
		"driver":                clusterDriverEKSV2,
		"eks_config_v2":         testClusterEKSConfigV2Interface,
		"system_project_id":     "system_project_id",
	}
	testClusterConfK3S = &Cluster{}
	testClusterConfK3S.Name = "test"
	testClusterConfK3S.Description = "description"
	testClusterConfK3S.K3sConfig = testClusterK3SConfigConf
	testClusterConfK3S.Driver = clusterDriverK3S
	testClusterConfK3S.AgentEnvVars = testClusterEnvVarsConf
	testClusterConfK3S.DefaultPodSecurityAdmissionConfigurationTemplateName = "default_pod_security_admission_configuration_template_name"
	testClusterConfK3S.EnableNetworkPolicy = newTrue()
	testClusterConfK3S.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceK3S = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"agent_env_vars":             testClusterEnvVarsInterface,
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_admission_configuration_template_name": "default_pod_security_admission_configuration_template_name",
		"enable_network_policy":    true,
		"kube_config":              "kube_config",
		"driver":                   clusterDriverK3S,
		"k3s_config":               testClusterK3SConfigInterface,
		"system_project_id":        "system_project_id",
		"windows_prefered_cluster": false,
	}
	testClusterConfGKEV2 = &Cluster{}
	testClusterConfGKEV2.GKEConfig = testClusterGKEConfigV2Conf
	testClusterConfGKEV2.Name = "test"
	testClusterConfGKEV2.Description = "description"
	testClusterConfGKEV2.Driver = clusterDriverGKEV2
	testClusterConfGKEV2.AgentEnvVars = testClusterEnvVarsConf
	testClusterConfGKEV2.DefaultPodSecurityAdmissionConfigurationTemplateName = "default_pod_security_admission_configuration_template_name"
	testClusterConfGKEV2.EnableNetworkPolicy = newTrue()
	testClusterConfGKEV2.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceGKEV2 = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"agent_env_vars":             testClusterEnvVarsInterface,
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_admission_configuration_template_name": "default_pod_security_admission_configuration_template_name",
		"enable_network_policy": true,
		"kube_config":           "kube_config",
		"driver":                clusterDriverGKEV2,
		"gke_config_v2":         testClusterGKEConfigV2Interface,
		"system_project_id":     "system_project_id",
	}
	testClusterConfRKE = &Cluster{}
	testClusterConfRKE.Name = "test"
	testClusterConfRKE.Description = "description"
	testClusterConfRKE.RancherKubernetesEngineConfig = testClusterRKEConfigConf
	testClusterConfRKE.Driver = clusterDriverRKE
	testClusterConfRKE.AgentEnvVars = testClusterEnvVarsConf
	testClusterConfRKE.ClusterAgentDeploymentCustomization = testClusterAgentDeploymentCustomizationConf
	testClusterConfRKE.FleetAgentDeploymentCustomization = testFleetAgentDeploymentCustomizationConf
	testClusterConfRKE.DefaultPodSecurityAdmissionConfigurationTemplateName = "default_pod_security_admission_configuration_template_name"
	testClusterConfRKE.FleetWorkspaceName = "fleet-test"
	testClusterConfRKE.EnableNetworkPolicy = newTrue()
	testClusterConfRKE.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceRKE = map[string]interface{}{
		"id":                                     "id",
		"name":                                   "test",
		"agent_env_vars":                         testClusterEnvVarsInterface,
		"cluster_agent_deployment_customization": testClusterAgentDeploymentCustomizationInterface,
		"fleet_agent_deployment_customization":   testFleetAgentDeploymentCustomizationInterface,
		"default_project_id":                     "default_project_id",
		"description":                            "description",
		"cluster_auth_endpoint":                  testLocalClusterAuthEndpointInterface,
		"cluster_registration_token":             testClusterRegistrationTokenInterface,
		"default_pod_security_admission_configuration_template_name": "default_pod_security_admission_configuration_template_name",
		"enable_network_policy":    true,
		"fleet_workspace_name":     "fleet-test",
		"kube_config":              "kube_config",
		"driver":                   clusterDriverRKE,
		"rke_config":               testClusterRKEConfigInterface,
		"system_project_id":        "system_project_id",
		"windows_prefered_cluster": false,
	}
	testClusterConfRKE2 = &Cluster{}
	testClusterConfRKE2.Name = "test"
	testClusterConfRKE2.Description = "description"
	testClusterConfRKE2.Rke2Config = testClusterRKE2ConfigConf
	testClusterConfRKE2.Driver = clusterDriverRKE2
	testClusterConfRKE2.AgentEnvVars = testClusterEnvVarsConf
	testClusterConfRKE2.ClusterAgentDeploymentCustomization = testClusterAgentDeploymentCustomizationConf
	testClusterConfRKE2.FleetAgentDeploymentCustomization = testFleetAgentDeploymentCustomizationConf
	testClusterConfRKE2.DefaultPodSecurityAdmissionConfigurationTemplateName = "default_pod_security_admission_configuration_template_name"
	testClusterConfRKE2.EnableNetworkPolicy = newTrue()
	testClusterConfRKE2.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceRKE2 = map[string]interface{}{
		"id":                                     "id",
		"name":                                   "test",
		"agent_env_vars":                         testClusterEnvVarsInterface,
		"cluster_agent_deployment_customization": testClusterAgentDeploymentCustomizationInterface,
		"fleet_agent_deployment_customization":   testFleetAgentDeploymentCustomizationInterface,
		"default_project_id":                     "default_project_id",
		"description":                            "description",
		"cluster_auth_endpoint":                  testLocalClusterAuthEndpointInterface,
		"cluster_registration_token":             testClusterRegistrationTokenInterface,
		"default_pod_security_admission_configuration_template_name": "default_pod_security_admission_configuration_template_name",
		"enable_network_policy":    true,
		"kube_config":              "kube_config",
		"driver":                   clusterDriverRKE2,
		"rke2_config":              testClusterRKE2ConfigInterface,
		"system_project_id":        "system_project_id",
		"windows_prefered_cluster": false,
	}
	testClusterConfTemplate = &Cluster{}
	testClusterConfTemplate.Name = "test"
	testClusterConfTemplate.Description = "description"
	testClusterConfTemplate.ClusterTemplateAnswers = testClusterAnswersConf
	testClusterConfTemplate.ClusterTemplateID = "cluster_template_id"
	testClusterConfTemplate.ClusterTemplateQuestions = testClusterQuestionsConf
	testClusterConfTemplate.ClusterTemplateRevisionID = "cluster_template_revision_id"
	testClusterConfTemplate.Driver = clusterDriverRKE
	testClusterConfTemplate.AgentEnvVars = testClusterEnvVarsConf
	testClusterConfTemplate.DefaultPodSecurityAdmissionConfigurationTemplateName = "default_pod_security_admission_configuration_template_name"
	testClusterConfTemplate.EnableNetworkPolicy = newTrue()
	testClusterConfTemplate.LocalClusterAuthEndpoint = testLocalClusterAuthEndpointConf
	testClusterInterfaceTemplate = map[string]interface{}{
		"id":                         "id",
		"name":                       "test",
		"agent_env_vars":             testClusterEnvVarsInterface,
		"default_project_id":         "default_project_id",
		"description":                "description",
		"cluster_auth_endpoint":      testLocalClusterAuthEndpointInterface,
		"cluster_registration_token": testClusterRegistrationTokenInterface,
		"default_pod_security_admission_configuration_template_name": "default_pod_security_admission_configuration_template_name",
		"enable_network_policy":        true,
		"kube_config":                  "kube_config",
		"driver":                       clusterDriverRKE,
		"cluster_template_answers":     testClusterAnswersInterface,
		"cluster_template_id":          "cluster_template_id",
		"cluster_template_questions":   testClusterQuestionsInterface,
		"cluster_template_revision_id": "cluster_template_revision_id",
		"rke_config":                   []interface{}{},
		"system_project_id":            "system_project_id",
		"windows_prefered_cluster":     false,
	}
}

func TestFlattenClusterRegistrationToken(t *testing.T) {
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
		output, err := flattenClusterRegistrationToken(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenCluster(t *testing.T) {
	testCluster()
	cases := []struct {
		Input          *Cluster
		InputToken     *managementClient.ClusterRegistrationToken
		InputKube      *managementClient.GenerateKubeConfigOutput
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterConfK3S,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceK3S,
		},
		{
			testClusterConfGKEV2,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceGKEV2,
		},
		{
			testClusterConfRKE,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceRKE,
		},
		{
			testClusterConfRKE2,
			testClusterRegistrationTokenConf,
			testClusterGenerateKubeConfigOutput,
			testClusterInterfaceRKE2,
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
		err := flattenCluster(output, tc.Input, tc.InputToken, tc.InputKube, tc.ExpectedOutput["default_project_id"].(string), tc.ExpectedOutput["system_project_id"].(string))
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)

		}
		if tc.ExpectedOutput["driver"] == clusterDriverRKE {
			expectedOutput["rke_config"], _ = flattenClusterRKEConfig(tc.Input.RancherKubernetesEngineConfig, []interface{}{})
		}
		if tc.ExpectedOutput["cluster_agent_deployment_customization"] != nil {
			expectedOutput["cluster_agent_deployment_customization"] = flattenAgentDeploymentCustomization(tc.Input.ClusterAgentDeploymentCustomization, true)
		}
		if tc.ExpectedOutput["fleet_agent_deployment_customization"] != nil {
			expectedOutput["fleet_agent_deployment_customization"] = flattenAgentDeploymentCustomization(tc.Input.FleetAgentDeploymentCustomization, false)
		}
		expectedOutput["id"] = "id"

		assert.Equal(t, tc.ExpectedOutput, expectedOutput, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRegistrationToken(t *testing.T) {
	testCluster()
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
		output, err := expandClusterRegistrationToken(tc.Input, tc.ExpectedOutput.ClusterID)
		tc.ExpectedOutput.ID = "id"
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandCluster(t *testing.T) {
	testCluster()
	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *Cluster
	}{
		{
			testClusterInterfaceEKSV2,
			testClusterConfEKSV2,
		},
		{
			testClusterInterfaceK3S,
			testClusterConfK3S,
		},
		{
			testClusterInterfaceGKEV2,
			testClusterConfGKEV2,
		},
		{
			testClusterInterfaceRKE,
			testClusterConfRKE,
		},
		{
			testClusterInterfaceRKE2,
			testClusterConfRKE2,
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
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestFlattenClusterWithPreservedClusterTemplateAnswers(t *testing.T) {
	testCluster()
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
		err := flattenCluster(output, tc.Input, tc.InputToken, tc.InputKube, tc.ExpectedOutput["default_project_id"].(string), tc.ExpectedOutput["system_project_id"].(string))
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)

		}
		if tc.ExpectedOutput["driver"] == clusterDriverRKE {
			expectedOutput["rke_config"], _ = flattenClusterRKEConfig(tc.Input.RancherKubernetesEngineConfig, []interface{}{})
		}
		expectedOutput["id"] = "id"

		assert.Equal(t, tc.ExpectedOutput, expectedOutput, "Unexpected output from flattener.")
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
	assert.Equal(t, expectedOutput, result, "Unexpected result from preserved answers.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}
