package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterTemplateQuestionsConf                        []managementClient.Question
	testClusterTemplateQuestionsInterface                   []interface{}
	testClusterTemplateRevisionsConfigRKEConf               *managementClient.RancherKubernetesEngineConfig
	testClusterTemplateRevisionsConfigAuthEndpointConf      *managementClient.LocalClusterAuthEndpoint
	testClusterTemplateRevisionsConfigRKEInterface          []interface{}
	testClusterTemplateRevisionsConfigAuthEndpointInterface []interface{}
	testClusterTemplateRevisionsConfigConf                  *managementClient.ClusterSpecBase
	testClusterTemplateRevisionsConfigInterface             []interface{}
	testClusterTemplateRevisionsConf                        []managementClient.ClusterTemplateRevision
	testClusterTemplateRevisionsInterface                   []interface{}
	testClusterTemplateMembersConf                          []managementClient.Member
	testClusterTemplateMembersInterface                     []interface{}
	testClusterTemplateConf                                 *managementClient.ClusterTemplate
	testClusterTemplateInterface                            map[string]interface{}
)

func testClusterTemplate() {
	k8sVersion = testAccRancher2ClusterRKEK8SDefaultVersion
	if len(testAccRancher2ClusterRKEK8SDefaultVersion) == 0 {
		k8sVersion = "test"
	}
	testClusterTemplateRevisionsConfigRKEConf = &managementClient.RancherKubernetesEngineConfig{
		AddonJobTimeout:     30,
		Addons:              "addons",
		AddonsInclude:       []string{"addon1", "addon2"},
		Authentication:      testClusterRKEConfigAuthenticationConf,
		Authorization:       testClusterRKEConfigAuthorizationConf,
		BastionHost:         testClusterRKEConfigBastionHostConf,
		CloudProvider:       testClusterRKEConfigCloudProviderConf,
		IgnoreDockerVersion: newTrue(),
		Ingress:             testClusterRKEConfigIngressConf,
		Monitoring:          testClusterRKEConfigMonitoringConf,
		Network:             testClusterRKEConfigNetworkConfCanal,
		Nodes:               testClusterRKEConfigNodesConf,
		PrefixPath:          "terraform-test",
		PrivateRegistries:   testClusterRKEConfigPrivateRegistriesConf,
		Services:            testClusterRKEConfigServicesConf,
		SSHAgentAuth:        true,
		SSHKeyPath:          "/home/user/.ssh",
		Version:             k8sVersion,
	}
	testClusterTemplateRevisionsConfigRKEInterface = []interface{}{
		map[string]interface{}{
			"addon_job_timeout":     30,
			"addons":                "addons",
			"addons_include":        []interface{}{"addon1", "addon2"},
			"authentication":        testClusterRKEConfigAuthenticationInterface,
			"authorization":         testClusterRKEConfigAuthorizationInterface,
			"bastion_host":          testClusterRKEConfigBastionHostInterface,
			"cloud_provider":        testClusterRKEConfigCloudProviderInterface,
			"ignore_docker_version": true,
			"ingress":               testClusterRKEConfigIngressInterface,
			"monitoring":            testClusterRKEConfigMonitoringInterface,
			"network":               testClusterRKEConfigNetworkInterfaceCanal,
			"nodes":                 testClusterRKEConfigNodesInterface,
			"prefix_path":           "terraform-test",
			"private_registries":    testClusterRKEConfigPrivateRegistriesInterface,
			"services":              testClusterRKEConfigServicesInterface,
			"ssh_agent_auth":        true,
			"ssh_key_path":          "/home/user/.ssh",
			"kubernetes_version":    k8sVersion,
		},
	}
	testClusterTemplateQuestionsConf = []managementClient.Question{
		{
			Default:  "default",
			Required: true,
			Type:     "string",
			Variable: "variable",
		},
	}
	testClusterTemplateQuestionsInterface = []interface{}{
		map[string]interface{}{
			"default":  "default",
			"required": true,
			"type":     "string",
			"variable": "variable",
		},
	}
	testClusterTemplateRevisionsConfigAuthEndpointConf = &managementClient.LocalClusterAuthEndpoint{
		CACerts: "cacerts",
		Enabled: true,
		FQDN:    "fqdn",
	}
	testClusterTemplateRevisionsConfigAuthEndpointInterface = []interface{}{
		map[string]interface{}{
			"ca_certs": "cacerts",
			"enabled":  true,
			"fqdn":     "fqdn",
		},
	}
	testClusterTemplateRevisionsConfigConf = &managementClient.ClusterSpecBase{
		DefaultClusterRoleForProjectMembers: "default_cluster_role_for_project_members",
		DefaultPodSecurityPolicyTemplateID:  "default_pod_security_policy_template_id",
		DesiredAgentImage:                   "desired_agent_image",
		DesiredAuthImage:                    "desired_auth_image",
		DockerRootDir:                       "docker_root_dir",
		EnableClusterAlerting:               true,
		EnableClusterMonitoring:             true,
		EnableNetworkPolicy:                 newTrue(),
		LocalClusterAuthEndpoint:            testClusterTemplateRevisionsConfigAuthEndpointConf,
		RancherKubernetesEngineConfig:       testClusterTemplateRevisionsConfigRKEConf,
		WindowsPreferedCluster:              true,
	}
	testClusterTemplateRevisionsConfigInterface = []interface{}{
		map[string]interface{}{
			"cluster_auth_endpoint":                    testClusterTemplateRevisionsConfigAuthEndpointInterface,
			"default_cluster_role_for_project_members": "default_cluster_role_for_project_members",
			"default_pod_security_policy_template_id":  "default_pod_security_policy_template_id",
			"desired_agent_image":                      "desired_agent_image",
			"desired_auth_image":                       "desired_auth_image",
			"docker_root_dir":                          "docker_root_dir",
			"enable_cluster_alerting":                  true,
			"enable_cluster_monitoring":                true,
			"enable_network_policy":                    true,
			"rke_config":                               testClusterTemplateRevisionsConfigRKEInterface,
			"windows_prefered_cluster":                 true,
		},
	}
	testClusterTemplateRevisionsConf = []managementClient.ClusterTemplateRevision{
		{
			ClusterConfig:     testClusterTemplateRevisionsConfigConf,
			ClusterTemplateID: "cluster_template_id",
			Enabled:           newTrue(),
			Name:              "test",
			Questions:         testClusterTemplateQuestionsConf,
		},
	}
	testClusterTemplateRevisionsConf[0].ID = "default_revision_id"
	testClusterTemplateRevisionsInterface = []interface{}{
		map[string]interface{}{
			"id":                  "default_revision_id",
			"cluster_config":      testClusterTemplateRevisionsConfigInterface,
			"cluster_template_id": "cluster_template_id",
			"default":             true,
			"enabled":             true,
			"name":                "test",
			"questions":           testClusterTemplateQuestionsInterface,
		},
	}
	testClusterTemplateMembersConf = []managementClient.Member{
		{
			AccessType:       "access_type",
			GroupPrincipalID: "group_principal_id",
			UserPrincipalID:  "user_principal_id",
		},
	}
	testClusterTemplateMembersInterface = []interface{}{
		map[string]interface{}{
			"access_type":        "access_type",
			"group_principal_id": "group_principal_id",
			"user_principal_id":  "user_principal_id",
		},
	}
	testClusterTemplateConf = &managementClient.ClusterTemplate{
		DefaultRevisionID: "default_revision_id",
		Description:       "description",
		Name:              "name",
		Members:           testClusterTemplateMembersConf,
	}
	testClusterTemplateInterface = map[string]interface{}{
		"default_revision_id": "default_revision_id",
		"description":         "description",
		"members":             testClusterTemplateMembersInterface,
		"name":                "name",
		"template_revisions":  testClusterTemplateRevisionsInterface,
	}
}

