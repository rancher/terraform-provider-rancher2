package pretty_print

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// PrettyPrint doesn't redact any data, use judiciously.
func PrettyPrint(data any) string {
	printable := toPrintable(data)
	s, err := json.MarshalIndent(printable, "", "  ")
	if err != nil {
		return fmt.Sprintf("failed to pretty-print object: %v", err)
	}
	return string(s)
}

type unknownAndNullable interface {
	IsUnknown() bool
	IsNull() bool
}

func toPrintableFramework[T any](v unknownAndNullable, valuer func() T) any {
	if v.IsUnknown() {
		return "<unknown>"
	}
	if v.IsNull() {
		return nil
	}
	return valuer()
}

// toPrintable is a helper that recursively converts data structures,
// especially those containing terraform-plugin-framework types and tftypes.Value,
// into a representation that is friendly for JSON marshaling.
func toPrintable(data any) any {
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
			printable["AttributePath"] = toPrintable(dp.Path())
		}
		return printable
	}

	if marshaller, ok := data.(json.Marshaler); ok {
		return marshaller
	}

	val := reflect.ValueOf(data)

	// Dereference pointer first.
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil
		}
		return toPrintable(val.Elem().Interface())
	}

	// Handle special Terraform types that need conversion.
	switch v := data.(type) {
	case tftypes.Value:
		return toPrintableTftypesValue(v)
	case types.String:
		return toPrintableFramework(v, v.ValueString)
	case types.Bool:
		return toPrintableFramework(v, v.ValueBool)
	case types.Int64:
		return toPrintableFramework(v, v.ValueInt64)
	case types.Int32:
		return toPrintableFramework(v, v.ValueInt32)
	case types.Float64:
		return toPrintableFramework(v, v.ValueFloat64)
	case types.Float32:
		return toPrintableFramework(v, v.ValueFloat32)
	case types.Number:
		return toPrintableFramework(v, v.ValueBigFloat)
	}

	// Handle native Go kinds.
	switch val.Kind() {
	case reflect.Struct:
		out := make(map[string]any)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			// Skip unexported fields.
			if field.PkgPath != "" {
				continue
			}
			value := val.Field(i)
			out[field.Name] = toPrintable(value.Interface())
		}
		return out

	case reflect.Slice, reflect.Array:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return string(val.Bytes())
		}
		out := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			out[i] = toPrintable(val.Index(i).Interface())
		}
		return out

	case reflect.Map:
		out := make(map[string]any)
		for _, key := range val.MapKeys() {
			out[fmt.Sprintf("%v", key.Interface())] = toPrintable(val.MapIndex(key).Interface())
		}
		return out
	}

	// For primitive types, return as is.
	return data
}

// toPrintableTftypesValue handles the conversion of tftypes.Value to a JSON-marshalable format.
func toPrintableTftypesValue(v tftypes.Value) any {
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
			result[k] = toPrintable(v)
		}
		return result
	case tftypes.List, tftypes.Set: // Sets can be treated as lists for printing
		elems := make([]tftypes.Value, 0)
		if err := v.As(&elems); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.List/Set: %s>", err)
		}
		result := make([]any, len(elems))
		for i, v := range elems {
			result[i] = toPrintable(v)
		}
		return result
	case tftypes.Map:
		elems := make(map[string]tftypes.Value)
		if err := v.As(&elems); err != nil {
			return fmt.Sprintf("<unmarshallable tftypes.Map: %s>", err)
		}
		result := make(map[string]any)
		for k, v := range elems {
			result[k] = toPrintable(v)
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
	return toPrintable(goValue)
}
