package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testProjectContainerResourceLimitConf           *managementClient.ContainerResourceLimit
	testProjectContainerResourceLimitInterface      []interface{}
	testProjectResourceQuotaLimitConf               *managementClient.ResourceQuotaLimit
	testProjectResourceQuotaLimitInterface          []interface{}
	testProjectResourceQuotaLimitNamespaceConf      *managementClient.ResourceQuotaLimit
	testProjectResourceQuotaLimitNamespaceInterface []interface{}
	testProjectResourceQuotaConf                    *managementClient.ProjectResourceQuota
	testProjectNamespaceResourceQuotaConf           *managementClient.NamespaceResourceQuota
	testProjectResourceQuotaInterface               []interface{}
	testProjectConf                                 *managementClient.Project
	testProjectInterface                            map[string]interface{}
)

func init() {
	testProjectContainerResourceLimitConf = &managementClient.ContainerResourceLimit{
		LimitsCPU:      "limits_cpu",
		LimitsMemory:   "limits_memory",
		RequestsCPU:    "requests_cpu",
		RequestsMemory: "requests_memory",
	}
	testProjectContainerResourceLimitInterface = []interface{}{
		map[string]interface{}{
			"limits_cpu":      "limits_cpu",
			"limits_memory":   "limits_memory",
			"requests_cpu":    "requests_cpu",
			"requests_memory": "requests_memory",
		},
	}
	testProjectResourceQuotaLimitConf = &managementClient.ResourceQuotaLimit{
		ConfigMaps:             "config",
		LimitsCPU:              "cpu",
		LimitsMemory:           "memory",
		PersistentVolumeClaims: "pvc",
		Pods:                   "pods",
		ReplicationControllers: "rc",
		RequestsCPU:            "r_cpu",
		RequestsMemory:         "r_memory",
		RequestsStorage:        "r_storage",
		Secrets:                "secrets",
		Services:               "services",
		ServicesLoadBalancers:  "lb",
		ServicesNodePorts:      "np",
		Extended: map[string]string{
			"count/gpu": "anumber",
		},
	}
	testProjectResourceQuotaLimitInterface = []interface{}{
		map[string]interface{}{
			"config_maps":              "config",
			"limits_cpu":               "cpu",
			"limits_memory":            "memory",
			"persistent_volume_claims": "pvc",
			"pods":                     "pods",
			"replication_controllers":  "rc",
			"requests_cpu":             "r_cpu",
			"requests_memory":          "r_memory",
			"requests_storage":         "r_storage",
			"secrets":                  "secrets",
			"services":                 "services",
			"services_load_balancers":  "lb",
			"services_node_ports":      "np",
			"extended": map[string]string{
				"count/gpu": "anumber",
			},
		},
	}
	testProjectResourceQuotaLimitNamespaceConf = &managementClient.ResourceQuotaLimit{
		ConfigMaps:             "config",
		LimitsCPU:              "cpu",
		LimitsMemory:           "memory",
		PersistentVolumeClaims: "pvc",
		Pods:                   "pods",
		ReplicationControllers: "rc",
		RequestsCPU:            "r_cpu",
		RequestsMemory:         "r_memory",
		RequestsStorage:        "r_storage",
		Secrets:                "secrets",
		Services:               "services",
		ServicesLoadBalancers:  "lb",
		ServicesNodePorts:      "np",
	}
	testProjectResourceQuotaLimitNamespaceInterface = []interface{}{
		map[string]interface{}{
			"config_maps":              "config",
			"limits_cpu":               "cpu",
			"limits_memory":            "memory",
			"persistent_volume_claims": "pvc",
			"pods":                     "pods",
			"replication_controllers":  "rc",
			"requests_cpu":             "r_cpu",
			"requests_memory":          "r_memory",
			"requests_storage":         "r_storage",
			"secrets":                  "secrets",
			"services":                 "services",
			"services_load_balancers":  "lb",
			"services_node_ports":      "np",
		},
	}
	testProjectResourceQuotaConf = &managementClient.ProjectResourceQuota{
		Limit: testProjectResourceQuotaLimitConf,
	}
	testProjectNamespaceResourceQuotaConf = &managementClient.NamespaceResourceQuota{
		Limit: testProjectResourceQuotaLimitNamespaceConf,
	}
	testProjectResourceQuotaInterface = []interface{}{
		map[string]interface{}{
			"project_limit":           testProjectResourceQuotaLimitInterface,
			"namespace_default_limit": testProjectResourceQuotaLimitNamespaceInterface,
		},
	}
	testProjectConf = &managementClient.Project{
		ClusterID:                     "cluster-test",
		Name:                          "test",
		ContainerDefaultResourceLimit: testProjectContainerResourceLimitConf,
		Description:                   "description",
		ResourceQuota:                 testProjectResourceQuotaConf,
		NamespaceDefaultResourceQuota: testProjectNamespaceResourceQuotaConf,
	}
	testProjectInterface = map[string]interface{}{
		"cluster_id":               "cluster-test",
		"name":                     "test",
		"container_resource_limit": testProjectContainerResourceLimitInterface,
		"description":              "description",
		"resource_quota":           testProjectResourceQuotaInterface,
	}
}