func TestFlattenQuestions(t *testing.T) {
	testClusterTemplate()
	cases := []struct {
		Input          []managementClient.Question
		ExpectedOutput []interface{}
	}{
		{
			testClusterTemplateQuestionsConf,
			testClusterTemplateQuestionsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenQuestions(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterSpecBase(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterSpecBase
		ExpectedOutput []interface{}
	}{
		{
			testClusterTemplateRevisionsConfigConf,
			testClusterTemplateRevisionsConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterSpecBase(tc.Input, tc.ExpectedOutput)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterTemplateRevisions(t *testing.T) {

	cases := []struct {
		Input          []managementClient.ClusterTemplateRevision
		ExpectedOutput []interface{}
	}{
		{
			testClusterTemplateRevisionsConf,
			testClusterTemplateRevisionsInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterTemplateRevisions(tc.Input, "default_revision_id", tc.ExpectedOutput)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterTemplate(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterTemplate
		Revisions      []managementClient.ClusterTemplateRevision
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterTemplateConf,
			testClusterTemplateRevisionsConf,
			testClusterTemplateInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterTemplateFields(), tc.ExpectedOutput)
		err := flattenClusterTemplate(output, tc.Input, tc.Revisions)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}

		expectedOutput["template_revisions"].([]interface{})[0].(map[string]interface{})["cluster_config"].([]interface{})[0].(map[string]interface{})["rke_config"], _ = flattenClusterRKEConfig(testClusterTemplateRevisionsConfigRKEConf, []interface{}{})
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandQuestions(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Question
	}{
		{
			testClusterTemplateQuestionsInterface,
			testClusterTemplateQuestionsConf,
		},
	}

	for _, tc := range cases {
		output := expandQuestions(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterSpecBase(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ClusterSpecBase
	}{
		{
			testClusterTemplateRevisionsConfigInterface,
			testClusterTemplateRevisionsConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterSpecBase(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterTemplateRevisions(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.ClusterTemplateRevision
	}{
		{
			testClusterTemplateRevisionsInterface,
			testClusterTemplateRevisionsConf,
		},
	}

	for _, tc := range cases {
		_, output, err := expandClusterTemplateRevisions(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterTemplate(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ClusterTemplate
	}{
		{
			testClusterTemplateInterface,
			testClusterTemplateConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterTemplateFields(), tc.Input)
		_, output, _, err := expandClusterTemplate(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expnader: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
