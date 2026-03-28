package rancher2_metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MetadataAttribute defines the generic Kubernetes metadata schema as a Terraform schema.
func MetadataAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "The metadata of the resource.",
		Required:            true,
		Validators: []validator.Object{
			nameAndGenerateNameValidator{},
		},
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the resource.",
				Optional:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace of the resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
			},
			"generate_name": schema.StringAttribute{
				MarkdownDescription: "The generate name of the resource.",
				Optional:            true,
			},
			"annotations": schema.MapAttribute{
				MarkdownDescription: "The annotations of the resource.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					metadataMergeModifier{},
				},
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: "The labels of the resource.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					metadataMergeModifier{},
				},
			},
			"finalizers": schema.ListAttribute{
				MarkdownDescription: "Advanced use cases only. The finalizers of the resource.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					finalizerModifier{},
				},
			},
			"owner_references": schema.ListNestedAttribute{
				MarkdownDescription: "Advanced use cases only. The owner references of the resource.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_version":          schema.StringAttribute{Optional: true},
						"kind":                 schema.StringAttribute{Optional: true},
						"name":                 schema.StringAttribute{Optional: true},
						"uid":                  schema.StringAttribute{Optional: true},
						"controller":           schema.BoolAttribute{Optional: true},
						"block_owner_deletion": schema.BoolAttribute{Optional: true},
					},
				},
				PlanModifiers: []planmodifier.List{
					ownerReferenceModifier{},
				},
			},
			"uid":                           schema.StringAttribute{Computed: true},
			"generation":                    schema.Int64Attribute{Computed: true},
			"creation_timestamp":            schema.StringAttribute{Computed: true},
			"deletion_grace_period_seconds": schema.Int64Attribute{Computed: true},
			"deletion_timestamp":            schema.StringAttribute{Computed: true},
			"managed_fields":                schema.StringAttribute{Computed: true},
			"resource_version":              schema.StringAttribute{Computed: true},
			"self_link":                     schema.StringAttribute{Computed: true},
		},
	}
}

// metadataMergeModifier implements the "Sticky" map modifier to ensure Terraform
// only reconciles what it "owns", and adopts cluster state for API-only keys.
type metadataMergeModifier struct{}

func (m metadataMergeModifier) Description(ctx context.Context) string {
	return "Merges state and config maps, keeping unconfigured state keys (adopted external changes)."
}

func (m metadataMergeModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m metadataMergeModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.ConfigValue.IsNull() {
		// if the user doesn't specify a config, set plan = state
		resp.PlanValue = req.StateValue
		return
	}

	var configMap, stateMap map[string]types.String
	req.ConfigValue.ElementsAs(ctx, &configMap, false)
	req.StateValue.ElementsAs(ctx, &stateMap, false)

	plannedMap := make(map[string]types.String)
	for k, v := range configMap {
		plannedMap[k] = v // If key exists in config, save value in plan (save user intent)
	}

	for k, v := range stateMap {
		if _, managedByTF := configMap[k]; !managedByTF {
			plannedMap[k] = v // If key isn't in config, add value to plan (adopt unplanned API keys)
		}
	}
	resp.PlanValue, _ = types.MapValueFrom(ctx, types.StringType, plannedMap)
}

// finalizerModifier implements the "Big Red Button" list modifier.
// Merge by default. If HCL is an explicit empty list [], overwrite State to clear blockers.
type finalizerModifier struct{}

func (m finalizerModifier) Description(ctx context.Context) string {
	return "Merges finalizers. If the user explicitly sets [], clears all finalizers."
}

func (m finalizerModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m finalizerModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If user explicitly provides an empty list, bypass merging to clear finalizers
	if !req.ConfigValue.IsNull() && len(req.ConfigValue.Elements()) == 0 {
		resp.PlanValue = req.ConfigValue
		return
	}

	// Otherwise, perform a Union merge of Config and State
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.ConfigValue.IsNull() {
		resp.PlanValue = req.StateValue
		return
	}

	var configList, stateList []types.String
	req.ConfigValue.ElementsAs(ctx, &configList, false)
	req.StateValue.ElementsAs(ctx, &stateList, false)

	plannedSet := make(map[string]struct{})
	var plannedList []types.String

	for _, v := range configList {
		plannedList = append(plannedList, v)
		plannedSet[v.ValueString()] = struct{}{}
	}
	for _, v := range stateList {
		if _, exists := plannedSet[v.ValueString()]; !exists {
			plannedList = append(plannedList, v)
			plannedSet[v.ValueString()] = struct{}{}
		}
	}
	resp.PlanValue, _ = types.ListValueFrom(ctx, types.StringType, plannedList)
}

// ownerReferenceModifier matches by uid. Merges API-side refs into the Plan; does not delete them.
type ownerReferenceModifier struct{}

func (m ownerReferenceModifier) Description(ctx context.Context) string {
	return "Merges OwnerReferences, matching by UID to prevent orphaning cluster-side attachments."
}

func (m ownerReferenceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m ownerReferenceModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.ConfigValue.IsNull() {
		resp.PlanValue = req.StateValue
		return
	}

	var configList, stateList []types.Object
	req.ConfigValue.ElementsAs(ctx, &configList, false)
	req.StateValue.ElementsAs(ctx, &stateList, false)

	plannedList := append([]types.Object{}, configList...)

	for _, stateObj := range stateList {
		stateAttrs := stateObj.Attributes()
		stateUidAttr, hasUid := stateAttrs["uid"]
		if !hasUid {
			continue
		}
		var stateUid string
		if stateUidValue, ok := stateUidAttr.(types.String); ok {
			stateUid = stateUidValue.ValueString()
		} else {
			continue
		}

		found := false
		for _, configObj := range configList {
			configAttrs := configObj.Attributes()
			if configUidAttr, ok := configAttrs["uid"]; ok {
				var configUid string
				if configUidValue, ok := configUidAttr.(types.String); ok {
					configUid = configUidValue.ValueString()
				}
				if stateUid == configUid {
					found = true
					break
				}
			}
		}
		if !found {
			plannedList = append(plannedList, stateObj)
		}
	}

	resp.PlanValue, _ = types.ListValueFrom(ctx, req.StateValue.ElementType(ctx), plannedList)
}

type nameAndGenerateNameValidator struct{}

func (v nameAndGenerateNameValidator) Description(ctx context.Context) string {
	return "Validates that name and generate_name are mutually exclusive and that one is set."
}

func (v nameAndGenerateNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v nameAndGenerateNameValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	var name, generateName types.String
	dgs := req.Config.GetAttribute(ctx, path.Root("metadata").AtName("name"), &name)
	if dgs.HasError() {
		resp.Diagnostics.Append(dgs...)
		return
	}
	dgs = req.Config.GetAttribute(ctx, path.Root("metadata").AtName("generate_name"), &generateName)
	if dgs.HasError() {
		resp.Diagnostics.Append(dgs...)
		return
	}

	isNameSet := !name.IsNull() && !name.IsUnknown()
	isGenerateNameSet := !generateName.IsNull() && !generateName.IsUnknown()

	if isNameSet && isGenerateNameSet {
		resp.Diagnostics.AddError(
			"mutually_exclusive_attributes",
			"name and generate_name are mutually exclusive, but both were set.",
		)
	}

	if !isNameSet && !isGenerateNameSet {
		resp.Diagnostics.AddError(
			"one_of_required",
			"one of name or generate_name must be set, but neither were.",
		)
	}
}
