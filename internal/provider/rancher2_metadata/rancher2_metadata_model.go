package rancher2_metadata

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type OwnerReference struct {
	APIVersion         string `tfsdk:"api_version" json:"api_version"`
	Kind               string `tfsdk:"kind" json:"kind"`
	Name               string `tfsdk:"name" json:"name"`
	UID                string `tfsdk:"uid" json:"uid"`
	Controller         bool   `tfsdk:"controller" json:"controller"`
	BlockOwnerDeletion bool   `tfsdk:"block_owner_deletion" json:"block_owner_deletion"`
}

var OwnerReferenceAttrTypes = map[string]attr.Type{
	"api_version":          types.StringType,
	"kind":                 types.StringType,
	"name":                 types.StringType,
	"uid":                  types.StringType,
	"controller":           types.BoolType,
	"block_owner_deletion": types.BoolType,
}

type Metadata struct {
	Name                      string            `tfsdk:"name" json:"name"`
	Namespace                 string            `tfsdk:"namespace" json:"namespace"`
	GenerateName              string            `tfsdk:"generate_name" json:"generate_name"`
	Annotations               map[string]string `tfsdk:"annotations" json:"annotations"`
	Labels                    map[string]string `tfsdk:"labels" json:"labels"`
	Finalizers                []string          `tfsdk:"finalizers" json:"finalizers"`
	OwnerReferences           []OwnerReference  `tfsdk:"owner_references" json:"owner_references"`
	UID                       string            `tfsdk:"uid" json:"uid"`
	Generation                int64             `tfsdk:"generation" json:"generation"`
	CreationTimestamp         string            `tfsdk:"creation_timestamp" json:"creation_timestamp"`
	DeletionGracePeriodSecond int64             `tfsdk:"deletion_grace_period_seconds" json:"deletion_grace_period_seconds"`
	DeletionTimestamp         string            `tfsdk:"deletion_timestamp" json:"deletion_timestamp"`
	ManagedFields             json.RawMessage   `tfsdk:"managed_fields" json:"managed_fields"`
	ResourceVersion           string            `tfsdk:"resource_version" json:"resource_version"`
	SelfLink                  string            `tfsdk:"self_link" json:"self_link"`
}

var MetadataAttrTypes = map[string]attr.Type{
	"name":                          types.StringType,
	"namespace":                     types.StringType,
	"generate_name":                 types.StringType,
	"annotations":                   types.MapType{ElemType: types.StringType},
	"labels":                        types.MapType{ElemType: types.StringType},
	"finalizers":                    types.ListType{ElemType: types.StringType},
	"owner_references":              types.ListType{ElemType: types.ObjectType{AttrTypes: OwnerReferenceAttrTypes}},
	"uid":                           types.StringType,
	"generation":                    types.Int64Type,
	"creation_timestamp":            types.StringType,
	"deletion_grace_period_seconds": types.Int64Type,
	"deletion_timestamp":            types.StringType,
	"managed_fields":                types.StringType,
	"resource_version":              types.StringType,
	"self_link":                     types.StringType,
}

func (m *Metadata) ToTypesObject(ctx context.Context, diags *diag.Diagnostics) types.Object {
	var ownerReferencesValue attr.Value = types.ListNull(types.ObjectType{AttrTypes: OwnerReferenceAttrTypes})
	if len(m.OwnerReferences) > 0 {
		var ownerReferences []attr.Value
		for _, or := range m.OwnerReferences {
			orObj := or.ToTypesObject(ctx, diags)
			if diags.HasError() {
				return types.ObjectNull(MetadataAttrTypes)
			}
			ownerReferences = append(ownerReferences, orObj)
		}
		var d diag.Diagnostics
		ownerReferencesValue, d = basetypes.NewListValue(types.ObjectType{AttrTypes: OwnerReferenceAttrTypes}, ownerReferences)
		diags.Append(d...)
	}

	var annotationsValue attr.Value = types.MapNull(types.StringType)
	if m.Annotations != nil {
		var annotations = make(map[string]attr.Value)
		for k, v := range m.Annotations {
			annotations[k] = basetypes.NewStringValue(v)
		}
		var d diag.Diagnostics
		annotationsValue, d = basetypes.NewMapValue(types.StringType, annotations)
		diags.Append(d...)
	}

	var labelsValue attr.Value = types.MapNull(types.StringType)
	if m.Labels != nil {
		var labels = make(map[string]attr.Value)
		for k, v := range m.Labels {
			labels[k] = basetypes.NewStringValue(v)
		}
		var d diag.Diagnostics
		labelsValue, d = basetypes.NewMapValue(types.StringType, labels)
		diags.Append(d...)
	}

	var finalizersValue attr.Value = types.ListNull(types.StringType)
	if len(m.Finalizers) > 0 {
		var finalizers []attr.Value
		for _, f := range m.Finalizers {
			finalizers = append(finalizers, basetypes.NewStringValue(f))
		}
		var d diag.Diagnostics
		finalizersValue, d = basetypes.NewListValue(types.StringType, finalizers)
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectNull(MetadataAttrTypes)
	}

	attributes := map[string]attr.Value{
		"name":                          basetypes.NewStringValue(m.Name),
		"namespace":                     basetypes.NewStringValue(m.Namespace),
		"generate_name":                 basetypes.NewStringValue(m.GenerateName),
		"annotations":                   annotationsValue,
		"labels":                        labelsValue,
		"finalizers":                    finalizersValue,
		"owner_references":              ownerReferencesValue,
		"uid":                           basetypes.NewStringValue(m.UID),
		"generation":                    basetypes.NewInt64Value(m.Generation),
		"creation_timestamp":            basetypes.NewStringValue(m.CreationTimestamp),
		"deletion_grace_period_seconds": basetypes.NewInt64Value(m.DeletionGracePeriodSecond),
		"deletion_timestamp":            basetypes.NewStringValue(m.DeletionTimestamp),
		"managed_fields":                basetypes.NewStringValue(string(m.ManagedFields)),
		"resource_version":              basetypes.NewStringValue(m.ResourceVersion),
		"self_link":                     basetypes.NewStringValue(m.SelfLink),
	}
	var d diag.Diagnostics
	obj, d := basetypes.NewObjectValue(MetadataAttrTypes, attributes)
	diags.Append(d...)
	return obj
}

