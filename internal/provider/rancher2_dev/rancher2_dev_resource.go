package rancher2_dev

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	resp.TypeName = req.ProviderTypeName + "_dev" // rancher2_dev
}

func (r *RancherDevResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rancher Development resource. \n" +
			"This resource is used as a dummy for development purposes.",
		Attributes: map[string]schema.Attribute{
			// "id" is a special attribute in Terraform, it must be Computed (read only).
			// It must be represented as 'ID' in the model.
			// Terraform uses this attribute as a reference for other resources, it must be universally unique.
			// It is possible for this to be an attribute returned from the remote API, but it must be universally unique.
			// In most cases I recommend automatically setting this with a UUID in the provider (as done here).
			// If you have an id in your API that needs to be addressed, maybe call it "identifier" or "name".
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // when an attribute is set on creation and shouldn't change use this (UseStateForUnknown)
				},
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The id field returned from the API, used in GET, PUT, and DELETE requests.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // when an attribute is set on creation and shouldn't change use this (UseStateForUnknown)
				},
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
			"number_attribute": schema.NumberAttribute{
				MarkdownDescription: "A number attribute. Behind the scenes this is a big.Float.",
				Required:            true,
			},
			"bool_attribute": schema.BoolAttribute{
				MarkdownDescription: "A boolean attribute.",
				// Computed tells Terraform that this attribute shouldn't be empty in the state.
				// WARNING! If the function doesn't set this attribute in state there will be perpetual diffs.
				// Optional with Computed tells Terraform that the API or the provider always set this, even if it isn't in the config.
				// All attributes with a Default should be Optional and Computed.
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"int64_attribute": schema.Int64Attribute{
				MarkdownDescription: "A big int attribute.",
				// Optional without Computed means that this attribute may not exist in state, and that is OK
				Optional: true,
			},
			"int32_attribute": schema.Int32Attribute{
				MarkdownDescription: "A small int attribute.",
				Computed:            true, // This attribute is read only
			},
			"float64_attribute": schema.Float64Attribute{
				MarkdownDescription: "A float attribute.",
				Optional:            true,
			},
			"float32_attribute": schema.Float32Attribute{
				MarkdownDescription: "A small float attribute.",
				Optional:            true,
			},
			"map_attribute": schema.MapAttribute{
				MarkdownDescription: "A map of strings.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"list_attribute": schema.ListAttribute{
				MarkdownDescription: "A list of strings.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"set_attribute": schema.SetAttribute{
				MarkdownDescription: "A set of strings.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			// Don't use "object" type, use "single nested attribute" type instead,
			// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework@v1.19.0/resource/schema#SingleNestedAttribute
			// Nested Blocks are deprecated, use nested object attributes instead
			"nested_object": schema.SingleNestedAttribute{
				MarkdownDescription: "A single nested object." +
					"This represents a single object, not a list or set of objects.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"string_attribute": schema.StringAttribute{
						MarkdownDescription: "A string attribute.",
						Required:            true,
					},
					"nested_nested_object": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"string_attribute": schema.StringAttribute{
								MarkdownDescription: "A string attribute of an object within an object within a list of objects.",
								Required:            true,
							},
							"bool_attribute": schema.BoolAttribute{
								MarkdownDescription: "A read only boolean attribute of an object within an object within a list of objects.",
								Computed:            true,
							},
						},
					},
				},
			},
			// Don't use SetNestedAttribute, this isn't representable in json,
			// make the schema as close to the API representation as possible.
			// If possible avoid lists since it can be the source of dependency chain issues for users.
			"nested_object_list": schema.ListNestedAttribute{
				MarkdownDescription: "A list of nested objects.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attribute": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "A string attribute of an object within a list of objects.",
						},
						"nested_nested_object": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"string_attribute": schema.StringAttribute{
									MarkdownDescription: "A string attribute of an object within an object within a list of objects.",
									Required:            true,
								},
								"bool_attribute": schema.BoolAttribute{
									MarkdownDescription: "A read only boolean attribute of an object within an object within a list of objects.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"nested_object_map": schema.MapNestedAttribute{
				MarkdownDescription: "A map of nested objects.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attribute": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "A string attribute of an object within a set of objects.",
						},
						"nested_nested_object": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"string_attribute": schema.StringAttribute{
									MarkdownDescription: "A string attribute of an object within an object within a list of objects.",
									Required:            true,
								},
								"bool_attribute": schema.BoolAttribute{
									MarkdownDescription: "A read only boolean attribute of an object within an object within a list of objects.",
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

// configure runs at compile time, don't overload the context.
func (r *RancherDevResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Rancher2_Dev_Resource Configure request: %+v\n", pp.PrettyPrint(req)))
	tflog.Debug(ctx, fmt.Sprintf("Rancher2_Dev_Resource Configure object: %+v\n", pp.PrettyPrint(r)))
	if req.ProviderData == nil {
		return // Prevent panic if the provider has not been configured.
	}

	// Retrieving the client from the provider
	client, ok := req.ProviderData.(c.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
	tflog.Debug(ctx, fmt.Sprintf("Rancher2_Dev_Resource Configure object after config: %+v\n", pp.PrettyPrint(r)))
	tflog.Debug(ctx, fmt.Sprintf("Rancher2_Dev_Resource Configure response: %+v\n", pp.PrettyPrint(resp)))
}

// Create generates reality and state to match plan.
func (r *RancherDevResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Create Config: %+v\n", pp.PrettyPrint(req.Config.Raw)))
	tflog.Debug(ctx, fmt.Sprintf("Create Plan: %+v\n", pp.PrettyPrint(req.Plan.Raw)))
	var err error
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

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}
	endpoint := fmt.Sprintf("%s/%s", client.GetApiUrl(), endpointPath)

	// DEV ONLY
	//
	// This is special to the dev resource and
	// is only necessary due to the nature of the dev resource being a dummy resource for dev purposes.
	// When using the dev resource as a template for other resources, remove this.
	//
	var clnt c.TestClient
	rc, ok := client.(*c.HttpClient)
	if ok {
		// found a real client, need to inject a test client
		clnt = *c.NewTestClient(ctx, r.client.GetApiUrl(), "", false, false, 30, 10, rc.TokenStore)
		requestId := fmt.Sprintf("%s:%s:%s", endpoint, "POST", "")
		tflog.Debug(ctx, fmt.Sprintf("create requestId: %s", requestId))
		rbp := plan // copy the plan so it can still be used to generate the request
		rgm := rbp.ToGoModel(ctx)
		rgm.Int32Attribute = int32(1)        // simulating the API returning an attribute that isn't set by the provider
		rgm.Identifier = uuid.New().String() // simulating the API setting the id field
		responseBody, err := json.Marshal(rgm)
		if err != nil {
			resp.Diagnostics.AddError("Error marshalling dev plan for create response: ", err.Error())
			return
		}
		clnt.SetResponse(ctx, requestId, c.Response{
			StatusCode: 200,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: responseBody,
		})
		client = &clnt
	}
	//
	// END DEV ONLY

	var id types.String
	if plan.ID.IsNull() || plan.ID.IsUnknown() || plan.ID.ValueString() == "" {
		// don't blindly set uuid so that unit tests can set to a known id
		id = types.StringValue(uuid.New().String())
	} else {
		id = plan.ID
	}

	// All read only attributes should be set to their null value before sending the request.
	plan.ID = types.StringNull()
	plan.Identifier = types.StringNull()
	plan.Int32Attribute = types.Int32Null()

	tflog.Debug(ctx, fmt.Sprintf("Create Client: %+v", pp.PrettyPrint(client)))

	request := c.Request{
		Endpoint: endpoint,
		Method:   "POST",
		Body:     plan.ToGoModel(ctx),
	}

	response := c.Response{}

	err = client.Do(ctx, &request, &response)
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

	// Insert provider generated values before the state is saved.
	state.ID = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Create State After Set: %+v", pp.PrettyPrint(resp.State.Raw)))
}

// Read updates state to match reality.
// Read runs at refresh time which happens before all other functions and every time another function would be called.
// Don't call this function from one of the other functions (eg. don't call the Read function from within the Create function).
func (r *RancherDevResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Read Request: %+v", pp.PrettyPrint(req.State.Raw)))

	var state RancherDevResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var client c.Client

	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}

	var id types.String
	if state.ID.IsNull() || state.ID.IsUnknown() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"State Corruption",
			"The ID field in the state for this resource is null or unknown, "+
				"this shouldn't be possible, and is an indicator that there is a problem with your state file.",
		)
		return
	} else {
		id = state.ID
	}

	endpoint := fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, state.Identifier.ValueString())

	// DEV ONLY
	//
	// This is special to the dev resource and
	// is only necessary due to the nature of the dev resource being a dummy resource for dev purposes.
	// When using the dev resource as a template for other resources, remove this.
	//
	var clnt c.TestClient
	rc, ok := client.(*c.HttpClient)
	if ok {
		// found a real client, need to inject a test client
		clnt = *c.NewTestClient(ctx, r.client.GetApiUrl(), "", false, false, 30, 10, rc.TokenStore)
		requestId := fmt.Sprintf("%s:%s:%s", endpoint, "GET", "")
		tflog.Debug(ctx, fmt.Sprintf("read requestId: %s", requestId))
		responseBody, err := json.Marshal(state)
		if err != nil {
			resp.Diagnostics.AddError("Error marshalling dev plan for create response: ", err.Error())
			return
		}
		clnt.SetResponse(ctx, requestId, c.Response{
			StatusCode: 200,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: responseBody,
		})
		client = &clnt
	}
	//
	// END DEV ONLY

	request := c.Request{
		Endpoint: endpoint,
		Method:   "GET",
	}
	response := c.Response{}

	err := client.Do(ctx, &request, &response)
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

	// Inject provider only (non-API) attributes into the state before saving
	state.ID = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Read State After Set: %+v", pp.PrettyPrint(resp.State.Raw)))
}

