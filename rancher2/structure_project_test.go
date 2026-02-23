package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testProjectContainerResourceLimitConf           *managementClient.ContainerResourceLimit
	testProjectContainerResourceLimitInterface      []any
	testProjectResourceQuotaLimitConf               *managementClient.ResourceQuotaLimit
	testProjectResourceQuotaLimitInterface          []any
	testProjectResourceQuotaLimitNamespaceConf      *managementClient.ResourceQuotaLimit
	testProjectResourceQuotaLimitNamespaceInterface []any
	testProjectResourceQuotaConf                    *managementClient.ProjectResourceQuota
	testProjectNamespaceResourceQuotaConf           *managementClient.NamespaceResourceQuota
	testProjectResourceQuotaInterface               []any
	testProjectConf                                 *managementClient.Project
	testProjectInterface                            map[string]any
)

func init() {
	testProjectContainerResourceLimitConf = &managementClient.ContainerResourceLimit{
		LimitsCPU:      "limits_cpu",
		LimitsMemory:   "limits_memory",
		RequestsCPU:    "requests_cpu",
		RequestsMemory: "requests_memory",
	}
	testProjectContainerResourceLimitInterface = []any{
		map[string]any{
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
	testProjectResourceQuotaLimitInterface = []any{
		map[string]any{
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
			"extended": map[string]any{
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
	testProjectResourceQuotaLimitNamespaceInterface = []any{
		map[string]any{
			"config_maps":              "config",
			"extended":                 map[string]any{},
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
	testProjectResourceQuotaInterface = []any{
		map[string]any{
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
	testProjectInterface = map[string]any{
		"cluster_id":               "cluster-test",
		"name":                     "test",
		"container_resource_limit": testProjectContainerResourceLimitInterface,
		"description":              "description",
		"resource_quota":           testProjectResourceQuotaInterface,
	}
}

func TestFlattenProjectContainerResourceLimit(t *testing.T) {
	output := flattenProjectContainerResourceLimit(testProjectContainerResourceLimitConf)

	assert.Equal(t, testProjectContainerResourceLimitInterface, output, "Unexpected output from flattener")
}

func TestFlattenProjectResourceQuotaLimit(t *testing.T) {
	output := flattenProjectResourceQuotaLimit(testProjectResourceQuotaLimitConf)

	assert.Equal(t, testProjectResourceQuotaLimitInterface, output, "Unexpected output from flattener")
}

func TestFlattenProjectResourceQuota(t *testing.T) {
	output := flattenProjectResourceQuota(testProjectResourceQuotaConf, testProjectNamespaceResourceQuotaConf)
	// testProjectNamespaceResourceQuotaConf has no extended limits and so
	// namespace_default_limit has no extended quotas.
	want := []any{
		map[string]any{
			"project_limit": testProjectResourceQuotaLimitInterface,
			"namespace_default_limit": []any{
				map[string]any{
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
			},
		},
	}

	assert.Equal(t, want, output, "Unexpected output from flattener")
}

func TestFlattenProject(t *testing.T) {
	output := schema.TestResourceDataRaw(t, projectFields(), map[string]any{})
	err := flattenProject(output, testProjectConf)
	assert.NoError(t, err)

	result := map[string]any{}
	for k := range testProjectInterface {
		result[k] = output.Get(k)
	}
	assert.Equal(t, testProjectInterface, result, "Unexpected output from flattener")
}

func TestExpandProjectContainerResourceLimit(t *testing.T) {
	output := expandProjectContainerResourceLimit(testProjectContainerResourceLimitInterface)

	assert.Equal(t, testProjectContainerResourceLimitConf, output, "Unexpected output from expander")
}

func TestExpandProjectResourceQuotaLimit(t *testing.T) {
	output := expandProjectResourceQuotaLimit(testProjectResourceQuotaLimitInterface)

	assert.Equal(t, testProjectResourceQuotaLimitConf, output, "Unexpected output from expander")
}

func TestExpandProjectResourceQuota(t *testing.T) {
	output1, output2 := expandProjectResourceQuota(testProjectResourceQuotaInterface)

	assert.Equal(t, testProjectResourceQuotaConf, output1, "Unexpected output from expander")
	assert.Equal(t, testProjectNamespaceResourceQuotaConf, output2, "Unexpected output from expander on namespace quota")
}

func TestExpandProject(t *testing.T) {
	inputResourceData := schema.TestResourceDataRaw(t, projectFields(), testProjectInterface)
	output := expandProject(inputResourceData)

	assert.Equal(t, testProjectConf, output, "Unexpected output from expander")
}
