package qdrant_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/qdrant/go-client/qdrant"
)

func TestNewValue(t *testing.T) {
	null := qdrant.NewValueNull()
	value := qdrant.NewValueInt(1)
	structValue := &qdrant.Struct{Fields: map[string]*qdrant.Value{"k": qdrant.NewValueString("v")}}
	listValue := &qdrant.ListValue{Values: []*qdrant.Value{qdrant.NewValueString("v")}}
	one, two := qdrant.NewValueInt(1), qdrant.NewValueInt(2)
	oneAndTwo := qdrant.NewValueFromList(one, two)
	field := func(v *qdrant.Value) *qdrant.Value {
		return qdrant.NewValueFromFields(map[string]*qdrant.Value{"k": v})
	}

	tests := []struct {
		input any
		want  *qdrant.Value
	}{
		{nil, null},
		{(*qdrant.Value)(nil), null},
		{(*qdrant.Struct)(nil), null},
		{(*qdrant.ListValue)(nil), null},
		{value, value},
		{structValue, qdrant.NewValueStruct(structValue)},
		{listValue, qdrant.NewValueList(listValue)},
		{true, qdrant.NewValueBool(true)},
		{int(-1), qdrant.NewValueInt(-1)},
		{int8(-1), qdrant.NewValueInt(-1)},
		{int16(-1), qdrant.NewValueInt(-1)},
		{int32(-1), qdrant.NewValueInt(-1)},
		{int64(-1), qdrant.NewValueInt(-1)},
		{uint(1), one},
		{uint8(1), one},
		{uint16(1), one},
		{uint32(1), one},
		{uint64(math.MaxInt64), qdrant.NewValueInt(math.MaxInt64)},
		{float32(1.5), qdrant.NewValueDouble(1.5)},
		{float64(1.5), qdrant.NewValueDouble(1.5)},
		{"v", qdrant.NewValueString("v")},
		{[]byte("raw"), qdrant.NewValueString("cmF3")},
		{map[string]any{"k": "v"}, qdrant.NewValueStruct(structValue)},
		{[]any{1, "v"}, qdrant.NewValueFromList(one, qdrant.NewValueString("v"))},
		{[]string{"v"}, qdrant.NewValueList(listValue)},
		{[]bool{true}, qdrant.NewValueFromList(qdrant.NewValueBool(true))},
		{[]int{1, 2}, oneAndTwo},
		{[]int8{1, 2}, oneAndTwo},
		{[]int16{1, 2}, oneAndTwo},
		{[]int32{1, 2}, oneAndTwo},
		{[]int64{1, 2}, oneAndTwo},
		{[]uint{1, 2}, oneAndTwo},
		{[]uint16{1, 2}, oneAndTwo},
		{[]uint32{1, 2}, oneAndTwo},
		{[]uint64{1, 2}, oneAndTwo},
		{[]float32{1.5}, qdrant.NewValueFromList(qdrant.NewValueDouble(1.5))},
		{[]float64{1.5}, qdrant.NewValueFromList(qdrant.NewValueDouble(1.5))},
		{[]*qdrant.Value{one, nil}, qdrant.NewValueFromList(one, null)},
		{map[string]*qdrant.Value{"k": nil}, field(null)},
		{map[string]string{"k": "v"}, qdrant.NewValueStruct(structValue)},
		{map[string]bool{"k": true}, field(qdrant.NewValueBool(true))},
		{map[string]int{"k": 1}, field(one)},
		{map[string]int8{"k": 1}, field(one)},
		{map[string]int16{"k": 1}, field(one)},
		{map[string]int32{"k": 1}, field(one)},
		{map[string]int64{"k": 1}, field(one)},
		{map[string]uint{"k": 1}, field(one)},
		{map[string]uint8{"k": 1}, field(one)},
		{map[string]uint16{"k": 1}, field(one)},
		{map[string]uint32{"k": 1}, field(one)},
		{map[string]uint64{"k": 1}, field(one)},
		{map[string]float32{"k": 1.5}, field(qdrant.NewValueDouble(1.5))},
		{map[string]float64{"k": 1.5}, field(qdrant.NewValueDouble(1.5))},
	}
	for _, tt := range tests {
		got, err := qdrant.NewValue(tt.input)
		require.NoError(t, err, "%T", tt.input)
		require.True(t, proto.Equal(tt.want, got), "%T: got %v, want %v", tt.input, got, tt.want)
	}
}

func TestValue_AsInterface(t *testing.T) {
	input := map[string]any{
		"null":   nil,
		"bool":   true,
		"int":    int64(1),
		"double": 1.5,
		"string": "v",
		"list":   []any{int64(1), "v"},
		"struct": map[string]any{"k": "v"},
	}
	value, err := qdrant.NewValue(input)
	require.NoError(t, err)
	require.Equal(t, input, value.AsInterface())
	require.Equal(t, input, qdrant.ValueMapToMap(qdrant.NewValueMap(input)))
	require.Nil(t, (&qdrant.Value{}).AsInterface())
	require.Equal(t, []any{}, (&qdrant.Value{Kind: &qdrant.Value_ListValue{}}).AsInterface())
	require.Equal(t, map[string]any{}, (&qdrant.Value{Kind: &qdrant.Value_StructValue{}}).AsInterface())
}

func TestNewValue_Errors(t *testing.T) {
	inputs := []any{
		struct{}{},
		string([]byte{0xff}),
		uint64(math.MaxUint64),
		[]uint64{math.MaxUint64},
		map[string]uint64{"k": math.MaxUint64},
		map[string]*qdrant.Value{string([]byte{0xff}): nil},
	}
	for _, input := range inputs {
		_, err := qdrant.NewValue(input)
		require.Error(t, err, "%T", input)
	}
	_, err := qdrant.TryValueMap(map[string]any{string([]byte{0xff}): 1})
	require.Error(t, err)
}
