package pretty_print

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// PrettyPrint doesn't redact any data, use judiciously.
func PrettyPrint(data any) string {
	printable := ToPrintable(data)
	s, err := json.MarshalIndent(printable, "", "  ")
	if err != nil {
		fmt.Printf("%+v", printable)
		return fmt.Sprintf("failed to pretty-print object: %v", err)
	}
	return string(s)
}

type unknownAndNullable interface {
	IsUnknown() bool
	IsNull() bool
}

func ToPrintableFramework[T any](v unknownAndNullable, valuer func() T) any {
	if v.IsUnknown() {
		return "<unknown>"
	}
	if v.IsNull() {
		return nil
	}
	return valuer()
}

// ToPrintable is a helper that recursively converts data structures,
// especially those containing terraform-plugin-framework types and tftypes.Value,
// into a representation that is friendly for JSON marshaling.
func ToPrintable(data any) any {
	if data == nil {
		return nil
	}

	if d, ok := data.(diag.Diagnostic); ok {
		severity := ""
		switch d.Severity() {
		case diag.SeverityError:
			severity = "Error"
		case diag.SeverityWarning:
			severity = "Warning"
		}
		printable := map[string]any{
			"Severity": severity,
			"Summary":  d.Summary(),
			"Detail":   d.Detail(),
		}
		if dp, ok := d.(diag.DiagnosticWithPath); ok {
			printable["AttributePath"] = ToPrintable(dp.Path())
		}
		return printable
	}

	if raw, ok := data.(json.RawMessage); ok {
		var parsed any
		if err := json.Unmarshal(raw, &parsed); err == nil {
			return ToPrintable(parsed)
		}
		// Fallback to basic string conversion if Unmarshal fails (e.g., invalid JSON)
		// The goal is for these logs not to fail even if the data is the wrong type
		return string(raw)
	}

	val := reflect.ValueOf(data)

	// Dereference pointer first.
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil
		}
		return ToPrintable(val.Elem().Interface())
	}

	// Handle special Terraform types that need conversion.
	switch v := data.(type) {
	case tftypes.Value:
		// object, map, set, and list are handled here
		return ToPrintableTftypesValue(v)
	case types.String:
		return ToPrintableFramework(v, v.ValueString)
	case types.Bool:
		return ToPrintableFramework(v, v.ValueBool)
	case types.Int64:
		return ToPrintableFramework(v, v.ValueInt64)
	case types.Int32:
		return ToPrintableFramework(v, v.ValueInt32)
	case types.Float64:
		return ToPrintableFramework(v, v.ValueFloat64)
	case types.Float32:
		return ToPrintableFramework(v, v.ValueFloat32)
	case types.Number:
		return ToPrintableFramework(v, v.ValueBigFloat)
	case basetypes.ObjectValue:
		if v.IsUnknown() {
			return "<unknown>"
		}
		if v.IsNull() {
			return nil
		}
		out := make(map[string]any)
		for k, val := range v.Attributes() {
			out[k] = ToPrintable(val)
		}
		return out
	case basetypes.MapValue:
		if v.IsUnknown() {
			return "<unknown>"
		}
		if v.IsNull() {
			return nil
		}
		out := make(map[string]any)
		for k, val := range v.Elements() {
			out[k] = ToPrintable(val)
		}
		return out
	case basetypes.ListValue:
		if v.IsUnknown() {
			return "<unknown>"
		}
		if v.IsNull() {
			return nil
		}
		out := make([]any, 0, len(v.Elements()))
		for _, val := range v.Elements() {
			out = append(out, ToPrintable(val))
		}
		return out
	case basetypes.SetValue:
		if v.IsUnknown() {
			return "<unknown>"
		}
		if v.IsNull() {
			return nil
		}
		out := make([]any, 0, len(v.Elements()))
		for _, val := range v.Elements() {
			out = append(out, ToPrintable(val))
		}
		return out
	}

	// Handle native Go kinds.
	switch val.Kind() {
	case reflect.Struct:
		// fmt.Printf("pretty printing Go native struct: \n%+v\n", val)
		out := make(map[string]any)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			// Skip unexported fields to prevent panics.
			if field.PkgPath != "" {
				continue
			}
			value := val.Field(i)
			out[field.Name] = ToPrintable(value.Interface())
		}
		return out

	case reflect.Slice, reflect.Array:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return string(val.Bytes())
		}
		out := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			out[i] = ToPrintable(val.Index(i).Interface())
		}
		return out

	case reflect.Map:
		out := make(map[string]any)
		for _, key := range val.MapKeys() {
			out[fmt.Sprintf("%v", key.Interface())] = ToPrintable(val.MapIndex(key).Interface())
		}
		return out
	}
	// For primitive types, return as is.
	return data
}

// ToPrintableTftypesValue handles the conversion of tftypes.Value to a JSON-marshalable format.
func ToPrintableTftypesValue(v tftypes.Value) any {
	// fmt.Printf("\nhandling tftypes value: \n%+v\n", v)
	if !v.IsKnown() {
		return "<unknown>"
	}
	if v.IsNull() {
		return nil
	}

	ty := v.Type()

	// Handle collection types by recursively calling this function.
	switch ty.(type) {
	case tftypes.Object:
		attrs := make(map[string]tftypes.Value)
		if err := v.As(&attrs); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.Object: %s>", err)
		}
		result := make(map[string]any)
		for k, v := range attrs {
			result[k] = ToPrintable(v)
		}
		return result
		// Handle collection types by recursively calling this function.
	case tftypes.List, tftypes.Set: // Sets can be treated as lists for printing
		elems := make([]tftypes.Value, 0)
		if err := v.As(&elems); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.List/Set: %s>", err)
		}
		result := make([]any, len(elems))
		for i, v := range elems {
			result[i] = ToPrintable(v)
		}
		return result
	case tftypes.Map:
		elems := make(map[string]tftypes.Value)
		if err := v.As(&elems); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.Map: %s>", err)
		}
		result := make(map[string]any)
		for k, v := range elems {
			result[k] = ToPrintable(v)
		}
		return result
	}

	// It wasn't a collection type, so it must be a primitive type.
	// We can't type-switch on the primitive type variables, so we'll
	// compare against them directly.
	if ty.Equal(tftypes.String) {
		var s string
		if err := v.As(&s); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.String: %s>", err)
		}
		return s
	}
	if ty.Equal(tftypes.Number) {
		var bf big.Float
		if err := v.As(&bf); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.Number: %s>", err)
		}
		return &bf
	}
	if ty.Equal(tftypes.Bool) {
		var b bool
		if err := v.As(&b); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.Bool: %s>", err)
		}
		return b
	}

	// For any other types, attempt to unmarshal into a generic interface.
	var goValue any
	if err := v.As(&goValue); err != nil {
		return fmt.Sprintf("<unmarshallable tftypes.Value: %s>", v.Type().String())
	}
	return ToPrintable(goValue)
}
