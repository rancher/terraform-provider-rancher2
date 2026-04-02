package rancher2_dev2

// Rancher2 Dev2 Resource is an example Terraform resource that represents a Kubernetes CRD.
// This is a dummy resource, it functions as a real resource, but it doesn't make API calls.
// Use this as a template for new resources and in testing.

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	"github.com/rancher/terraform-provider-rancher2/internal/provider/rancher2_metadata"
)

// The "var _" is a special Go construct that results in an unusable variable.
// The purpose of these lines is to make sure our LocalFileResource correctly implements the "resource.Resource“ interface.
// These will fail at compilation time if the implementation is not satisfied.
var _ resource.Resource = &Rancher2Dev2Resource{}
var _ resource.ResourceWithImportState = &Rancher2Dev2Resource{}

const (
	endpointPath = "dev2"
)

func NewRancher2Dev2Resource() resource.Resource {
	return &Rancher2Dev2Resource{}
}

type Rancher2Dev2Resource struct {
	client c.Client // client is an interface holding a pointer to a struct
}

// Rancher2Dev2ResourceModel is in rancher2_dev2_resource_model.go

func (r *Rancher2Dev2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev2" // rancher2_dev2
}

func (r *Rancher2Dev2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: "Rancher Development 2 resource. \n" +
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
				Computed:            true, // Computed without Optional/Required tells Terraform that this attribute is read-only.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // UseStateForUnknown tells Terraform that this attribute should only be set once when the object is created.
					stringplanmodifier.RequiresReplace(),    // RequiresReplace tells Terraform that if this value changes the resource needs to be recreated.
				},
			},
			"api_version": schema.StringAttribute{
				MarkdownDescription: "The API version of the resource.",
				Required:            true,
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "The kind of the resource.",
				Required:            true,
			},
			"metadata": rancher2_metadata.MetadataAttribute(),
			"spec": schema.SingleNestedAttribute{
				MarkdownDescription: "The specification of the object.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"string":  schema.StringAttribute{Optional: true},
					"bool":    schema.BoolAttribute{Optional: true},
					"number":  schema.NumberAttribute{Optional: true},
					"int32":   schema.Int32Attribute{Optional: true},
					"int64":   schema.Int64Attribute{Optional: true},
					"float32": schema.Float32Attribute{Optional: true},
					"float64": schema.Float64Attribute{Optional: true},
					"map": schema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"list": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"object": schema.SingleNestedAttribute{
						MarkdownDescription: "An object within the spec.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"string_attribute": schema.StringAttribute{
								MarkdownDescription: "A string attribute.",
								Optional:            true,
							},
						},
					},
					"object_list": schema.ListNestedAttribute{
						MarkdownDescription: "This is a list of objects within the spec.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"string_attribute": schema.StringAttribute{
									MarkdownDescription: "A string attribute.",
									Optional:            true,
								},
							},
						},
					},
					"object_map": schema.MapNestedAttribute{
						MarkdownDescription: "This is a map of objects within the spec.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"string_attribute": schema.StringAttribute{
									MarkdownDescription: "A string attribute.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the resource (JSON blob).",
				Computed:            true,
			},
			// Ignore this attribute when using this resource as a template.
			"api_responses": schema.MapNestedAttribute{
				MarkdownDescription: "Map of function to response, eg. create, read, update, delete.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"headers": schema.MapAttribute{ // map[string][]string
							Optional: true,
							ElementType: types.ListType{
								ElemType: types.StringType,
							},
						},
						"body": schema.StringAttribute{
							Optional: true,
						},
						"status_code": schema.Int64Attribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *Rancher2Dev2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Terraform may call this method before the provider has been configured (to make sure the resource is valid),
	//  so we need to gracefully return an empty resource.
	if req.ProviderData == nil {
		return
	}

	client, dgs := client(req.ProviderData)
	if dgs.HasError() {
		resp.Diagnostics.Append(*dgs...)
		return
	}
	r.client = client
}

