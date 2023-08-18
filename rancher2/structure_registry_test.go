package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testRegistryCredentialConf              *projectClient.RegistryCredential
	testRegistryCredentialConfInterface     []interface{}
	testDockerCredentialConf                *projectClient.DockerCredential
	testDockerCredentialInterface           map[string]interface{}
	testNamespacedDockerCredentialConf      *projectClient.NamespacedDockerCredential
	testNamespacedDockerCredentialInterface map[string]interface{}
)

func init() {
	testRegistryCredentialConf = &projectClient.RegistryCredential{
		Username: "username",
		Password: "password",
	}
	testRegistryCredentialConfInterface = []interface{}{
		map[string]interface{}{
			"address":  "address",
			"username": "username",
			"password": "password",
		},
	}
	testDockerCredentialConf = &projectClient.DockerCredential{
		ProjectID:   "project:test",
		Name:        "name",
		Description: "description",
		Registries: map[string]projectClient.RegistryCredential{
			"address": *testRegistryCredentialConf,
		},
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testDockerCredentialInterface = map[string]interface{}{
		"project_id":  "project:test",
		"name":        "name",
		"description": "description",
		"registries":  testRegistryCredentialConfInterface,
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedDockerCredentialConf = &projectClient.NamespacedDockerCredential{
		ProjectID:   "project:test",
		Name:        "name",
		Description: "description",
		NamespaceId: "namespace_id",
		Registries: map[string]projectClient.RegistryCredential{
			"address": *testRegistryCredentialConf,
		},
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedDockerCredentialInterface = map[string]interface{}{
		"project_id":   "project:test",
		"name":         "name",
		"description":  "description",
		"namespace_id": "namespace_id",
		"registries":   testRegistryCredentialConfInterface,
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

func TestFlattenRegistry(t *testing.T) {

	cases := []struct {
		Input          interface{}
		ExpectedOutput map[string]interface{}
	}{
		{
			testDockerCredentialConf,
			testDockerCredentialInterface,
		},
		{
			testNamespacedDockerCredentialConf,
			testNamespacedDockerCredentialInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, registryFields(), tc.ExpectedOutput)
		err := flattenRegistry(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		assert.Equal(t, tc.ExpectedOutput, expectedOutput, "Unexpected output from flattener.")
	}
}

func TestExpandRegistry(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput interface{}
	}{
		{
			testDockerCredentialInterface,
			testDockerCredentialConf,
		},
		{
			testNamespacedDockerCredentialInterface,
			testNamespacedDockerCredentialConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, registryFields(), tc.Input)
		output := expandRegistry(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
