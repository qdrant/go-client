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
//	╔════════════════════════╤════════════════════════════════════════════╗
//	║ Go type                │ Conversion                                 ║
//	╠════════════════════════╪════════════════════════════════════════════╣
//	║ nil                    │ stored as NullValue                        ║
//	║ bool                   │ stored as BoolValue                        ║
//	║ int, int32, int64      │ stored as IntegerValue                     ║
//	║ uint, uint32, uint64   │ stored as IntegerValue                     ║
//	║ float32, float64       │ stored as DoubleValue                      ║
//	║ string                 │ stored as StringValue; must be valid UTF-8 ║
//	║ []byte                 │ stored as StringValue; base64-encoded      ║
//	║ map[string]interface{} │ stored as StructValue                      ║
//	║ []interface{}          │ stored as ListValue                        ║
//	╚════════════════════════╧════════════════════════════════════════════╝
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
	case bool:
		return NewValueBool(v), nil
	case int:
		return NewValueInt(int64(v)), nil
	case int32:
		return NewValueInt(int64(v)), nil
	case int64:
		return NewValueInt(v), nil
	case uint:
		return NewValueInt(int64(v)), nil
	case uint32:
		return NewValueInt(int64(v)), nil
	case uint64:
		return NewValueInt(int64(v)), nil
	case float32:
		return NewValueDouble(float64(v)), nil
	case float64:
		return NewValueDouble(float64(v)), nil
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
	default:
		return nil, fmt.Errorf("invalid type: %T", v)
	}
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

// Constructs a new list Value.
func NewValueList(v *ListValue) *Value {
	return &Value{Kind: &Value_ListValue{ListValue: v}}
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
