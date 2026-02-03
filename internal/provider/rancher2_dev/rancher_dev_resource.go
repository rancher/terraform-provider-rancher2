package rancher2_dev

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	Id               types.String                   `tfsdk:"id"`
	BoolAttribute    types.Bool                     `tfsdk:"bool_attribute"`
	NumberAttribute  types.Number                   `tfsdk:"number_attribute"`
	Int64Attribute   types.Int64                    `tfsdk:"int64_attribute"`
	Int32Attribute   types.Int32                    `tfsdk:"int32_attribute"`
	Float64Attribute types.Float64                  `tfsdk:"float64_attribute"`
	Float32Attribute types.Float32                  `tfsdk:"float32_attribute"`
	StringAttribute  types.String                   `tfsdk:"string_attribute"`
	ListAttribute    types.List                     `tfsdk:"list_attribute"`
	SetAttribute     types.Set                      `tfsdk:"set_attribute"`
	MapAttribute     types.Map                      `tfsdk:"map_attribute"`
	NestedObject     NestedResourceModel            `tfsdk:"nested_object"`
	NestedObjectList []NestedResourceModel          `tfsdk:"nested_object_list"`
	NestedObjectMap  map[string]NestedResourceModel `tfsdk:"nested_object_map"`
}
type NestedResourceModel struct {
	StringAttribute    types.String              `tfsdk:"string_attribute"`
	NestedNestedObject NestedNestedResourceModel `tfsdk:"nested_nested_object"`
}
type NestedNestedResourceModel struct {
	StringAttribute types.String `tfsdk:"string_attribute"`
	BoolAttribute   types.Bool   `tfsdk:"bool_attribute"`
}

// Fills the target with the values set in the current NestedResourceModel
func (r *NestedResourceModel) ToGoModel(ctx context.Context, target *NestedObject) diag.Diagnostics {
	dgs := diag.Diagnostics{}

	if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}
	if r.StringAttribute.IsNull() || r.StringAttribute.IsUnknown() {
		target.StringAttribute = ""
	}
	target.StringAttribute = r.StringAttribute.ValueString()

	diags := r.NestedNestedObject.ToGoModel(ctx, &target.NestedNestedObject)
	dgs.Append(diags...)
	return dgs
}

// Fills the target with the values set in the current NestedNestedResourceModel
func (r *NestedNestedResourceModel) ToGoModel(ctx context.Context, target *NestedNestedObject) diag.Diagnostics {
	dgs := diag.Diagnostics{}

	if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}
	if r.StringAttribute.IsNull() || r.StringAttribute.IsUnknown() {
		target.StringAttribute = ""
	}
	target.StringAttribute = r.StringAttribute.ValueString()
	if r.BoolAttribute.IsNull() || r.BoolAttribute.IsUnknown() {
		target.BoolAttribute = false
	}
	target.BoolAttribute = r.BoolAttribute.ValueBool()

	return dgs
}

type RancherDevModel struct {
	Id               string                  `json:"id"`
	BoolAttribute    bool                    `json:"bool_attribute,omitempty"`
	NumberAttribute  *big.Float              `json:"number_attribute,omitempty"`
	Int64Attribute   int64                   `json:"int64_attribute,omitempty"`
	Int32Attribute   int32                   `json:"int32_attribute,omitempty"`
	Float64Attribute float64                 `json:"float64_attribute,omitempty"`
	Float32Attribute float32                 `json:"float32_attribute,omitempty"`
	StringAttribute  string                  `json:"string_attribute,omitempty"`
	ListAttribute    []string                `json:"list_attribute,omitempty"`
	SetAttribute     map[string]bool         `json:"set_attribute,omitempty"`
	MapAttribute     map[string]string       `json:"map_attribute,omitempty"`
	NestedObject     NestedObject            `json:"nested_object"`
	NestedObjectList []NestedObject          `json:"nested_object_list,omitempty"`
	NestedObjectMap  map[string]NestedObject `json:"nested_object_map,omitempty"`
}
type NestedObject struct {
	StringAttribute    string             `json:"string_attribute,omitempty"`
	NestedNestedObject NestedNestedObject `json:"nested_nested_object"`
}
type NestedNestedObject struct {
	StringAttribute string `json:"string_attribute,omitempty"`
	BoolAttribute   bool   `json:"bool_attribute,omitempty"`
}

