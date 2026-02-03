package helpers

import (
	// "fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FixMapType is a generalized way to correct a map attribute with no elements.
//
// If the elements of a map aren't set then tfsdk.Set won't know how to properly propagate.
func FixMapType(v *types.Map, atr schema.MapAttribute, diags *diag.Diagnostics) types.Map {
	elements := v.Elements()
	if elements == nil {
		elements = map[string]attr.Value{}
	}

	newValue, valDiags := types.MapValue(atr.ElementType, elements)
	diags.Append(valDiags...)

	return newValue
}

func FixListType(v *types.List, atr schema.ListAttribute, diags *diag.Diagnostics) types.List {

	elements := v.Elements()
	if len(elements) == 0 {
		elements = []attr.Value{}
	}
	newValue, valDiags := types.ListValue(atr.ElementType, elements)
	diags.Append(valDiags...)

	return newValue
}

func FixSetType(v *types.Set, atr schema.SetAttribute, diags *diag.Diagnostics) types.Set {

	elements := v.Elements()
	if elements == nil {
		elements = []attr.Value{}
	}

	newValue, valDiags := types.SetValue(atr.ElementType, elements)
	diags.Append(valDiags...)

	return newValue
}
