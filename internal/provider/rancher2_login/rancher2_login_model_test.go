package rancher2_login

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRancherLoginModel(t *testing.T) {
	t.Run("ToResourceModel", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherLoginModel
			want RancherLoginResourceModel
		}{
			{
				"Basic",
				RancherLoginModel{
					Id:                          "test",
					Username:                    "user",
					Password:                    "pass",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserToken:                   "user-token",
					UserTokenStartDate:          "start",
					UserTokenEndDate:            "end",
					UserTokenRefreshDate:        "refresh",
				},
				RancherLoginResourceModel{
					Id:                          types.StringValue("test"),
					Username:                    types.StringValue("user"),
					Password:                    types.StringValue("pass"),
					UsernameEnvironmentVariable: types.StringValue("RANCHER_USERNAME"),
					PasswordEnvironmentVariable: types.StringValue("RANCHER_PASSWORD"),
					TokenTtl:                    types.StringValue("90d"),
					RefreshAt:                   types.StringValue("10d"),
					IgnoreToken:                 types.BoolValue(false),
					UserToken:                   types.StringValue("user-token"),
					UserTokenStartDate:          types.StringValue("start"),
					UserTokenEndDate:            types.StringValue("end"),
					UserTokenRefreshDate:        types.StringValue("refresh"),
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.fit.ToResourceModel(context.Background(), &diag.Diagnostics{})
				if diff := cmp.Diff(tc.want, *got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})
}