// ToGoModel converts a types.Object to a Metadata struct.
func ToGoModel(ctx context.Context, diags *diag.Diagnostics, metadataObj types.Object) *Metadata {
	if metadataObj.IsNull() || metadataObj.IsUnknown() {
		return &Metadata{}
	}

	// Need a temporary struct to use with As
	type TmpMetadata struct {
		Name                      types.String `tfsdk:"name"`
		Namespace                 types.String `tfsdk:"namespace"`
		GenerateName              types.String `tfsdk:"generate_name"`
		Annotations               types.Map    `tfsdk:"annotations"`
		Labels                    types.Map    `tfsdk:"labels"`
		Finalizers                types.List   `tfsdk:"finalizers"`
		OwnerReferences           types.List   `tfsdk:"owner_references"`
		UID                       types.String `tfsdk:"uid"`
		Generation                types.Int64  `tfsdk:"generation"`
		CreationTimestamp         types.String `tfsdk:"creation_timestamp"`
		DeletionGracePeriodSecond types.Int64  `tfsdk:"deletion_grace_period_seconds"`
		DeletionTimestamp         types.String `tfsdk:"deletion_timestamp"`
		ManagedFields             types.String `tfsdk:"managed_fields"`
		ResourceVersion           types.String `tfsdk:"resource_version"`
		SelfLink                  types.String `tfsdk:"self_link"`
	}
	var tmpMetadata TmpMetadata
	diags.Append(metadataObj.As(ctx, &tmpMetadata, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	obj := &Metadata{}

	obj.Name = tmpMetadata.Name.ValueString()
	obj.Namespace = tmpMetadata.Namespace.ValueString()
	obj.GenerateName = tmpMetadata.GenerateName.ValueString()
	diags.Append(tmpMetadata.Annotations.ElementsAs(ctx, &obj.Annotations, false)...)
	diags.Append(tmpMetadata.Labels.ElementsAs(ctx, &obj.Labels, false)...)
	diags.Append(tmpMetadata.Finalizers.ElementsAs(ctx, &obj.Finalizers, false)...)

	// OwnerReferences
	if !tmpMetadata.OwnerReferences.IsNull() && !tmpMetadata.OwnerReferences.IsUnknown() {
		type TmpOwnerReference struct {
			APIVersion         types.String `tfsdk:"api_version"`
			Kind               types.String `tfsdk:"kind"`
			Name               types.String `tfsdk:"name"`
			UID                types.String `tfsdk:"uid"`
			Controller         types.Bool   `tfsdk:"controller"`
			BlockOwnerDeletion types.Bool   `tfsdk:"block_owner_deletion"`
		}
		var tmpOwnerReferences []TmpOwnerReference
		diags.Append(tmpMetadata.OwnerReferences.ElementsAs(ctx, &tmpOwnerReferences, false)...)
		if diags.HasError() {
			return nil
		}
		for _, tmpOR := range tmpOwnerReferences {
			or := OwnerReference{
				APIVersion:         tmpOR.APIVersion.ValueString(),
				Kind:               tmpOR.Kind.ValueString(),
				Name:               tmpOR.Name.ValueString(),
				UID:                tmpOR.UID.ValueString(),
				Controller:         tmpOR.Controller.ValueBool(),
				BlockOwnerDeletion: tmpOR.BlockOwnerDeletion.ValueBool(),
			}
			obj.OwnerReferences = append(obj.OwnerReferences, or)
		}
	}

	obj.UID = tmpMetadata.UID.ValueString()
	obj.Generation = tmpMetadata.Generation.ValueInt64()
	obj.CreationTimestamp = tmpMetadata.CreationTimestamp.ValueString()
	obj.DeletionGracePeriodSecond = tmpMetadata.DeletionGracePeriodSecond.ValueInt64()
	obj.DeletionTimestamp = tmpMetadata.DeletionTimestamp.ValueString()
	obj.ManagedFields = json.RawMessage(tmpMetadata.ManagedFields.ValueString())
	obj.ResourceVersion = tmpMetadata.ResourceVersion.ValueString()
	obj.SelfLink = tmpMetadata.SelfLink.ValueString()

	return obj
}

// This only includes attributes which are able to be sent to the API.
type ApiRequestMetadata struct {
	Name            string            `json:"name,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	GenerateName    string            `json:"generate_name,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Finalizers      []string          `json:"finalizers,omitempty"`
	OwnerReferences []OwnerReference  `json:"owner_references,omitempty"`
}

func (m *OwnerReference) ToTypesObject(ctx context.Context, diags *diag.Diagnostics) types.Object {
	obj, d := basetypes.NewObjectValueFrom(ctx, OwnerReferenceAttrTypes, m)
	diags.Append(d...)
	return obj
}

// Sample data functions all return the same data in different forms.
func SampleMetadataTypesObject() types.Object {
	return types.ObjectValueMust(
		MetadataAttrTypes,
		map[string]attr.Value{
			"name":          types.StringValue("test_metadata"),
			"namespace":     types.StringValue("test_namespace"),
			"generate_name": types.StringValue(""), // mutually exclusive to the name attribute
			"annotations": types.MapValueMust(types.StringType, map[string]attr.Value{
				"ann_key": types.StringValue("ann_value"),
			}),
			"labels": types.MapValueMust(types.StringType, map[string]attr.Value{
				"label_key": types.StringValue("label_value"),
			}),
			"finalizers": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("finalizer_a"),
			}),
			"owner_references": types.ListValueMust(types.ObjectType{AttrTypes: OwnerReferenceAttrTypes}, []attr.Value{
				types.ObjectValueMust(OwnerReferenceAttrTypes, map[string]attr.Value{
					"api_version":          types.StringValue("v1"),
					"kind":                 types.StringValue("some_kind"),
					"name":                 types.StringValue("owner"),
					"uid":                  types.StringValue("some_uid"),
					"controller":           types.BoolValue(true),
					"block_owner_deletion": types.BoolValue(true),
				},
				)}),
			"uid":                           types.StringValue("test_uid"),
			"generation":                    types.Int64Value(1),
			"creation_timestamp":            types.StringValue("2023-01-01T00:00:00Z"),
			"deletion_grace_period_seconds": types.Int64Value(30),
			"deletion_timestamp":            types.StringValue("2023-01-01T01:00:00Z"),
			"managed_fields":                types.StringValue("{\"field\": \"test_managed_fields\"}"),
			"resource_version":              types.StringValue("test_resource_version"),
			"self_link":                     types.StringValue("/api/v1/namespaces/default/rancher2_dev2s/test"),
		},
	)
}

