package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	testResourceRequirementsV2Conf      *corev1.ResourceRequirements
	testResourceRequirementsV2Interface []interface{}
)

func init() {
	testVal := "500"
	testQuantity, _ := resource.ParseQuantity(testVal)
	testResourceRequirementsV2Conf = &corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    testQuantity,
			corev1.ResourceMemory: testQuantity,
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    testQuantity,
			corev1.ResourceMemory: testQuantity,
		},
	}
	testResourceRequirementsV2Interface = []interface{}{
		map[string]interface{}{
			"cpu_limit":      testVal,
			"memory_limit":   testVal,
			"cpu_request":    testVal,
			"memory_request": testVal,
		},
	}
}

func TestFlattenResourceRequirementsV2(t *testing.T) {

	cases := []struct {
		Input          *corev1.ResourceRequirements
		ExpectedOutput []interface{}
	}{
		{
			testResourceRequirementsV2Conf,
			testResourceRequirementsV2Interface,
		},
	}

	for _, tc := range cases {
		output := flattenResourceRequirementsV2(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandResourceRequirementsV2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *corev1.ResourceRequirements
	}{
		{
			testResourceRequirementsV2Interface,
			testResourceRequirementsV2Conf,
		},
	}

	for _, tc := range cases {
		output, _ := expandResourceRequirementsV2(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
