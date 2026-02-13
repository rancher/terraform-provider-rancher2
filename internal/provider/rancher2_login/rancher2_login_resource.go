package rancher2_login

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

// The "var _" is a special Go construct that results in an unusable variable.
// The purpose of these lines is to make sure our LocalFileResource correctly implements the "resource.Resource“ interface.
// These will fail at compilation time if the implementation is not satisfied.
var _ resource.Resource = &RancherLoginResource{}
var _ resource.ResourceWithImportState = &RancherLoginResource{}

const (
	endpointPath = "login"
)

func NewRancherLoginResource() resource.Resource {
	return &RancherLoginResource{}
}

type RancherLoginResource struct {
	client c.Client // client is an interface holding a pointer to a struct
}

// RancherLoginResourceModel is in rancher2_login_resource_model.go

func (r *RancherLoginResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_login" // rancher2_login
}

func (r *RancherLoginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rancher Login resource. \n" +
			"This resource manages how the provider logs into Rancher. \n" +
			"It is different from most resources in that it alters the provider's running client.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "The username to use when logging in." +
					"Optionally you can pass this value in an environment variable and it won't be saved in state." +
					"You can control the environment variable with the `username_environment_variable` attribute, which defaults to `RANCHER_USERNAME`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to use when logging in." +
					"Optionally you can pass this value in an environment variable and it won't be saved in state." +
					"You can control the environment variable with the `password_environment_variable` attribute, which defaults to `RANCHER_PASSWORD`.",
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username_environment_variable": schema.StringAttribute{
				MarkdownDescription: "The environment variable where your Rancher username is stored, defaults to `RANCHER_USERNAME`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("RANCHER_USERNAME"),
			},
			"password_environment_variable": schema.StringAttribute{
				MarkdownDescription: "The environment variable where your Rancher password is stored defaults to `RANCHER_PASSWORD`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("RANCHER_PASSWORD"),
			},
			"token_ttl": schema.StringAttribute{
				MarkdownDescription: "The Go time string before the token expires. Defaults to 90 days (`90d`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("90d"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9]+[smhd]$`),
						"must be a valid Go time string (e.g. '90d', '24h', '30m')",
					),
				},
			},
			"refresh_at": schema.StringAttribute{
				MarkdownDescription: "The Go time string before the token expires at which the token will be refreshed. Defaults to 10 days (`10d`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("10d"),
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9]+[smhd]$`),
						"must be a valid Go time string (e.g. '10d', '24h', '30m')",
					),
				},
			},
			"ignore_token": schema.BoolAttribute{
				MarkdownDescription: "Whether to save the resulting token to state." +
					"If set to true this won't save the resulting token to state and will always recreate on plan/apply.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			// read only attributes below here
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier generated on create for the resource." +
					"This value will be used as the name of the user token.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_token": schema.StringAttribute{
				MarkdownDescription: "The user token retrieved from login.",
				Computed:            true,
				Sensitive:           true,
			},
			"session_token": schema.StringAttribute{
				MarkdownDescription: "The session token retrieved from login.",
				Computed:            true,
				Sensitive:           true,
			},
			"user_token_start_date": schema.StringAttribute{
				MarkdownDescription: "The unix timestamp of when the user token was created.",
				Computed:            true,
			},
			"user_token_end_date": schema.StringAttribute{
				MarkdownDescription: "The unix timestamp of when the user token expires.",
				Computed:            true,
			},
			"user_token_refresh_date": schema.StringAttribute{
				MarkdownDescription: "The unix timestamp of when the user token will need to be refreshed. Before this time the resource will only validate the token in state, not refresh it. After this date, the resource will use the current token to attempt to get a new one with the same ttl, 'refreshing' it.",
				Computed:            true,
			},
		},
	}
}

// configure runs at compile time, don't overload the context.
func (r *RancherLoginResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return // Prevent panic if the provider has not been configured.
	}

	// Retrieving the client from the provider.
	client, ok := req.ProviderData.(c.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

// Create generates reality and state to match plan.
func (r *RancherLoginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Create Config: %+v", pp.PrettyPrint(req.Config.Raw)))
	tflog.Debug(ctx, fmt.Sprintf("Create Plan: %+v", pp.PrettyPrint(req.Plan.Raw)))
	var err error

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly.
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}

	plan := RancherLoginResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = validateData(&plan) // we validate data here because tests will bypass the schema validators.
	if err != nil {
		resp.Diagnostics.AddError("Error validating data: ", err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Create Plan after Validate: %s", pp.PrettyPrint(plan)))

	var username string
	var userIsInEnv bool
	var password string
	var passIsInEnv bool

	if plan.Username.IsNull() || plan.Username.IsUnknown() {
		userIsInEnv, username = getValueFromEnv(plan.Username, plan.UsernameEnvironmentVariable.ValueString())
	}
	if !userIsInEnv {
		username = plan.Username.ValueString()
	}
	if plan.Password.IsNull() || plan.Password.IsUnknown() {
		passIsInEnv, password = getValueFromEnv(plan.Password, plan.PasswordEnvironmentVariable.ValueString())
	}
	if !passIsInEnv {
		password = plan.Password.ValueString()
	}

	// Login
	loginReqBody := map[string]any{
		"type":         "localProvider",
		"username":     username,
		"password":     password,
		"responseType": "json",
	}

	loginRequest := c.Request{
		Endpoint: fmt.Sprintf("%s/v1-public/login", client.GetApiUrl()),
		Method:   "POST",
		Body:     loginReqBody,
	}

	loginResponse := c.Response{}
	err = client.Do(ctx, &loginRequest, &loginResponse)
	if err != nil {
		resp.Diagnostics.AddError("Error logging in to Rancher: ", err.Error())
		return
	}

	var loginData struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(loginResponse.Body, &loginData)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshaling login response:", err.Error())
		return
	}

	if loginData.Token == "" {
		resp.Diagnostics.AddError("Login failed", "No session token returned from Rancher")
		return
	}
	plan.SessionToken = types.StringValue(loginData.Token)

	ttlDuration, err := parseCustomDuration(plan.TokenTtl.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error parsing token_ttl duration", err.Error())
		return
	}
	ttlMillis := ttlDuration.Milliseconds()

	// Create the user token
	tokenReqBody := map[string]any{
		"apiVersion": "ext.cattle.io/v1",
		"kind":       "Token",
		"metadata": map[string]any{
			"name": plan.Id.ValueString(),
		},
		"spec": map[string]any{
			"description": "Terraform login token.",
			"ttl":         ttlMillis,
		},
	}
	tokenRequest := c.Request{
		Endpoint: fmt.Sprintf("%s/apis/ext.cattle.io/v1/tokens", client.GetApiUrl()),
		Method:   "POST",
		Body:     tokenReqBody,
		Token:    loginData.Token,
	}

	tokenResponse := c.Response{}
	err = client.Do(ctx, &tokenRequest, &tokenResponse)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user token: ", err.Error())
	}
	diags := processTokenResponse(&plan, tokenResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.IgnoreToken.ValueBool() {
		plan.UserToken = types.StringValue("")
		plan.SessionToken = types.StringValue("")
	}
	if passIsInEnv {
		plan.Password = types.StringValue("")
	}
	if userIsInEnv {
		plan.Username = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("Create State After Set: %+v", pp.PrettyPrint(resp.State.Raw)))
}

// Read updates state to match reality.
// Read runs at refresh time which happens before all other functions and every time another function would be called.
// Don't call this function from one of the other functions (eg. don't call the Read function from within the Create function).
func (r *RancherLoginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Read Request: %+v", pp.PrettyPrint(req.State.Raw)))

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}

	var state RancherLoginResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	validateData(&state)
	tflog.Debug(ctx, fmt.Sprintf("State after validate in read function: %+v", pp.PrettyPrint(state)))

	if state.Id.ValueString() == "" {
		resp.Diagnostics.AddWarning("Id missing", "State missing id during read, recreating.")
		resp.State.RemoveResource(ctx)
		return
	}

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/apis/ext.cattle.io/v1/tokens/%s", client.GetApiUrl(), state.Id.ValueString()),
		Method:   "GET",
		Token:    state.UserToken.ValueString(),
	}
	response := c.Response{}

	err := client.Do(ctx, &request, &response)
	if err != nil {
		if e, ok := err.(*c.ApiError); ok && e.StatusCode == 404 {
			resp.Diagnostics.AddWarning("Token not found", "Got 404 from API when attempting to get token in Read function, recreating.")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading dev resource: ", err.Error())
		return
	}

	diags := processTokenResponse(&state, response)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddWarning("Error processing response", "Error found when processing get token response.")
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Read State After Set: %+v", pp.PrettyPrint(resp.State.Raw)))
}

// Update changes reality and state to match plan (best practice is don't compare old state, just override).
func (r *RancherLoginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Update Request Config: %+v", pp.PrettyPrint(req.Config.Raw)))
	tflog.Debug(ctx, fmt.Sprintf("Update Request Plan: %+v", pp.PrettyPrint(req.Plan.Raw)))
	tflog.Debug(ctx, fmt.Sprintf("Update Request State: %+v", pp.PrettyPrint(req.State.Raw)))
	var err error

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly.
		resp.Diagnostics.AddError("initial client not found, please configure the provider", "")
		return
	}

	var plan RancherLoginResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// attributes that can be updated:
	// username_environment_variable,
	// password_environment_variable,
	// refresh_at

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, plan.Id.ValueString()),
		Method:   "PUT",
		Body:     plan.ToGoModel(ctx),
	}

	response := c.Response{}

	err = client.Do(ctx, &request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dev resource: ", err.Error())
		return
	}

	var respBody RancherLoginModel
	err = json.Unmarshal(response.Body, &respBody)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshalling dev resource: ", err.Error())
		return
	}

	state := *respBody.ToResourceModel(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Update Response State After Set: %+v", pp.PrettyPrint(resp.State.Raw)))
}

// Destroy destroys reality (state is handled automatically).
func (r *RancherLoginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Delete Request State: %+v", pp.PrettyPrint(req.State.Raw)))
	var err error

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("initial client not found, please configure the provider", "")
		return
	}

	var state RancherLoginResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, state.Id.ValueString()),
		Method:   "DELETE",
	}

	response := c.Response{}

	err = client.Do(ctx, &request, &response)
	if err != nil {
		if e, ok := err.(*c.ApiError); ok && e.StatusCode == 404 {
			// resource already deleted
			return
		}
		resp.Diagnostics.AddError("Error deleting dev resource: ", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete Response Object: %+v", pp.PrettyPrint(*resp)))
}

func (r *RancherLoginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getValueFromEnv(s types.String, env string) (bool, string) {
	if os.Getenv(env) != "" {
		return true, os.Getenv(env)
	}
	return false, s.ValueString()
}

// note this function also enforces default values.
func validateData(data *RancherLoginResourceModel) error {
	u, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate UUID for login resource id: %v", err)
	}
	if data.Id.ValueString() == "" {
		data.Id = types.StringValue(u.String())
	}
	if data.UsernameEnvironmentVariable.ValueString() == "" {
		data.UsernameEnvironmentVariable = types.StringValue("RANCHER_USERNAME")
	}
	if data.PasswordEnvironmentVariable.ValueString() == "" {
		data.PasswordEnvironmentVariable = types.StringValue("RANCHER_PASSWORD")
	}
	if data.Username.IsNull() || data.Username.IsUnknown() {
		isInEnv, _ := getValueFromEnv(data.Username, data.UsernameEnvironmentVariable.ValueString())
		if !isInEnv {
			return fmt.Errorf("username not found in config or environment variable")
		}
	}
	if data.Password.IsNull() || data.Password.IsUnknown() {
		isInEnv, _ := getValueFromEnv(data.Password, data.PasswordEnvironmentVariable.ValueString())
		if !isInEnv {
			return fmt.Errorf("password not found in config or environment variable")
		}
	}
	if !data.TokenTtl.IsNull() && !data.TokenTtl.IsUnknown() {
		if _, err := parseCustomDuration(data.TokenTtl.ValueString()); err != nil {
			return fmt.Errorf("invalid token_ttl, must be a valid Go time string (e.g. '90d', '24h', '30m'): %w", err)
		}
	}
	if !data.RefreshAt.IsNull() && !data.RefreshAt.IsUnknown() {
		if _, err := parseCustomDuration(data.RefreshAt.ValueString()); err != nil {
			return fmt.Errorf("invalid refresh_at, must be a valid Go time string (e.g. '10d', '24h', '30m'): %w", err)
		}
	}
	if data.TokenTtl.ValueString() == "" {
		data.TokenTtl = types.StringValue("90d")
	}
	if data.RefreshAt.ValueString() == "" {
		data.RefreshAt = types.StringValue("10d")
	}
	if data.IgnoreToken.IsNull() {
		data.IgnoreToken = types.BoolValue(false)
	}

	return nil
}

func parseCustomDuration(durationStr string) (time.Duration, error) {
	dayRegex := regexp.MustCompile(`(\d+)d`)
	var totalDuration time.Duration

	dayParts := dayRegex.FindAllStringSubmatch(durationStr, -1)
	for _, match := range dayParts {
		if len(match) < 2 {
			continue
		}
		days, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid day value in duration: %s", match[0])
		}
		totalDuration += time.Duration(days) * 24 * time.Hour
	}

	remainingDurationStr := dayRegex.ReplaceAllString(durationStr, "")

	if remainingDurationStr != "" {
		parsedRemaining, err := time.ParseDuration(remainingDurationStr)
		if err != nil {
			if remainingDurationStr == "d" {
				return 0, fmt.Errorf("invalid duration format, standalone 'd' is not supported")
			}
			return 0, fmt.Errorf("invalid time format: %w", err)
		}
		totalDuration += parsedRemaining
	}

	return totalDuration, nil
}

func processTokenResponse(data *RancherLoginResourceModel, tokenResponse c.Response) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error
	var tokenData struct {
		Metadata struct {
			CreationTimestamp string `json:"creationTimestamp"`
			Name              string `json:"name"`
		} `json:"metadata"`
		Status struct {
			BearerToken string `json:"bearerToken"`
			ExpiresAt   string `json:"expiresAt"`
		} `json:"status"`
	}
	err = json.Unmarshal(tokenResponse.Body, &tokenData)
	if err != nil {
		diags.AddError("Error unmarshaling token response:", err.Error())
		return diags
	}
	// id
	data.Id = types.StringValue(tokenData.Metadata.Name)
	// user token
	if tokenData.Status.BearerToken != "" {
		data.UserToken = types.StringValue(tokenData.Status.BearerToken)
	}
	// start date
	creationTime, err := time.Parse(time.RFC3339, tokenData.Metadata.CreationTimestamp)
	if err != nil {
		diags.AddError("Error parsing creation timestamp", err.Error())
		return diags
	}
	data.UserTokenStartDate = types.StringValue(creationTime.Format(time.RFC3339))
	// end date
	expiresAt, err := time.Parse(time.RFC3339, tokenData.Status.ExpiresAt)
	if err != nil {
		diags.AddError("Error parsing token expiration date", err.Error())
		return diags
	}
	data.UserTokenEndDate = types.StringValue(expiresAt.Format(time.RFC3339))
	// refresh date
	refreshDuration, err := parseCustomDuration(data.RefreshAt.ValueString())
	if err != nil {
		diags.AddError("Error parsing refresh_at duration", err.Error())
		return diags
	}
	refreshDate := expiresAt.Add(-refreshDuration)
	data.UserTokenRefreshDate = types.StringValue(refreshDate.Format(time.RFC3339))

	return diags
}
