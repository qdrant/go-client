package qdrant

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewValue_Primitives(t *testing.T) {
	// nil
	val, err := NewValue(nil)
	require.NoError(t, err)
	require.Equal(t, NullValue_NULL_VALUE, val.GetNullValue())

	// bool
	val, err = NewValue(true)
	require.NoError(t, err)
	require.True(t, val.GetBoolValue())

	// signed ints
	val, err = NewValue(int(42))
	require.NoError(t, err)
	require.Equal(t, int64(42), val.GetIntegerValue())

	val, err = NewValue(int8(-8))
	require.NoError(t, err)
	require.Equal(t, int64(-8), val.GetIntegerValue())

	val, err = NewValue(int16(-16))
	require.NoError(t, err)
	require.Equal(t, int64(-16), val.GetIntegerValue())

	val, err = NewValue(int32(-32))
	require.NoError(t, err)
	require.Equal(t, int64(-32), val.GetIntegerValue())

	val, err = NewValue(int64(-64))
	require.NoError(t, err)
	require.Equal(t, int64(-64), val.GetIntegerValue())

	// unsigned ints
	val, err = NewValue(uint(100))
	require.NoError(t, err)
	require.Equal(t, int64(100), val.GetIntegerValue())

	val, err = NewValue(uint8(255))
	require.NoError(t, err)
	require.Equal(t, int64(255), val.GetIntegerValue())

	val, err = NewValue(uint16(65535))
	require.NoError(t, err)
	require.Equal(t, int64(65535), val.GetIntegerValue())

	val, err = NewValue(uint32(100000))
	require.NoError(t, err)
	require.Equal(t, int64(100000), val.GetIntegerValue())

	val, err = NewValue(uint64(5000000))
	require.NoError(t, err)
	require.Equal(t, int64(5000000), val.GetIntegerValue())

	// floats
	val, err = NewValue(float32(3.14))
	require.NoError(t, err)
	require.InDelta(t, float64(float32(3.14)), val.GetDoubleValue(), 0.0001)

	val, err = NewValue(float64(2.71828))
	require.NoError(t, err)
	require.Equal(t, 2.71828, val.GetDoubleValue())

	// string
	val, err = NewValue("hello world")
	require.NoError(t, err)
	require.Equal(t, "hello world", val.GetStringValue())

	// []byte (base64)
	bytesInput := []byte("binary data")
	val, err = NewValue(bytesInput)
	require.NoError(t, err)
	require.Equal(t, base64.StdEncoding.EncodeToString(bytesInput), val.GetStringValue())
}

func TestNewValue_ProtobufPointers(t *testing.T) {
	// *Value
	existing := NewValueInt(123)
	val, err := NewValue(existing)
	require.NoError(t, err)
	require.Same(t, existing, val)

	var nilVal *Value
	val, err = NewValue(nilVal)
	require.NoError(t, err)
	require.Equal(t, NullValue_NULL_VALUE, val.GetNullValue())

	// *Struct
	existingStruct := &Struct{Fields: map[string]*Value{"k": NewValueString("v")}}
	val, err = NewValue(existingStruct)
	require.NoError(t, err)
	require.Equal(t, existingStruct, val.GetStructValue())

	var nilStruct *Struct
	val, err = NewValue(nilStruct)
	require.NoError(t, err)
	require.Equal(t, NullValue_NULL_VALUE, val.GetNullValue())

	// *ListValue
	existingList := &ListValue{Values: []*Value{NewValueString("item")}}
	val, err = NewValue(existingList)
	require.NoError(t, err)
	require.Equal(t, existingList, val.GetListValue())

	var nilList *ListValue
	val, err = NewValue(nilList)
	require.NoError(t, err)
	require.Equal(t, NullValue_NULL_VALUE, val.GetNullValue())
}