// Create changes reality and state to match plan.
func (r *Rancher2Dev2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if resp.State.Schema == nil {
		resp.Diagnostics.AddError(
			"State Schema Not Found",
			"The state schema is missing from the response container.",
		)
		return
	}
  // Plan modifiers are the mechanism the provider uses to convert the config to the plan.
	// req.Config shouldn't be used in this function, req.Plan should convey user intent, if there is any confusion use Plan Modifiers.
	plan, dgs := plan(ctx, req.Plan)
	if dgs.HasError() {
		resp.Diagnostics.Append(*dgs...)
		return
	}
	client, dgs := client(r.client)
	if dgs.HasError() {
		resp.Diagnostics.Append(*dgs...)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Create Request Plan: %+v\n", req.Plan))

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
		planGoModel := plan.ToGoModel(ctx, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		clnt.SetResponse(ctx, requestId, c.Response{
			StatusCode: int(planGoModel.APIResponses["create"].StatusCode),
			Headers:    planGoModel.APIResponses["create"].Headers,
			Body:       json.RawMessage(planGoModel.APIResponses["create"].Body),
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

	rb := plan.ToGoModel(ctx, &resp.Diagnostics).ToApiRequestBody(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.Request{
		Endpoint: endpoint,
		Method:   "POST",
		Body:     rb,
	}

	response := c.Response{}

	err := client.Do(ctx, &request, &response)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dev resource: ", err.Error())
		return
	}
	// process the response here
	var model Rancher2Dev2Model
	model.FromAPIResponseBody(ctx, response.Body, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	state := model.ToResourceModel(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Insert provider generated values before the state is saved.
	state.ID = id

	// Insert user generated values that only exist in state before the state is saved.
	state.APIResponses = plan.APIResponses // API responses isn't something that exists in a remote API, it is only saved in the state.

	d := resp.State.Set(ctx, state)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
	}
	tflog.Debug(ctx, fmt.Sprintf("Create State After Set: %+v", pp.PrettyPrint(resp.State.Raw)))
}

// Read updates the state to match reality.
func (r *Rancher2Dev2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if resp.State.Schema == nil {
		resp.Diagnostics.AddError(
			"State Schema Not Found",
			"The state schema is missing from the response container.",
		)
		return
	}
  state, dgs := state(ctx, req.State)
	if dgs.HasError() {
		resp.Diagnostics.Append(*dgs...)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Read Request State: %+v\n", pp.PrettyPrint(state)))
	tflog.Debug(ctx, fmt.Sprintf("Read Request Provider Config: %+v\n", pp.PrettyPrint(client)))
}

// Update changes reality and state to match plan.
// Best practice is not to compare, just overwrite.
func (r *Rancher2Dev2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if resp.State.Schema == nil {
		resp.Diagnostics.AddError(
			"State Schema Not Found",
			"The state schema is missing from the response container.",
		)
		return
	}
  // req.Config shouldn't be used in this function, req.Plan should convey user intent, if there is any confusion use Plan Modifiers.
	plan, dgs := plan(ctx, req.Plan)
	if dgs.HasError() {
		resp.Diagnostics.Append(*dgs...)
		return
	}
	state, dgs := state(ctx, req.State)
	if dgs.HasError() {
		resp.Diagnostics.Append(*dgs...)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Update Request Plan: %+v\n", pp.PrettyPrint(plan)))
	tflog.Debug(ctx, fmt.Sprintf("Update Request State: %+v\n", pp.PrettyPrint(state)))
	tflog.Debug(ctx, fmt.Sprintf("Update Request Provider Config: %+v\n", pp.PrettyPrint(client)))
}

// Delete removes reality, the framework automatically handles the state.
func (r *Rancher2Dev2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
  state, dgs := state(ctx, req.State)
	if dgs.HasError() {
		resp.Diagnostics.Append(*dgs...)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Delete Request State: %+v\n", pp.PrettyPrint(state)))
	tflog.Debug(ctx, fmt.Sprintf("Delete Request Provider Config: %+v\n", pp.PrettyPrint(client)))
}

func (r *Rancher2Dev2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Import Error",
			"Import request missing ID.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Validators.
func state(ctx context.Context, state tfsdk.State) (Rancher2Dev2ResourceModel, *diag.Diagnostics) {
	var dgs diag.Diagnostics
	var model Rancher2Dev2ResourceModel

	emptyState := tfsdk.State{}
	if state == emptyState {
		dgs.AddError(
			"State not found",
			"The state is missing from the request.",
		)
		return model, &dgs
	}

	dgs.Append(state.Get(ctx, &model)...)
	if dgs.HasError() {
		return model, &dgs
	}
	tflog.Debug(ctx, fmt.Sprintf("State: %+v\n", pp.PrettyPrint(model)))
	return model, &dgs
}

func plan(ctx context.Context, plan tfsdk.Plan) (Rancher2Dev2ResourceModel, *diag.Diagnostics) {
	var dgs diag.Diagnostics
	var model Rancher2Dev2ResourceModel

	emptyPlan := tfsdk.Plan{}
	if plan == emptyPlan {
		dgs.AddError(
			"Plan not found",
			"The plan is missing from the request.",
		)
		return model, &dgs
	}

	dgs.Append(plan.Get(ctx, &model)...)
	if dgs.HasError() {
		return model, &dgs
	}

	tflog.Debug(ctx, fmt.Sprintf("Plan: %+v\n", pp.PrettyPrint(model)))
	return model, &dgs
}

func client(client any) (c.Client, *diag.Diagnostics) {
	var dgs diag.Diagnostics
	var clnt c.Client
	if client == nil {
		dgs.AddError(
			"Client not found",
			"The client is missing from the request, please configure the provider.",
		)
		return clnt, &dgs
	}

	clnt, ok := client.(c.Client)
	if !ok {
		dgs.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected c.Client, got: %T. Please report this issue to the provider developers.", clnt),
		)
		return clnt, &dgs
	}
	tflog.Debug(context.Background(), fmt.Sprintf("Client: %+v\n", pp.PrettyPrint(clnt)))
	return clnt, &dgs
}
