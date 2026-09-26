// This file contains methods to convert a generic map to map of string to *grpc.Value(Qdrant's payload type).
// This is a custom implementatation based on "google.golang.org/protobuf/types/known/structpb".
// It extends the original implementation to support IntegerValue and DoubleValue instead of a single NumberValue.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/json_with_int.proto
//
// USAGE:
//
//	jsonMap := map[string]any{
//		"some_null":    nil,
//		"some_bool":    true,
//		"some_int":     42,
//		"some_float":   3.14,
//		"some_string":  "hello",
//		"some_bytes":   []byte("world"),
//		"some_nested":  map[string]any{"key": "value"},
//		"some_list":    []any{"foo", 32},
//	}
//
//	valueMap := NewValueMap(jsonMap)

package qdrant

import (
	"encoding/base64"
	"fmt"
	"unicode/utf8"
)

// Converts a map of string to any to a map of string to *grpc.Value
// NOTE: This function panics if the conversion fails. Use TryValueMap() to have errors returned.
//
//	╔════════════════════════════════════════════════════════╤════════════════════════════════════════════╗
//	║ Go type                                                │ Conversion                                 ║
//	╠════════════════════════════════════════════════════════╪════════════════════════════════════════════╣
//	║ nil                                                    │ stored as NullValue                        ║
//	║ bool                                                   │ stored as BoolValue                        ║
//	║ int, int8, int16, int32, int64                         │ stored as IntegerValue                     ║
//	║ uint, uint8, uint16, uint32, uint64                    │ stored as IntegerValue                     ║
//	║ float32, float64                                       │ stored as DoubleValue                      ║
//	║ string                                                 │ stored as StringValue; must be valid UTF-8 ║
//	║ []byte                                                 │ stored as StringValue; base64-encoded      ║
//	║ *Value                                                 │ returned as-is (nil as NullValue)          ║
//	║ *Struct                                                │ stored as StructValue                      ║
//	║ *ListValue                                             │ stored as ListValue                        ║
//	║ map[string]interface{}                                 │ stored as StructValue                      ║
//	║ []interface{}, []string, []bool, []int, []int64, etc.  │ stored as ListValue                        ║
//	║ []*Value                                               │ stored as ListValue                        ║
//	╚════════════════════════════════════════════════════════╧════════════════════════════════════════════╝
func NewValueMap(inputMap map[string]any) map[string]*Value {
	valueMap, err := TryValueMap(inputMap)
	if err != nil {
		panic(err)
	}
	return valueMap
}

// Converts a map of string to any to a map of string to *grpc.Value
// Returns an error if the conversion fails.
func TryValueMap(inputMap map[string]any) (map[string]*Value, error) {
	valueMap := make(map[string]*Value)
	for key, val := range inputMap {
		value, err := NewValue(val)
		if err != nil {
			return nil, err
		}
		valueMap[key] = value
	}
	return valueMap, nil
}

// Constructs a *Value from a generic Go interface.
func NewValue(v any) (*Value, error) {
	switch v := v.(type) {
	case nil:
		return NewValueNull(), nil
	case *Value:
		if v == nil {
			return NewValueNull(), nil
		}
		return v, nil
	case *Struct:
		if v == nil {
			return NewValueNull(), nil
		}
		return NewValueStruct(v), nil
	case *ListValue:
		if v == nil {
			return NewValueNull(), nil
		}
		return NewValueList(v), nil
	case bool:
		return NewValueBool(v), nil
	case string:
		if !utf8.ValidString(v) {
			return nil, fmt.Errorf("invalid UTF-8 in string: %q", v)
		}
		return NewValueString(v), nil
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return NewValueString(s), nil
	case map[string]interface{}:
		v2, err := NewStruct(v)
		if err != nil {
			return nil, err
		}
		return NewValueStruct(v2), nil
	case []interface{}:
		v2, err := NewListValue(v)
		if err != nil {
			return nil, err
		}
		return NewValueList(v2), nil
	case []*Value:
		return NewValueFromList(v...), nil
	default:
		return newValueNumeric(v)
	}
}

// newValueNumeric converts numeric scalars and typed slices into a *Value.
// Anything else is delegated to newValueTypedSlice, which reports the
// invalid-type error for genuinely unsupported values.
func newValueNumeric(v any) (*Value, error) {
	switch v := v.(type) {
	case int:
		return NewValueInt(int64(v)), nil
	case int8:
		return NewValueInt(int64(v)), nil
	case int16:
		return NewValueInt(int64(v)), nil
	case int32:
		return NewValueInt(int64(v)), nil
	case int64:
		return NewValueInt(v), nil
	case uint:
		return NewValueInt(int64(v)), nil
	case uint8:
		return NewValueInt(int64(v)), nil
	case uint16:
		return NewValueInt(int64(v)), nil
	case uint32:
		return NewValueInt(int64(v)), nil
	case uint64:
		return NewValueInt(int64(v)), nil
	case float32:
		return NewValueDouble(float64(v)), nil
	case float64:
		return NewValueDouble(v), nil
	default:
		return newValueTypedSlice(v)
	}
}

