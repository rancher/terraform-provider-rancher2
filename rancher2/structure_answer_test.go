package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testAnswersConf      []managementClient.Answer
	testAnswersInterface []interface{}
)

func init() {
	testAnswersConf = []managementClient.Answer{
		{
			ClusterID: "cluster_id",
			ProjectID: "project_id",
			Values: map[string]string{
				"value1": "one",
				"value2": "two",
			},
		},
	}
	testAnswersInterface = []interface{}{
		map[string]interface{}{
			"cluster_id": "cluster_id",
			"project_id": "project_id",
			"values": map[string]interface{}{
				"value1": "one",
				"value2": "two",
			},
		},
	}
}

func TestFlattenAnswers(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Answer
		ExpectedOutput []interface{}
	}{
		{
			testAnswersConf,
			testAnswersInterface,
		},
	}

	for _, tc := range cases {
		output := flattenAnswers(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandAnswers(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Answer
	}{
		{
			testAnswersInterface,
			testAnswersConf,
		},
	}

	for _, tc := range cases {
		output := expandAnswers(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