// Fills the target with types appropriate for a resource model.
func (m *NestedObject) ToResourceModel(ctx context.Context, target *NestedResourceModel) diag.Diagnostics {
	dgs := diag.Diagnostics{}

	target.StringAttribute = types.StringValue(m.StringAttribute)
	diags := m.NestedNestedObject.ToResourceModel(ctx, &target.NestedNestedObject)
	dgs.Append(diags...)
	return dgs
}

// Fills the target with types appropriate for a resource model.
func (m *NestedNestedObject) ToResourceModel(ctx context.Context, target *NestedNestedResourceModel) diag.Diagnostics {
	dgs := diag.Diagnostics{}
	target.StringAttribute = types.StringValue(m.StringAttribute)
	target.BoolAttribute = types.BoolValue(m.BoolAttribute)
	return dgs
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

	request := c.Request{
		Endpoint: fmt.Sprintf("%s/%s", client.GetApiUrl(), endpointPath),
		Method:   "POST",
		Body:     plan.ToGoModel(ctx),
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

	state := model.ToResource(ctx, &resp.Diagnostics)
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

	state = *respBody.ToResource(ctx, &resp.Diagnostics)
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

	state := *respBody.ToResource(ctx, &resp.Diagnostics)
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

// Conversion Functions

// ToGoModel converts a RancherDevResourceModel to a RancherDevModel
//
// This is helpful when building a json representation of the model for API requests.
func (data *RancherDevResourceModel) ToGoModel(ctx context.Context) RancherDevModel {
	obj := RancherDevModel{
		Id:               data.Id.ValueString(),
		BoolAttribute:    data.BoolAttribute.ValueBool(),
		Int32Attribute:   data.Int32Attribute.ValueInt32(),
		Int64Attribute:   data.Int64Attribute.ValueInt64(),
		Float64Attribute: data.Float64Attribute.ValueFloat64(),
		Float32Attribute: data.Float32Attribute.ValueFloat32(),
		StringAttribute:  data.StringAttribute.ValueString(),
		NumberAttribute:  data.NumberAttribute.ValueBigFloat(),
	}

	if len(data.MapAttribute.Elements()) > 0 {
		obj.MapAttribute = map[string]string{}
		data.MapAttribute.ElementsAs(ctx, &obj.MapAttribute, false)
	}

	if len(data.ListAttribute.Elements()) > 0 {
		obj.ListAttribute = []string{}
		data.ListAttribute.ElementsAs(ctx, &obj.ListAttribute, false)
	}

	var setAttrs []string
	if len(data.SetAttribute.Elements()) > 0 {
		obj.SetAttribute = map[string]bool{}
		data.SetAttribute.ElementsAs(ctx, &setAttrs, false)
		for _, v := range setAttrs {
			obj.SetAttribute[v] = true
		}
	}

	data.NestedObject.ToGoModel(ctx, &obj.NestedObject)

	if len(data.NestedObjectList) > 0 {
		obj.NestedObjectList = []NestedObject{}
		for _, nestedObjectFromList := range data.NestedObjectList {
			m := NestedObject{}
			nestedObjectFromList.ToGoModel(ctx, &m)
			obj.NestedObjectList = append(obj.NestedObjectList, m)
		}
	}

	if len(data.NestedObjectMap) > 0 {
		obj.NestedObjectMap = map[string]NestedObject{}
		for k, v := range data.NestedObjectMap {
			m := NestedObject{}
			v.ToGoModel(ctx, &m)
			obj.NestedObjectMap[k] = m
		}
	}

	return obj
}

// ToResource converts a RancherDevModel to a RancherDevResourceModel
//
// This is useful for processing json marshalled response bodies.
func (obj *RancherDevModel) ToResource(ctx context.Context, diags *diag.Diagnostics) *RancherDevResourceModel {
	if diags.HasError() {
		return nil
	}
	// var err error
	var data RancherDevResourceModel

	// primitive types (string, bool, int, etc)
	data.Id = types.StringValue(obj.Id)
	data.BoolAttribute = types.BoolValue(obj.BoolAttribute)
	data.Int32Attribute = types.Int32Value(obj.Int32Attribute)
	data.Int64Attribute = types.Int64Value(obj.Int64Attribute)
	data.Float64Attribute = types.Float64Value(obj.Float64Attribute)
	data.Float32Attribute = types.Float32Value(obj.Float32Attribute)
	data.StringAttribute = types.StringValue(obj.StringAttribute)
	data.NumberAttribute = types.NumberValue(obj.NumberAttribute)
	if diags.HasError() {
		return &data
	}

	// simple types (map, list, set)
	mapElems := make(map[string]attr.Value)
	if obj.MapAttribute != nil {
		for k, v := range obj.MapAttribute {
			mapElems[k] = basetypes.NewStringValue(v)
		}
	}
	mapVal, d := basetypes.NewMapValue(types.StringType, mapElems)
	diags.Append(d...)
	if mapVal.IsNull() {
		diags.AddError("Map Creation Error", "basetypes.NewMapValue returned null")
	}
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting map from value.")
		return &data
	}
	data.MapAttribute = mapVal

	var listElems []attr.Value
	for _, v := range obj.ListAttribute {
		listElems = append(listElems, basetypes.NewStringValue(v))
	}
	listVal, d := basetypes.NewListValue(types.StringType, listElems)
	diags.Append(d...)
	if listVal.IsNull() {
		diags.AddError("List Creation Error", "basetypes.NewListValue returned null")
	}
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting list from value.")
		return &data
	}
	data.ListAttribute = listVal

	var setAttributeElems []attr.Value
	for k := range obj.SetAttribute {
		setAttributeElems = append(setAttributeElems, basetypes.NewStringValue(k))
	}
	setVal, d := basetypes.NewSetValue(types.StringType, setAttributeElems)
	diags.Append(d...)
	if setVal.IsNull() {
		diags.AddError("Set Creation Error", "basetypes.NewSetValue returned null")
	}
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting set from value.")
		return &data
	}
	data.SetAttribute = setVal

	// complex types (nested objects)
	data.NestedObject = NestedResourceModel{}
	d = obj.NestedObject.ToResourceModel(ctx, &data.NestedObject)
	diags.Append(d...)
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting resource from nested object.")
		return &data
	}

	if len(obj.NestedObjectList) > 0 {
		for _, v := range obj.NestedObjectList {
			r := NestedResourceModel{}
			d = v.ToResourceModel(ctx, &r)
			diags.Append(d...)
			if diags.HasError() {
				diags.AddWarning("Progress:", "Error getting resource from list of nested objects.")
				return &data
			}
			data.NestedObjectList = append(data.NestedObjectList, r)
		}
	}

	if len(obj.NestedObjectMap) > 0 {
		data.NestedObjectMap = make(map[string]NestedResourceModel)

		for k, v := range obj.NestedObjectMap {
			r := NestedResourceModel{}
			d = v.ToResourceModel(ctx, &r)
			diags.Append(d...)
			if diags.HasError() {
				diags.AddWarning("Progress:", "Error getting resource from map of nested objects.")
				return &data
			}
			data.NestedObjectMap[k] = r
		}
	}

	err := validateData(&data)
	if err != nil {
		diags.AddError("Error validating data: ", err.Error())
	}

	return &data
}

// ToState returns a tfsdk.State with the data from the model.
func (m *RancherDevResourceModel) ToState(ctx context.Context, diags *diag.Diagnostics) tfsdk.State {
	if diags.HasError() {
		return tfsdk.State{}
	}
	diags.AddWarning("Model given ToState:", pp.PrettyPrint(m))
	r := NewRancherDevResource()
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

// ToPlan returns a tfsdk.Plan with the data from the model.
func (m *RancherDevResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
	if diags.HasError() {
		return tfsdk.Plan{}
	}
	// diags.AddWarning("Model given ToPlan:", fmt.Sprintf("%#v", m))
	r := NewRancherDevResource()
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