// newValueTypedSlice converts typed Go slices into a ListValue.
func newValueTypedSlice(v any) (*Value, error) {
	switch v := v.(type) {
	case []string:
		list := make([]*Value, len(v))
		for i, item := range v {
			if !utf8.ValidString(item) {
				return nil, fmt.Errorf("invalid UTF-8 in string: %q", item)
			}
			list[i] = NewValueString(item)
		}
		return NewValueFromList(list...), nil
	case []bool:
		return newValueBoolSlice(v), nil
	case []int:
		return newValueIntSlice(v), nil
	case []int8:
		return newValueIntSlice(v), nil
	case []int16:
		return newValueIntSlice(v), nil
	case []int32:
		return newValueIntSlice(v), nil
	case []int64:
		return newValueIntSlice(v), nil
	case []uint:
		return newValueIntSlice(v), nil
	case []uint16:
		return newValueIntSlice(v), nil
	case []uint32:
		return newValueIntSlice(v), nil
	case []uint64:
		return newValueIntSlice(v), nil
	case []float32:
		return newValueFloatSlice(v), nil
	case []float64:
		return newValueFloatSlice(v), nil
	default:
		return nil, fmt.Errorf("invalid type: %T", v)
	}
}

// newValueIntSlice converts a typed integer slice into a ListValue
// where each element is stored as an IntegerValue.
func newValueIntSlice[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](
	items []T,
) *Value {
	list := make([]*Value, len(items))
	for i, item := range items {
		list[i] = NewValueInt(int64(item))
	}
	return NewValueFromList(list...)
}

// newValueFloatSlice converts a typed float slice into a ListValue
// where each element is stored as a DoubleValue.
func newValueFloatSlice[T float32 | float64](items []T) *Value {
	list := make([]*Value, len(items))
	for i, item := range items {
		list[i] = NewValueDouble(float64(item))
	}
	return NewValueFromList(list...)
}

// newValueBoolSlice converts a bool slice into a ListValue.
func newValueBoolSlice(items []bool) *Value {
	list := make([]*Value, len(items))
	for i, item := range items {
		list[i] = NewValueBool(item)
	}
	return NewValueFromList(list...)
}

// Constructs a new null Value.
func NewValueNull() *Value {
	return &Value{Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}}
}

// Constructs a new boolean Value.
func NewValueBool(v bool) *Value {
	return &Value{Kind: &Value_BoolValue{BoolValue: v}}
}

// Constructs a new integer Value.
func NewValueInt(v int64) *Value {
	return &Value{Kind: &Value_IntegerValue{IntegerValue: v}}
}

// Constructs a new double Value.
func NewValueDouble(v float64) *Value {
	return &Value{Kind: &Value_DoubleValue{DoubleValue: v}}
}

// Constructs a new string Value.
func NewValueString(v string) *Value {
	return &Value{Kind: &Value_StringValue{StringValue: v}}
}

// Constructs a new struct Value.
func NewValueStruct(v *Struct) *Value {
	return &Value{Kind: &Value_StructValue{StructValue: v}}
}

// Constructs a new struct Value from the provided field map.
func NewValueFromFields(fields map[string]*Value) *Value {
	return &Value{Kind: &Value_StructValue{StructValue: &Struct{
		Fields: fields,
	}}}
}

// Constructs a new list Value.
func NewValueList(v *ListValue) *Value {
	return &Value{Kind: &Value_ListValue{ListValue: v}}
}

// Constructs a new list Value from the provided elements.
func NewValueFromList(values ...*Value) *Value {
	return &Value{Kind: &Value_ListValue{ListValue: &ListValue{
		Values: values,
	}}}
}

// Constructs a ListValue from a general-purpose Go slice.
// The slice elements are converted using NewValue().
func NewListValue(v []interface{}) (*ListValue, error) {
	x := &ListValue{Values: make([]*Value, len(v))}
	for i, v := range v {
		var err error
		x.Values[i], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// Constructs a Struct from a general-purpose Go map.
// The map keys must be valid UTF-8.
// The map values are converted using NewValue().
func NewStruct(v map[string]interface{}) (*Struct, error) {
	x := &Struct{Fields: make(map[string]*Value, len(v))}
	for k, v := range v {
		if !utf8.ValidString(k) {
			return nil, fmt.Errorf("invalid UTF-8 in string: %q", k)
		}
		var err error
		x.Fields[k], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}
