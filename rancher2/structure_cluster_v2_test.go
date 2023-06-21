package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	testClusterV2EnvVarConf                       []rkev1.EnvVar
	testClusterV2EnvVarInterface                  []interface{}
	testClusterV2AgentDeploymentCustomizationConf *provisionv1.AgentDeploymentCustomization
	testClusterV2AgentCustomizationInterface      []interface{}
	testClusterV2Conf                             *ClusterV2
	testClusterV2Interface                        map[string]interface{}
)

func init() {
	testClusterV2EnvVarConf = []rkev1.EnvVar{
		{
			Name:  "name1",
			Value: "value1",
		},
		{
			Name:  "name2",
			Value: "value2",
		},
	}
	testClusterV2EnvVarInterface = []interface{}{
		map[string]interface{}{
			"name":  "name1",
			"value": "value1",
		},
		map[string]interface{}{
			"name":  "name2",
			"value": "value2",
		},
	}
	testClusterV2Conf = &ClusterV2{}

	testClusterV2Conf.TypeMeta.Kind = clusterV2Kind
	testClusterV2Conf.TypeMeta.APIVersion = clusterV2APIVersion

	testClusterV2Conf.ObjectMeta.Name = "name"
	testClusterV2Conf.ObjectMeta.Namespace = "fleet_namespace"
	testClusterV2Conf.ObjectMeta.Annotations = map[string]string{
		"value1": "one",
		"value2": "two",
	}
	testClusterV2Conf.ObjectMeta.Labels = map[string]string{
		"label1": "one",
		"label2": "two",
	}
	testClusterV2Conf.Spec.KubernetesVersion = "kubernetes_version"
	testClusterV2Conf.Spec.LocalClusterAuthEndpoint = testClusterV2LocalAuthEndpointConf
	testClusterV2Conf.Spec.RKEConfig = testClusterV2RKEConfigConf
	testClusterV2Conf.Spec.AgentEnvVars = testClusterV2EnvVarConf
	testClusterV2Conf.Spec.CloudCredentialSecretName = "cloud_credential_secret_name"
	testClusterV2Conf.Spec.DefaultPodSecurityPolicyTemplateName = "default_pod_security_policy_template_name"
	testClusterV2Conf.Spec.DefaultClusterRoleForProjectMembers = "default_cluster_role_for_project_members"
	testClusterV2Conf.Spec.EnableNetworkPolicy = newTrue()

	// cluster and fleet agent customization
	testClusterV2AppendTolerations := []corev1.Toleration{{
		Effect:   corev1.TaintEffectNoSchedule,
		Key:      "tolerate/test",
		Operator: corev1.TolerationOpEqual,
		Value:    "true",
	},
	}
	testClusterV2OverrideAffinity := &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "not.this/nodepool",
								Operator: corev1.NodeSelectorOpNotIn,
								Values:   []string{"true"},
							},
						},
					},
				},
			},
		},
	}
	testVal := "500"
	testQuantity, _ := resource.ParseQuantity(testVal)
	testClusterV2OverrideResourceRequirements := &corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    testQuantity,
			corev1.ResourceMemory: testQuantity,
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    testQuantity,
			corev1.ResourceMemory: testQuantity,
		},
	}
	testClusterV2AgentDeploymentCustomizationConf = &provisionv1.AgentDeploymentCustomization{
		AppendTolerations:            testClusterV2AppendTolerations,
		OverrideAffinity:             testClusterV2OverrideAffinity,
		OverrideResourceRequirements: testClusterV2OverrideResourceRequirements,
	}
	testClusterV2Conf.Spec.ClusterAgentDeploymentCustomization = testClusterV2AgentDeploymentCustomizationConf
	testClusterV2Conf.Spec.FleetAgentDeploymentCustomization = testClusterV2AgentDeploymentCustomizationConf

	testClusterV2AgentCustomizationInterface = []interface{}{
		map[string]interface{}{
			"append_tolerations": []interface{}{
				map[string]interface{}{
					"effect":   "NoSchedule",
					"key":      "tolerate/test",
					"operator": "Equal",
					"seconds":  0,
					"value":    "true",
				},
			},
			"override_affinity": `{
  				"nodeAffinity": {
    				"requiredDuringSchedulingIgnoredDuringExecution": {
      					"nodeSelectorTerms": [
        					{
          						"matchExpressions": [
            						{
              							"key": "not.this/nodepool",
              							"operator": "NotIn",
              							"values": [
                							"true"
              							]
            						}
          						]
        					}
      					]
    				}
  				}
			}`,
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

	testClusterV2Interface = map[string]interface{}{
		"name":                                      "name",
		"fleet_namespace":                           "fleet_namespace",
		"kubernetes_version":                        "kubernetes_version",
		"local_auth_endpoint":                       testClusterV2LocalAuthEndpointInterface,
		"rke_config":                                testClusterV2RKEConfigInterface,
		"agent_env_vars":                            testClusterV2EnvVarInterface,
		"cluster_agent_deployment_customization":    testClusterV2AgentCustomizationInterface,
		"fleet_agent_deployment_customization":      testClusterV2AgentCustomizationInterface,
		"cloud_credential_secret_name":              "cloud_credential_secret_name",
		"default_pod_security_policy_template_name": "default_pod_security_policy_template_name",
		"default_cluster_role_for_project_members":  "default_cluster_role_for_project_members",
		"enable_network_policy":                     true,
		"annotations": map[string]interface{}{
			"value1": "one",
			"value2": "two",
		},
		"labels": map[string]interface{}{
			"label1": "one",
			"label2": "two",
		},
	}
}

func TestFlattenClusterV2(t *testing.T) {

	cases := []struct {
		Input          *ClusterV2
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterV2Conf,
			testClusterV2Interface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterV2Fields(), tc.ExpectedOutput)
		err := flattenClusterV2(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
			if k == "rke_config" {
				// This is a hack to remove the deprecated field because it is not being set.
				rkeConfig := expectedOutput[k].([]interface{})[0].(map[string]interface{})
				delete(rkeConfig, "local_auth_endpoint")
			}
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandClusterV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *ClusterV2
	}{
		{
			testClusterV2Interface,
			testClusterV2Conf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterV2Fields(), tc.Input)
		output, err := expandClusterV2(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
