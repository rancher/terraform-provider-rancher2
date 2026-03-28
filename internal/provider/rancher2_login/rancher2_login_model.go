package rancher2_login

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RancherLoginModel struct {
	Id                          string `json:"id,omitempty"`
	Username                    string `json:"username,omitempty"`
	Password                    string `json:"password,omitempty"`
	UsernameEnvironmentVariable string `json:"username_environment_variable,omitempty"`
	PasswordEnvironmentVariable string `json:"password_environment_variable,omitempty"`
	TokenTtl                    string `json:"token_ttl,omitempty"`
	RefreshAt                   string `json:"refresh_at,omitempty"`
	IgnoreToken                 bool   `json:"ignore_token,omitempty"`
	UserToken                   string `json:"user_token,omitempty"`
	UserTokenStartDate          string `json:"user_token_start_date,omitempty"`
	UserTokenEndDate            string `json:"user_token_end_date,omitempty"`
	UserTokenRefreshDate        string `json:"user_token_refresh_date,omitempty"`
}

func (obj *RancherLoginModel) ToResourceModel(ctx context.Context, diags *diag.Diagnostics) *Rancher2LoginResourceModel {
	if diags.HasError() {
		return nil
	}
	data := Rancher2LoginResourceModel{
		Id:                          types.StringValue(obj.Id),
		Username:                    types.StringValue(obj.Username),
		Password:                    types.StringValue(obj.Password),
		UsernameEnvironmentVariable: types.StringValue(obj.UsernameEnvironmentVariable),
		PasswordEnvironmentVariable: types.StringValue(obj.PasswordEnvironmentVariable),
		TokenTtl:                    types.StringValue(obj.TokenTtl),
		RefreshAt:                   types.StringValue(obj.RefreshAt),
		IgnoreToken:                 types.BoolValue(obj.IgnoreToken),
		UserToken:                   types.StringValue(obj.UserToken),
		UserTokenStartDate:          types.StringValue(obj.UserTokenStartDate),
		UserTokenEndDate:            types.StringValue(obj.UserTokenEndDate),
		UserTokenRefreshDate:        types.StringValue(obj.UserTokenRefreshDate),
	}

	return &data
}