// Update changes reality and state to match plan (best practice is don't compare old state, just override).
func (r *RancherDevResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Update Request Config: %+v", pp.PrettyPrint(req.Config.Raw)))
	tflog.Debug(ctx, fmt.Sprintf("Update Request Plan: %+v", pp.PrettyPrint(req.Plan.Raw)))
	tflog.Debug(ctx, fmt.Sprintf("Update Request State: %+v", pp.PrettyPrint(req.State.Raw)))
	var err error
	var plan RancherDevResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state RancherDevResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
	if state.ID.IsNull() || state.ID.IsUnknown() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"State Corruption",
			"The ID field in the state for this resource is null or unknown, "+
				"this shouldn't be possible, and is an indicator that there is a problem with your state file.",
		)
		return
	} else {
		id = state.ID
	}

	err = validateData(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Error validating plan: ", err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Update Plan after validate: %+v", pp.PrettyPrint(plan)))

	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}

	endpoint := fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, plan.ID.ValueString())

	// DEV ONLY
	//
	// This is special to the dev resource and
	// is only necessary due to the nature of the dev resource being a dummy resource for dev purposes.
	// When using the dev resource as a template for other resources, remove this.
	//
	var clnt c.TestClient
	rc, ok := client.(*c.HttpClient)
	if ok {
		// found a real client, need to inject a test client
		clnt = *c.NewTestClient(ctx, r.client.GetApiUrl(), "", false, false, 30, 10, rc.TokenStore)
		requestId := fmt.Sprintf("%s:%s:%s", endpoint, "PUT", "")
		tflog.Debug(ctx, fmt.Sprintf("update requestId: %s", requestId))
		rbp := plan // copy the plan so it can still be used to generate the request
		rgm := rbp.ToGoModel(ctx)
		rgm.Int32Attribute = int32(1) // simulating the API returning an attribute that isn't set by the provider
		responseBody, err := json.Marshal(rgm)
		if err != nil {
			resp.Diagnostics.AddError("Error marshalling dev plan for create response: ", err.Error())
			return
		}
		clnt.SetResponse(ctx, requestId, c.Response{
			StatusCode: 200,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: responseBody,
		})
		client = &clnt
	}
	//
	// END DEV ONLY

	// Remove any attributes that shouldn't be sent in the request body (usually all read only attributes)
	plan.ID = types.StringNull()            // already pulled this from state
	plan.Identifier = types.StringNull()    // use the response object's value
	plan.Int32Attribute = types.Int32Null() // use the response object's value

	request := c.Request{
		Endpoint: endpoint,
		Method:   "PUT",
		Body:     plan.ToGoModel(ctx),
	}

	response := c.Response{}
	err = client.Do(ctx, &request, &response)
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

	state = *respBody.ToResourceModel(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Inject any provider only (non-API) attributes into the state before saving.
	state.ID = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Update Response State After Set: %+v", pp.PrettyPrint(resp.State.Raw)))
}

