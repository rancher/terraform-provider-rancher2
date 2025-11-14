package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupressFunc(t *testing.T) {
	cases := []struct {
		Name           string
		K              string
		Old            string
		New            string
		ExpectedResult bool
	}{
		{
			Name:           "Exception key change should not be suppressed",
			K:              "annotations.rancher.io/imported-cluster-version-management",
			Old:            "true",
			New:            "false",
			ExpectedResult: false,
		},
		{
			Name:           "Exception key removal should not be suppressed",
			K:              "labels.rancher.io/imported-cluster-version-management",
			Old:            "false",
			New:            "",
			ExpectedResult: false,
		},
		{
			Name:           "Rancher.io annotation/label removal should be suppressed",
			K:              "annotations.rancher.io/creator",
			Old:            "user-123",
			New:            "",
			ExpectedResult: true,
		},
		{
			Name:           "Cattle.io annotation/label removal should be suppressed",
			K:              "annotations.cattle.io/some-state",
			Old:            "some-val",
			New:            "",
			ExpectedResult: true,
		},
		{
			Name:           "User-managed change to a rancher.io annotation should not be suppressed",
			K:              "annotations.rancher.io/creator",
			Old:            "user-123",
			New:            "user-456",
			ExpectedResult: false,
		},
		{Name: "User-managed change to a cattle.io annotation should not be suppressed",
			K:              "annotations.cattle.io/creator",
			Old:            "user-123",
			New:            "user-456",
			ExpectedResult: false,
		},
		{
			Name:           "Non-Rancher annotation should not be suppressed",
			K:              "annotations.my-custom-annotation",
			Old:            "old-val",
			New:            "new-val",
			ExpectedResult: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result := supressFunc(tc.K, tc.Old, tc.New, nil)
			assert.Equal(t, tc.ExpectedResult, result)
		})
	}
}
