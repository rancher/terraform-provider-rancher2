package rancher2_dev

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

// RancherDevResourceModel is in rancher2_dev_resource_model.go

func (r *RancherDevResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev_resource" // rancher2_dev_resource
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_token": schema.StringAttribute{
				MarkdownDescription: "A token to inject into the client for any requests this resource needs to make.",
				Optional:            true,
				Computed:            true,
			},
			"bool_attribute": schema.BoolAttribute{
				MarkdownDescription: "A boolean attribute.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"number_attribute": schema.NumberAttribute{
				MarkdownDescription: "A number attribute. Behind the scenes this is a big.Float.",
				Required:            true,
			},
			"int64_attribute": schema.Int64Attribute{
				MarkdownDescription: "A big int attribute.",
				Optional:            true,
				Computed:            true,
			},
			"int32_attribute": schema.Int32Attribute{
				MarkdownDescription: "A small int attribute.",
				Optional:            true,
				Computed:            true,
			},
			"float64_attribute": schema.Float64Attribute{
				MarkdownDescription: "A float attribute.",
				Optional:            true,
				Computed:            true,
			},
			"float32_attribute": schema.Float32Attribute{
				MarkdownDescription: "A small float attribute.",
				Optional:            true,
				Computed:            true,
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
			"map_attribute": schema.MapAttribute{
				MarkdownDescription: "A map of strings.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"list_attribute": schema.ListAttribute{
				MarkdownDescription: "A list of strings.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"set_attribute": schema.SetAttribute{
				MarkdownDescription: "A set of strings.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			// Don't use "object" type, use "single nested attribute" type instead
			"nested_object": schema.SingleNestedAttribute{
				MarkdownDescription: "A single nested object." +
					"This represents a single object, not a list or set of objects.",
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"string_attribute": schema.StringAttribute{
						MarkdownDescription: "A string attribute.",
						Optional:            true,
						Computed:            true,
					},
					"nested_nested_object": schema.SingleNestedAttribute{
						MarkdownDescription: "A nested object within a nested object.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"string_attribute": schema.StringAttribute{
								MarkdownDescription: "A string attribute.",
								Optional:            true,
								Computed:            true,
							},
							"bool_attribute": schema.BoolAttribute{
								MarkdownDescription: "A boolean attribute.",
								Optional:            true,
								Computed:            true,
							},
						},
					},
				},
			},
			"nested_object_list": schema.ListNestedAttribute{
				MarkdownDescription: "A list of nested objects.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attribute": schema.StringAttribute{
							MarkdownDescription: "A string attribute.",
							Optional:            true,
							Computed:            true,
						},
						"nested_nested_object": schema.SingleNestedAttribute{
							MarkdownDescription: "A nested object within a nested object.",
							Optional:            true,
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"string_attribute": schema.StringAttribute{
									MarkdownDescription: "A string attribute.",
									Optional:            true,
									Computed:            true,
								},
								"bool_attribute": schema.BoolAttribute{
									MarkdownDescription: "A boolean attribute.",
									Optional:            true,
									Computed:            true,
								},
							},
						},
					},
				},
			},
			// Don't use nested object sets, use a map instead
			// You can generate an ordered map to use as a set with map[string]any{"1": any, "2": any}
			"nested_object_map": schema.MapNestedAttribute{
				MarkdownDescription: "A map of nested objects.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attribute": schema.StringAttribute{
							MarkdownDescription: "A string attribute.",
							Optional:            true,
							Computed:            true,
						},
						"nested_nested_object": schema.SingleNestedAttribute{
							MarkdownDescription: "A nested object within a nested object.",
							Optional:            true,
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"string_attribute": schema.StringAttribute{
									MarkdownDescription: "A string attribute.",
									Optional:            true,
									Computed:            true,
								},
								"bool_attribute": schema.BoolAttribute{
									MarkdownDescription: "A boolean attribute.",
									Optional:            true,
									Computed:            true,
								},
							},
						},
					},
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

// Create generates reality and state to match plan.
func (r *RancherDevResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Create Request Object: %+v", pp.PrettyPrint(req)))
	var err error

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}

	plan := RancherDevResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = validateData(&plan) // we validate data here because tests will bypass the schema validators
	if err != nil {
		resp.Diagnostics.AddError("Error validating data: ", err.Error())
		return
	}
	token := plan.UserToken.ValueString()
	plan.UserToken = types.StringValue("") // this implicitly removes the user_token from any request body because omitempty and jsonMarshal

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s", client.GetApiUrl(), endpointPath),
		Method:   "POST",
		Body:     plan.ToGoModel(ctx),
		Token:    token,
	}

	response := c.Response{}

	err = client.Do(&request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dev resource: ", err.Error())
		return
	}

	// process the response here
	var model RancherDevModel
	err = json.Unmarshal(response.Body, &model)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshaling response body:", err.Error())
		return
	}

	state := model.ToResourceModel(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Create Response Object: %+v", pp.PrettyPrint(*resp)))
}

// Read updates state to match reality.
// Read runs at refresh time which happens before all other functions and every time another function would be called.
// Don't call this function from one of the other functions (eg. don't call the Read function from within the Create function).
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
		if e, ok := err.(*c.ApiError); ok && e.StatusCode == 404 {
			// resource not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading dev resource: ", err.Error())
		return
	}

	var respBody RancherDevModel
	err = json.Unmarshal(response.Body, &respBody)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshalling dev resource: ", err.Error())
		return
	}

	state = *respBody.ToResourceModel(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Read Response Object: %+v", pp.PrettyPrint(*resp)))
}

// Update changes reality and state to match plan (best practice is don't compare old state, just override).
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

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, plan.Id.ValueString()),
		Method:   "PUT",
		Body:     plan.ToGoModel(ctx),
	}

	response := c.Response{}

	err = client.Do(&request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dev resource: ", err.Error())
		return
	}

	var respBody RancherDevModel
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
	tflog.Debug(ctx, fmt.Sprintf("Update Response Object: %+v", pp.PrettyPrint(*resp)))
}

// Destroy destroys reality (state is handled automatically).
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
		if e, ok := err.(*c.ApiError); ok && e.StatusCode == 404 {
			// resource already deleted
			return
		}
		resp.Diagnostics.AddError("Error deleting dev resource: ", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete Response Object: %+v", pp.PrettyPrint(*resp)))
}

func (r *RancherDevResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// note this function also enforces default values
func validateData(data *RancherDevResourceModel) error {
	if data.Id.ValueString() == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if regexp.MustCompile(`^dev-.*`).MatchString(data.Id.ValueString()) == false {
		return fmt.Errorf("Id must start with 'dev-'")
	}
	if data.BoolAttribute.IsNull() || data.BoolAttribute.IsUnknown() {
		data.BoolAttribute = types.BoolValue(true)
	}
	if data.MapAttribute.IsNull() || data.MapAttribute.IsUnknown() {
		defaultMap, diags := types.MapValue(types.MapType{ElemType: types.StringType}, map[string]attr.Value{})
		if diags.HasError() {
			return fmt.Errorf("error generating default map for partially tracked attribute")
		}
		data.MapAttribute = defaultMap
	}

	return nil
}