func TestNewValue_Slices(t *testing.T) {
	// []string
	val, err := NewValue([]string{"a", "b"})
	require.NoError(t, err)
	require.Len(t, val.GetListValue().GetValues(), 2)
	require.Equal(t, "a", val.GetListValue().GetValues()[0].GetStringValue())
	require.Equal(t, "b", val.GetListValue().GetValues()[1].GetStringValue())

	// []bool
	val, err = NewValue([]bool{true, false})
	require.NoError(t, err)
	require.Len(t, val.GetListValue().GetValues(), 2)
	require.True(t, val.GetListValue().GetValues()[0].GetBoolValue())
	require.False(t, val.GetListValue().GetValues()[1].GetBoolValue())

	// []int
	val, err = NewValue([]int{1, 2, 3})
	require.NoError(t, err)
	require.Len(t, val.GetListValue().GetValues(), 3)
	require.Equal(t, int64(1), val.GetListValue().GetValues()[0].GetIntegerValue())

	// []int64
	val, err = NewValue([]int64{10, 20})
	require.NoError(t, err)
	require.Len(t, val.GetListValue().GetValues(), 2)
	require.Equal(t, int64(10), val.GetListValue().GetValues()[0].GetIntegerValue())

	// []float64
	val, err = NewValue([]float64{1.1, 2.2})
	require.NoError(t, err)
	require.Len(t, val.GetListValue().GetValues(), 2)
	require.Equal(t, 1.1, val.GetListValue().GetValues()[0].GetDoubleValue())

	// []*Value
	val, err = NewValue([]*Value{NewValueInt(7), NewValueString("eight")})
	require.NoError(t, err)
	require.Len(t, val.GetListValue().GetValues(), 2)
	require.Equal(t, int64(7), val.GetListValue().GetValues()[0].GetIntegerValue())
	require.Equal(t, "eight", val.GetListValue().GetValues()[1].GetStringValue())
}

func TestNewValueMap_Complex(t *testing.T) {
	jsonMap := map[string]any{
		"null":    nil,
		"bool":    true,
		"int":     int(42),
		"int8":    int8(8),
		"int16":   int16(16),
		"int32":   int32(32),
		"int64":   int64(64),
		"uint":    uint(1),
		"uint8":   uint8(2),
		"uint16":  uint16(3),
		"uint32":  uint32(4),
		"uint64":  uint64(5),
		"float32": float32(1.5),
		"float64": float64(2.5),
		"string":  "test",
		"bytes":   []byte("raw"),
		"tags":    []string{"fast", "reliable"},
		"counts":  []int{10, 20},
		"nested":  map[string]any{"sub_key": "sub_value"},
	}

	vm := NewValueMap(jsonMap)
	require.Equal(t, NullValue_NULL_VALUE, vm["null"].GetNullValue())
	require.True(t, vm["bool"].GetBoolValue())
	require.Equal(t, int64(42), vm["int"].GetIntegerValue())
	require.Equal(t, int64(8), vm["int8"].GetIntegerValue())
	require.Equal(t, int64(16), vm["int16"].GetIntegerValue())
	require.Equal(t, int64(32), vm["int32"].GetIntegerValue())
	require.Equal(t, int64(64), vm["int64"].GetIntegerValue())
	require.Equal(t, int64(1), vm["uint"].GetIntegerValue())
	require.Equal(t, int64(2), vm["uint8"].GetIntegerValue())
	require.Equal(t, int64(3), vm["uint16"].GetIntegerValue())
	require.Equal(t, int64(4), vm["uint32"].GetIntegerValue())
	require.Equal(t, int64(5), vm["uint64"].GetIntegerValue())
	require.Equal(t, "test", vm["string"].GetStringValue())
	require.Equal(t, "fast", vm["tags"].GetListValue().GetValues()[0].GetStringValue())
	require.Equal(t, "reliable", vm["tags"].GetListValue().GetValues()[1].GetStringValue())
	require.Equal(t, int64(10), vm["counts"].GetListValue().GetValues()[0].GetIntegerValue())
	require.Equal(t, "sub_value", vm["nested"].GetStructValue().GetFields()["sub_key"].GetStringValue())
}

func TestNewValue_InvalidType(t *testing.T) {
	type CustomStruct struct{}
	_, err := NewValue(CustomStruct{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid type")

	invalidUTF8 := string([]byte{0xff, 0xfe, 0xfd})
	_, err = NewValue(invalidUTF8)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid UTF-8")
}