func TestFlattenProjectContainerResourceLimit(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ContainerResourceLimit
		ExpectedOutput []interface{}
	}{
		{
			testProjectContainerResourceLimitConf,
			testProjectContainerResourceLimitInterface,
		},
	}

	for _, tc := range cases {
		output := flattenProjectContainerResourceLimit(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenProjectResourceQuotaLimit(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ResourceQuotaLimit
		ExpectedOutput []interface{}
	}{
		{
			testProjectResourceQuotaLimitConf,
			testProjectResourceQuotaLimitInterface,
		},
	}

	for _, tc := range cases {
		output := flattenProjectResourceQuotaLimit(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenProjectResourceQuota(t *testing.T) {

	cases := []struct {
		Input1         *managementClient.ProjectResourceQuota
		Input2         *managementClient.NamespaceResourceQuota
		ExpectedOutput []interface{}
	}{
		{
			testProjectResourceQuotaConf,
			testProjectNamespaceResourceQuotaConf,
			testProjectResourceQuotaInterface,
		},
	}

	for _, tc := range cases {
		output := flattenProjectResourceQuota(tc.Input1, tc.Input2)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenProject(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Project
		ExpectedOutput map[string]interface{}
	}{
		{
			testProjectConf,
			testProjectInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, projectFields(), map[string]interface{}{})
		err := flattenProject(output, tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			assert.FailNow(t, "Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, output)
		}
	}
}

func TestExpandProjectContainerResourceLimit(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ContainerResourceLimit
	}{
		{
			testProjectContainerResourceLimitInterface,
			testProjectContainerResourceLimitConf,
		},
	}

	for _, tc := range cases {
		output := expandProjectContainerResourceLimit(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			assert.FailNow(t, "Unexpected output from expander.\nExpected: %#v\nGiven: %#v", tc.ExpectedOutput, output)
		}
	}
}

func TestExpandProjectResourceQuotaLimit(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ResourceQuotaLimit
	}{
		{
			testProjectResourceQuotaLimitInterface,
			testProjectResourceQuotaLimitConf,
		},
	}

	for _, tc := range cases {
		output := expandProjectResourceQuotaLimit(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandProjectResourceQuota(t *testing.T) {

	cases := []struct {
		Input           []interface{}
		ExpectedOutput1 *managementClient.ProjectResourceQuota
		ExpectedOutput2 *managementClient.NamespaceResourceQuota
	}{
		{
			testProjectResourceQuotaInterface,
			testProjectResourceQuotaConf,
			testProjectNamespaceResourceQuotaConf,
		},
	}

	for _, tc := range cases {
		output1, output2 := expandProjectResourceQuota(tc.Input)
		if !reflect.DeepEqual(output1, tc.ExpectedOutput1) {
			assert.FailNow(t, "Unexpected output from expander on project quota.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput1, output1)
		}
		if !reflect.DeepEqual(output2, tc.ExpectedOutput2) {
			assert.FailNow(t, "Unexpected output from expander on namespace quouta.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput2, output2)
		}
	}
}

func TestExpandProject(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.Project
	}{
		{
			testProjectInterface,
			testProjectConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, projectFields(), tc.Input)
		output := expandProject(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