// Destroy destroys reality (state is handled automatically).
func (r *RancherDevResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Delete Request State: %+v", pp.PrettyPrint(req.State.Raw)))
	var err error

	var state RancherDevResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var client c.Client
	if r.client != nil {
		client = r.client
	} else {
		// no client found, seems like the provider wasn't configured properly
		resp.Diagnostics.AddError("client not found, please configure the provider", "")
		return
	}

	endpoint := fmt.Sprintf("%s/%s/%s", client.GetApiUrl(), endpointPath, state.ID.ValueString())

	// DEV ONLY
	//
	// This is special to the dev resource and
	// is only necessary due to the nature of the dev resource being a dummy resource for dev purposes.
	// When using the dev resource as a template for other resources, remove this.
	//
	var clnt c.TestClient
	rc, ok := client.(*c.HttpClient)
	if ok {
		// found a real client, need to inject a test client
		clnt = *c.NewTestClient(ctx, r.client.GetApiUrl(), "", false, false, 30, 10, rc.TokenStore)
		requestId := fmt.Sprintf("%s:%s:%s", endpoint, "DELETE", "")
		tflog.Debug(ctx, fmt.Sprintf("delete requestId: %s", requestId))
		clnt.SetResponse(ctx, requestId, c.Response{StatusCode: 200})
		client = &clnt
	}
	//
	// END DEV ONLY

	request := c.Request{
		Endpoint: endpoint,
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
}

func (r *RancherDevResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// This function also enforces default values.
func validateData(data *RancherDevResourceModel) error {
	if data.StringAttribute.IsNull() || data.StringAttribute.IsUnknown() || data.StringAttribute.ValueString() == "" {
		return fmt.Errorf("string_attribute cannot be empty")
	}
	if !regexp.MustCompile(`^dev-.*`).MatchString(data.StringAttribute.ValueString()) {
		return fmt.Errorf("string must start with 'dev-'")
	}
	if data.NumberAttribute.IsNull() || data.NumberAttribute.IsUnknown() {
		return fmt.Errorf("number_attribute cannot be empty")
	}
	if data.BoolAttribute.IsNull() || data.BoolAttribute.IsUnknown() {
		data.BoolAttribute = types.BoolValue(true)
	}
	return nil
}
