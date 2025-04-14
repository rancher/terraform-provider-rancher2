package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	testClusterV2EnvVarConf                              []rkev1.EnvVar
	testClusterV2EnvVarInterface                         []interface{}
	testClusterV2ClusterAgentDeploymentCustomizationConf *provisionv1.AgentDeploymentCustomization
	testClusterV2FleetAgentDeploymentCustomizationConf   *provisionv1.AgentDeploymentCustomization
	testClusterV2ClusterAgentCustomizationInterface      []interface{}
	testClusterV2FleetAgentCustomizationInterface        []interface{}
	testClusterV2Conf                                    *ClusterV2
	testClusterV2Interface                               map[string]interface{}
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
	testClusterV2Conf.Spec.DefaultPodSecurityAdmissionConfigurationTemplateName = "default_pod_security_admission_configuration_template_name"
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

	neverPreemption := corev1.PreemptionPolicy("Never")
	testClusterV2ClusterAgentDeploymentCustomizationConf = &provisionv1.AgentDeploymentCustomization{
		AppendTolerations:            testClusterV2AppendTolerations,
		OverrideAffinity:             testClusterV2OverrideAffinity,
		OverrideResourceRequirements: testClusterV2OverrideResourceRequirements,
		SchedulingCustomization: &provisionv1.AgentSchedulingCustomization{
			PriorityClass: &provisionv1.PriorityClassSpec{
				Value:            123,
				PreemptionPolicy: &neverPreemption,
			},
			PodDisruptionBudget: &provisionv1.PodDisruptionBudgetSpec{
				MinAvailable: "1",
			},
		},
	}

	testClusterV2FleetAgentDeploymentCustomizationConf = &provisionv1.AgentDeploymentCustomization{
		AppendTolerations:            testClusterV2AppendTolerations,
		OverrideAffinity:             testClusterV2OverrideAffinity,
		OverrideResourceRequirements: testClusterV2OverrideResourceRequirements,
	}

	testClusterV2Conf.Spec.ClusterAgentDeploymentCustomization = testClusterV2ClusterAgentDeploymentCustomizationConf
	testClusterV2Conf.Spec.FleetAgentDeploymentCustomization = testClusterV2FleetAgentDeploymentCustomizationConf

	testClusterV2FleetAgentCustomizationInterface = []interface{}{
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

	testClusterV2ClusterAgentCustomizationInterface = []interface{}{
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
			"scheduling_customization": []interface{}{
				map[string]interface{}{
					"priority_class": []interface{}{
						map[string]interface{}{
							"value":             123,
							"preemption_policy": "Never",
						},
					},
					"pod_disruption_budget": []interface{}{
						map[string]interface{}{
							"min_available":   "1",
							"max_unavailable": "",
						},
					},
				},
			},
		},
	}

	testClusterV2Interface = map[string]interface{}{
		"name":                                   "name",
		"fleet_namespace":                        "fleet_namespace",
		"kubernetes_version":                     "kubernetes_version",
		"local_auth_endpoint":                    testClusterV2LocalAuthEndpointInterface,
		"rke_config":                             testClusterV2RKEConfigInterface,
		"agent_env_vars":                         testClusterV2EnvVarInterface,
		"cluster_agent_deployment_customization": testClusterV2ClusterAgentCustomizationInterface,
		"fleet_agent_deployment_customization":   testClusterV2FleetAgentCustomizationInterface,
		"cloud_credential_secret_name":           "cloud_credential_secret_name",
		"default_pod_security_admission_configuration_template_name": "default_pod_security_admission_configuration_template_name",
		"default_cluster_role_for_project_members":                   "default_cluster_role_for_project_members",
		"enable_network_policy":                                      true,
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
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		actualOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			actualOutput[k] = output.Get(k)
			if k == "rke_config" {
				// This is a hack to remove the deprecated field because it is not being set.
				rkeConfig := actualOutput[k].([]interface{})[0].(map[string]interface{})
				delete(rkeConfig, "local_auth_endpoint")
			}
		}
		assert.Equal(t, tc.ExpectedOutput, actualOutput)
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
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
