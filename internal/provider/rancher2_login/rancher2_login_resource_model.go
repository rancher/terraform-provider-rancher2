package rancher2_login

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

type RancherLoginResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Username                    types.String `tfsdk:"username"`
	Password                    types.String `tfsdk:"password"`
	UsernameEnvironmentVariable types.String `tfsdk:"username_environment_variable"`
	PasswordEnvironmentVariable types.String `tfsdk:"password_environment_variable"`
	TokenTtl                    types.String `tfsdk:"token_ttl"`
	RefreshAt                   types.String `tfsdk:"refresh_at"`
	IgnoreToken                 types.Bool   `tfsdk:"ignore_token"`
	UserToken                   types.String `tfsdk:"user_token"`
	SessionToken                types.String `tfsdk:"session_token"`
	UserTokenStartDate          types.String `tfsdk:"user_token_start_date"`
	UserTokenEndDate            types.String `tfsdk:"user_token_end_date"`
	UserTokenRefreshDate        types.String `tfsdk:"user_token_refresh_date"`
}

// ToPlan returns a tfsdk.Plan with the data from the model.
func (m *RancherLoginResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
	if diags.HasError() {
		return tfsdk.Plan{}
	}
	r := NewRancherLoginResource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	plan := tfsdk.Plan{
		Schema: s.Schema,
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

// ToState returns a tfsdk.State with the data from the model.
func (m *RancherLoginResourceModel) ToState(ctx context.Context, diags *diag.Diagnostics) tfsdk.State {
	if diags.HasError() {
		return tfsdk.State{}
	}
	diags.AddWarning("Model given ToState:", pp.PrettyPrint(m))
	r := NewRancherLoginResource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	state := tfsdk.State{
		Schema: s.Schema,
	}
	dgs := state.Set(ctx, m)
	if dgs.HasError() {
		diags.Append(dgs...)
	}
	return state
}

// ToGoModel converts a RancherLoginResourceModel to a RancherLoginModel.
func (data *RancherLoginResourceModel) ToGoModel(ctx context.Context) RancherLoginModel {
	return RancherLoginModel{
		Id:                          data.Id.ValueString(),
		Username:                    data.Username.ValueString(),
		Password:                    data.Password.ValueString(),
		UsernameEnvironmentVariable: data.UsernameEnvironmentVariable.ValueString(),
		PasswordEnvironmentVariable: data.PasswordEnvironmentVariable.ValueString(),
		TokenTtl:                    data.TokenTtl.ValueString(),
		RefreshAt:                   data.RefreshAt.ValueString(),
		IgnoreToken:                 data.IgnoreToken.ValueBool(),
		SessionToken:                data.SessionToken.ValueString(),
		UserToken:                   data.UserToken.ValueString(),
		UserTokenStartDate:          data.UserTokenStartDate.ValueString(),
		UserTokenEndDate:            data.UserTokenEndDate.ValueString(),
		UserTokenRefreshDate:        data.UserTokenRefreshDate.ValueString(),
	}
}
