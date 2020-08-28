package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testMultiClusterAppTargetsConf              []managementClient.Target
	testMultiClusterAppTargetsInterface         []interface{}
	testMultiClusterAppAnswersConf              []managementClient.Answer
	testMultiClusterAppAnswersInterface         []interface{}
	testMultiClusterAppMembersConf              []managementClient.Member
	testMultiClusterAppMembersInterface         []interface{}
	testMultiClusterAppStatusConf               *managementClient.MultiClusterAppStatus
	testMultiClusterAppRollingUpdateConf        *managementClient.RollingUpdate
	testMultiClusterAppRollingUpdateInterface   []interface{}
	testMultiClusterAppUpgradeStrategyConf      *managementClient.UpgradeStrategy
	testMultiClusterAppUpgradeStrategyInterface []interface{}
	testMultiClusterAppExternalID               string
	testMultiClusterAppConf                     *managementClient.MultiClusterApp
	testMultiClusterAppInterface                map[string]interface{}
)

func init() {
	testMultiClusterAppTargetsConf = []managementClient.Target{
		{
			ProjectID:   "project_id",
			AppID:       "app_id",
			Healthstate: "health_state",
			State:       "state",
		},
	}
	testMultiClusterAppTargetsInterface = []interface{}{
		map[string]interface{}{
			"project_id":   "project_id",
			"app_id":       "app_id",
			"health_state": "health_state",
			"state":        "state",
		},
	}
	testMultiClusterAppAnswersConf = []managementClient.Answer{
		{
			ClusterID: "cluster_id",
			ProjectID: "project_id",
			Values: map[string]string{
				"value1": "one",
				"value2": "two",
			},
		},
	}
	testMultiClusterAppAnswersInterface = []interface{}{
		map[string]interface{}{
			"cluster_id": "cluster_id",
			"project_id": "project_id",
			"values": map[string]interface{}{
				"value1": "one",
				"value2": "two",
			},
		},
	}
	testMultiClusterAppMembersConf = []managementClient.Member{
		{
			AccessType:       "access_type",
			GroupPrincipalID: "group_principal_id",
			UserPrincipalID:  "user_principal_id",
		},
	}
	testMultiClusterAppMembersInterface = []interface{}{
		map[string]interface{}{
			"access_type":        "access_type",
			"group_principal_id": "group_principal_id",
			"user_principal_id":  "user_principal_id",
		},
	}
	testMultiClusterAppStatusConf = &managementClient.MultiClusterAppStatus{
		RevisionID: "revision_id",
	}
	testMultiClusterAppRollingUpdateConf = &managementClient.RollingUpdate{
		BatchSize: 10,
		Interval:  10,
	}
	testMultiClusterAppRollingUpdateInterface = []interface{}{
		map[string]interface{}{
			"batch_size": 10,
			"interval":   10,
		},
	}
	testMultiClusterAppUpgradeStrategyConf = &managementClient.UpgradeStrategy{
		RollingUpdate: testMultiClusterAppRollingUpdateConf,
	}
	testMultiClusterAppUpgradeStrategyInterface = []interface{}{
		map[string]interface{}{
			"rolling_update": testMultiClusterAppRollingUpdateInterface,
		},
	}
	testMultiClusterAppExternalID = "catalog://?catalog=test&template=test-demo&version=1.23.0"
	testMultiClusterAppConf = &managementClient.MultiClusterApp{
		Name:                 "foo",
		Targets:              testMultiClusterAppTargetsConf,
		TemplateVersionID:    "cattle-global-data:test-test-demo-1.23.0",
		Answers:              testMultiClusterAppAnswersConf,
		Members:              testMultiClusterAppMembersConf,
		RevisionHistoryLimit: 10,
		Status:               testMultiClusterAppStatusConf,
		Roles: []string{
			"role1",
			"role2",
		},
		UpgradeStrategy: testMultiClusterAppUpgradeStrategyConf,
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testMultiClusterAppInterface = map[string]interface{}{
		"catalog_name":           "test",
		"name":                   "foo",
		"targets":                testMultiClusterAppTargetsInterface,
		"template_name":          "test-demo",
		"answers":                testMultiClusterAppAnswersInterface,
		"members":                testMultiClusterAppMembersInterface,
		"revision_history_limit": 10,
		"revision_id":            "revision_id",
		"roles": []interface{}{
			"role1",
			"role2",
		},
		"template_version": "1.23.0",
		"upgrade_strategy": testMultiClusterAppUpgradeStrategyInterface,
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
}

func TestFlattenMultiClusterApp(t *testing.T) {

	cases := []struct {
		Input          *managementClient.MultiClusterApp
		ExpectedOutput map[string]interface{}
	}{
		{
			testMultiClusterAppConf,
			testMultiClusterAppInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, multiClusterAppFields(), map[string]interface{}{})
		err := flattenMultiClusterApp(output, tc.Input, testMultiClusterAppExternalID)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven: %#v", tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandMultiClusterApp(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.MultiClusterApp
	}{
		{
			testMultiClusterAppInterface,
			testMultiClusterAppConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, multiClusterAppFields(), tc.Input)
		output, err := expandMultiClusterApp(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven: %#v", tc.ExpectedOutput, output)
		}
	}
}