func SampleMetadataGoModel() Metadata {
	return Metadata{
		Name:      "test_metadata",
		Namespace: "test_namespace",
		// GenerateName: null, // mutually exclusive to the Name attribute
		Annotations: map[string]string{"ann_key": "ann_value"},
		Labels:      map[string]string{"label_key": "label_value"},
		Finalizers:  []string{"finalizer_a"},
		OwnerReferences: []OwnerReference{
			{
				APIVersion:         "v1",
				Kind:               "some_kind",
				Name:               "owner",
				UID:                "some_uid",
				Controller:         true,
				BlockOwnerDeletion: true,
			},
		},
		UID:                       "test_uid",
		Generation:                1,
		CreationTimestamp:         "2023-01-01T00:00:00Z",
		DeletionGracePeriodSecond: 30,
		DeletionTimestamp:         "2023-01-01T01:00:00Z",
		ManagedFields:             json.RawMessage("{\"field\": \"test_managed_fields\"}"),
		ResourceVersion:           "test_resource_version",
		SelfLink:                  "/api/v1/namespaces/default/rancher2_dev2s/test",
	}
}

func SampleMetadataApiRequestJson() string {
	// "generate_name": "test_generate_name", mutually exclusive to the name attribute
	return `"metadata": {
    "name": "test_metadata",
    "namespace": "test_namespace",
    "annotations": {
      "ann_key": "ann_value"
    },
    "labels": {
      "label_key": "label_value"
    },
    "finalizers": [
      "finalizer_a"
    ],
    "owner_references": [
      {
        "api_version": "v1",
        "kind": "some_kind",
        "name": "owner",
        "uid": "some_uid",
        "controller": true,
        "block_owner_deletion": true
      }
    ]
  },`
}
