package rancher2_login

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRancherLoginResourceModel(t *testing.T) {
	t.Run("ToGoModel", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherLoginResourceModel
			want RancherLoginModel
		}{
			{
				"Basic",
				getDefaultResourceModel(),
				getDefaultGoModel(),
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.fit.ToGoModel(context.Background())
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
	t.Run("ToPlan", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherLoginResourceModel
			want tfsdk.Plan
		}{
			{
				"Basic",
				getDefaultResourceModel(),
				getPlan(context.Background(), getDefaultResourceModel(), &diag.Diagnostics{}),
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				model := getDefaultResourceModel()
				got := model.ToPlan(context.Background(), &diag.Diagnostics{})

				if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})
	t.Run("ToState", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherLoginResourceModel
			want tfsdk.State
		}{
			{
				"Basic",
				getDefaultResourceModel(),
				getState(context.Background(), getDefaultResourceModel(), &diag.Diagnostics{}),
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				model := getDefaultResourceModel()
				got := model.ToState(context.Background(), &diag.Diagnostics{})

				if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})

}

// helpers.
func getDefaultResourceModel() RancherLoginResourceModel {
	return RancherLoginResourceModel{
		Id:                          types.StringValue("test"),
		Username:                    types.StringValue("user"),
		Password:                    types.StringValue("pass"),
		UsernameEnvironmentVariable: types.StringValue("RANCHER_USERNAME"),
		PasswordEnvironmentVariable: types.StringValue("RANCHER_PASSWORD"),
		TokenTtl:                    types.StringValue("90d"),
		RefreshAt:                   types.StringValue("10d"),
		IgnoreToken:                 types.BoolValue(false),
		UserToken:                   types.StringValue("user-token"),
		SessionToken:                types.StringValue("session-token"),
		UserTokenStartDate:          types.StringValue("start"),
		UserTokenEndDate:            types.StringValue("end"),
		UserTokenRefreshDate:        types.StringValue("refresh"),
	}
}

func getDefaultGoModel() RancherLoginModel {
	return RancherLoginModel{
		Id:                          "test",
		Username:                    "user",
		Password:                    "pass",
		UsernameEnvironmentVariable: "RANCHER_USERNAME",
		PasswordEnvironmentVariable: "RANCHER_PASSWORD",
		TokenTtl:                    "90d",
		RefreshAt:                   "10d",
		IgnoreToken:                 false,
		UserToken:                   "user-token",
		SessionToken:                "session-token",
		UserTokenStartDate:          "start",
		UserTokenEndDate:            "end",
		UserTokenRefreshDate:        "refresh",
	}
}

func getPlan(ctx context.Context, m RancherLoginResourceModel, diags *diag.Diagnostics) tfsdk.Plan {
	plan := tfsdk.Plan{
		Schema: getSchema(ctx),
	}
	if diags.HasError() {
		return plan
	}

	dgs := plan.Set(ctx, m)
	if dgs.HasError() {
		diags.Append(dgs...)
	}
	return plan
}

func getState(ctx context.Context, m RancherLoginResourceModel, diags *diag.Diagnostics) tfsdk.State {
	state := tfsdk.State{
		Schema: getSchema(ctx),
	}
	if diags.HasError() {
		return state
	}

	dgs := state.Set(ctx, m)
	if dgs.HasError() {
		diags.Append(dgs...)
	}
	return state
}

func getSchema(ctx context.Context) schema.Schema {
	emptyResource := NewRancherLoginResource()
	schemaResponseContainer := &resource.SchemaResponse{}
	emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
	return schemaResponseContainer.Schema
}
