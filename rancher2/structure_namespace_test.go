package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterClient "github.com/rancher/rancher/pkg/client/generated/cluster/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testNamespaceContainerResourceLimitConf      *clusterClient.ContainerResourceLimit
	testNamespaceContainerResourceLimitInterface []interface{}
	testNamespaceResourceQuotaLimitConf          *clusterClient.ResourceQuotaLimit
	testNamespaceResourceQuotaLimitInterface     []interface{}
	testNamespaceResourceQuotaConf               *clusterClient.NamespaceResourceQuota
	testNamespaceResourceQuotaInterface          []interface{}
	testNamespaceConf                            *clusterClient.Namespace
	testNamespaceInterface                       map[string]interface{}
)

func init() {
	testNamespaceContainerResourceLimitConf = &clusterClient.ContainerResourceLimit{
		LimitsCPU:      "limits_cpu",
		LimitsMemory:   "limits_memory",
		RequestsCPU:    "requests_cpu",
		RequestsMemory: "requests_memory",
	}
	testNamespaceContainerResourceLimitInterface = []interface{}{
		map[string]interface{}{
			"limits_cpu":      "limits_cpu",
			"limits_memory":   "limits_memory",
			"requests_cpu":    "requests_cpu",
			"requests_memory": "requests_memory",
		},
	}
	testNamespaceResourceQuotaLimitConf = &clusterClient.ResourceQuotaLimit{
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
	testNamespaceResourceQuotaLimitInterface = []interface{}{
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
	testNamespaceResourceQuotaConf = &clusterClient.NamespaceResourceQuota{
		Limit: testNamespaceResourceQuotaLimitConf,
	}
	testNamespaceResourceQuotaInterface = []interface{}{
		map[string]interface{}{
			"limit": testNamespaceResourceQuotaLimitInterface,
		},
	}
	testNamespaceConf = &clusterClient.Namespace{
		ProjectID:                     "project:test",
		Name:                          "test",
		ContainerDefaultResourceLimit: testNamespaceContainerResourceLimitConf,
		Description:                   "description",
		ResourceQuota:                 testNamespaceResourceQuotaConf,
	}
	testNamespaceInterface = map[string]interface{}{
		"project_id":               "project:test",
		"name":                     "test",
		"container_resource_limit": testNamespaceContainerResourceLimitInterface,
		"description":              "description",
		"resource_quota":           testNamespaceResourceQuotaInterface,
	}
}

func TestFlattenNamespaceContainerResourceLimit(t *testing.T) {

	cases := []struct {
		Input          *clusterClient.ContainerResourceLimit
		ExpectedOutput []interface{}
	}{
		{
			testNamespaceContainerResourceLimitConf,
			testNamespaceContainerResourceLimitInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNamespaceContainerResourceLimit(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenNamespaceResourceQuotaLimit(t *testing.T) {

	cases := []struct {
		Input          *clusterClient.ResourceQuotaLimit
		ExpectedOutput []interface{}
	}{
		{
			testNamespaceResourceQuotaLimitConf,
			testNamespaceResourceQuotaLimitInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNamespaceResourceQuotaLimit(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenNamespaceResourceQuota(t *testing.T) {

	cases := []struct {
		Input          *clusterClient.NamespaceResourceQuota
		ExpectedOutput []interface{}
	}{
		{
			testNamespaceResourceQuotaConf,
			testNamespaceResourceQuotaInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNamespaceResourceQuota(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenNamespace(t *testing.T) {

	cases := []struct {
		Input          *clusterClient.Namespace
		ExpectedOutput map[string]interface{}
	}{
		{
			testNamespaceConf,
			testNamespaceInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, namespaceFields(), map[string]interface{}{})
		err := flattenNamespace(output, tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}

		assert.Equal(t, tc.ExpectedOutput, expectedOutput, "Unexpected output from flattener.")
	}
}

func TestExpandNamespaceContainerResourceLimit(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *clusterClient.ContainerResourceLimit
	}{
		{
			testNamespaceContainerResourceLimitInterface,
			testNamespaceContainerResourceLimitConf,
		},
	}

	for _, tc := range cases {
		output := expandNamespaceContainerResourceLimit(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandNamespaceResourceQuotaLimit(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *clusterClient.ResourceQuotaLimit
	}{
		{
			testNamespaceResourceQuotaLimitInterface,
			testNamespaceResourceQuotaLimitConf,
		},
	}

	for _, tc := range cases {
		output := expandNamespaceResourceQuotaLimit(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandNamespaceResourceQuota(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *clusterClient.NamespaceResourceQuota
	}{
		{
			testNamespaceResourceQuotaInterface,
			testNamespaceResourceQuotaConf,
		},
	}

	for _, tc := range cases {
		output := expandNamespaceResourceQuota(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandNamespace(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *clusterClient.Namespace
	}{
		{
			testNamespaceInterface,
			testNamespaceConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, namespaceFields(), tc.Input)
		output := expandNamespace(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
