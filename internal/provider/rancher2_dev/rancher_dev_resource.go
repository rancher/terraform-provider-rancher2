package rancher2_dev

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

// The "var _" is a special Go construct that results in an unusable variable.
// The purpose of these lines is to make sure our LocalFileResource correctly implements the "resource.Resource“ interface.
// These will fail at compilation time if the implementation is not satisfied.
var _ resource.Resource = &RancherDevResource{}
var _ resource.ResourceWithImportState = &RancherDevResource{}

const (
	endpointPath = "dev"
)

func NewRancherDevResource() resource.Resource {
	return &RancherDevResource{}
}

type RancherDevResource struct {
	client c.Client // client is an interface holding a pointer to a struct
}

type RancherDevResourceModel struct {
	Id              types.String `tfsdk:"id"`
	BoolAttribute   types.Bool   `tfsdk:"bool_attribute"`
	NumberAttribute types.Int64  `tfsdk:"number_attribute"`
	StringAttribute types.String `tfsdk:"string_attribute"`
}

type RancherDevModel struct {
	Id              string `json:"id"`
	BoolAttribute   bool   `json:"bool_attribute,omitempty"`
	NumberAttribute int64  `json:"number_attribute,omitempty"`
	StringAttribute string `json:"string_attribute,omitempty"`
}

func (r *RancherDevResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev" // rancher2_dev
}

func (r *RancherDevResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rancher Development resource. \n" +
			"This resource is used as a dummy for development purposes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the resource.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"bool_attribute": schema.BoolAttribute{
				MarkdownDescription: "A boolean attribute.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			// WARNING! Avoid the more generic "number" attribute type as it overcomplicates testing
			"number_attribute": schema.Int64Attribute{
				MarkdownDescription: "A number attribute.",
				Required:            true,
			},
			"string_attribute": schema.StringAttribute{
				MarkdownDescription: "A string attribute with validation.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^dev-.*`),
						"must start with 'dev-'",
					),
				},
			},
		},
	}
}

// configure runs at compile time, don't overload the context
func (r *RancherDevResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return // Prevent panic if the provider has not been configured.
	}

	// Retrieving the client from the provider
	r.client = req.ProviderData.(c.Client)
}

// - Read updates state to match reality.
// - Update changes reality and state to match plan/config (best practice is don't compare old state, just override).
// - Destroy destroys reality (state is handled automatically).

// - Create generates reality and state to match plan.
func (r *RancherDevResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Create Request Object: %+v", pp.PrettyPrint(req)))
	var err error

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("initial client not found, please configure the provider", "")
		return
	}

	var plan RancherDevResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = validateData(&plan) // we validate data here because tests will bypass the schema validators
	if err != nil {
		resp.Diagnostics.AddError("Error validating data: ", err.Error())
		return
	}

	// this object will form the body of the api request via json marshal
	obj := RancherDevModel{
		Id:              plan.Id.ValueString(),
		BoolAttribute:   plan.BoolAttribute.ValueBool(),
		NumberAttribute: plan.NumberAttribute.ValueInt64(),
		StringAttribute: plan.StringAttribute.ValueString(),
	}
	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s", client.GetApiUrl(), endpointPath),
		Method:   "POST",
		Body:     obj,
	}

	response := c.Response{}

	err = client.Do(&request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dev resource: ", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("Create Response Object: %+v", pp.PrettyPrint(*resp)))
}

// Read runs at refresh time which happens before all other functions and every time another function would be called.
func (r *RancherDevResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Read Request Object: %+v", pp.PrettyPrint(req)))

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}

	var state RancherDevResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, state.Id.ValueString()),
		Method:   "GET",
	}
	response := c.Response{}

	err := client.Do(&request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error reading dev resource: ", err.Error())
		return
	}

	var obj RancherDevModel
	err = json.Unmarshal(response.Body, &obj)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshalling dev resource: ", err.Error())
		return
	}

	if obj.Id == "" || obj.Id != state.Id.ValueString() {
		resp.Diagnostics.AddError(
			"Invalid object returned from API",
			fmt.Sprintf("Id must not be empty and must match the Id in the state, found: %s", obj.Id),
		)
		return
	}

	state.Id = types.StringValue(obj.Id)
	state.BoolAttribute = types.BoolValue(obj.BoolAttribute)
	state.NumberAttribute = types.Int64Value(obj.NumberAttribute)
	state.StringAttribute = types.StringValue(obj.StringAttribute)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Read Response Object: %+v", pp.PrettyPrint(*resp)))
}

func (r *RancherDevResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Update Request Object: %+v", pp.PrettyPrint(req)))
	var err error

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("initial client not found, please configure the provider", "")
		return
	}

	var plan RancherDevResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	numberVal := plan.NumberAttribute.ValueInt64()
	obj := RancherDevModel{
		Id:              plan.Id.ValueString(),
		BoolAttribute:   plan.BoolAttribute.ValueBool(),
		NumberAttribute: numberVal,
		StringAttribute: plan.StringAttribute.ValueString(),
	}
	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, plan.Id.ValueString()),
		Method:   "PUT",
		Body:     obj,
	}

	response := c.Response{}

	err = client.Do(&request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dev resource: ", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("Update Response Object: %+v", pp.PrettyPrint(*resp)))
}

func (r *RancherDevResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Delete Request Object: %+v", pp.PrettyPrint(req)))
	var err error

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("initial client not found, please configure the provider", "")
		return
	}

	var state RancherDevResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, state.Id.ValueString()),
		Method:   "DELETE",
	}

	response := c.Response{}

	err = client.Do(&request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting dev resource: ", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete Response Object: %+v", pp.PrettyPrint(*resp)))
}

func (r *RancherDevResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func validateData(data *RancherDevResourceModel) error {
	if data.Id.ValueString() == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if data.NumberAttribute.IsNull() {
		return fmt.Errorf("number_attribute cannot be empty")
	}
	if data.StringAttribute.ValueString() == "" {
		return fmt.Errorf("string_attribute cannot be empty")
	}
	return nil
}
