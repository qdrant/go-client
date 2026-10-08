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
		{map[string]*qdrant.Value{"k": qdrant.NewValueString("v")}, qdrant.NewValueStruct(structValue)},
		{map[string]string{"k": "v"}, qdrant.NewValueStruct(structValue)},
		{map[string]int{"k": 1}, qdrant.NewValueStruct(&qdrant.Struct{Fields: map[string]*qdrant.Value{"k": one}})},
		{map[string]int64{"k": 1}, qdrant.NewValueStruct(&qdrant.Struct{Fields: map[string]*qdrant.Value{"k": one}})},
		{map[string]float64{"k": 1.5},
			qdrant.NewValueStruct(&qdrant.Struct{Fields: map[string]*qdrant.Value{"k": qdrant.NewValueDouble(1.5)}})},
		{map[string]bool{"k": true},
			qdrant.NewValueStruct(&qdrant.Struct{Fields: map[string]*qdrant.Value{"k": qdrant.NewValueBool(true)}})},
	}
	for _, tt := range tests {
		got, err := qdrant.NewValue(tt.input)
		require.NoError(t, err, "%T", tt.input)
		require.True(t, proto.Equal(tt.want, got), "%T: got %v, want %v", tt.input, got, tt.want)
	}
}

func TestValue_AsInterface(t *testing.T) {
	valNull := qdrant.NewValueNull()
	require.Nil(t, valNull.AsInterface())

	valBool := qdrant.NewValueBool(true)
	require.Equal(t, true, valBool.AsInterface())

	valInt := qdrant.NewValueInt(42)
	require.Equal(t, int64(42), valInt.AsInterface())

	valDouble := qdrant.NewValueDouble(3.14)
	require.InDelta(t, 3.14, valDouble.AsInterface(), 0)

	valString := qdrant.NewValueString("hello")
	require.Equal(t, "hello", valString.AsInterface())

	listVal := qdrant.NewValueFromList(qdrant.NewValueInt(1), qdrant.NewValueString("item"))
	require.Equal(t, []any{int64(1), "item"}, listVal.AsInterface())

	fields := map[string]*qdrant.Value{
		"int_field": qdrant.NewValueInt(100),
		"str_field": qdrant.NewValueString("world"),
		"list":      listVal,
	}
	structVal := qdrant.NewValueFromFields(fields)
	expectedMap := map[string]any{
		"int_field": int64(100),
		"str_field": "world",
		"list":      []any{int64(1), "item"},
	}
	require.Equal(t, expectedMap, structVal.AsInterface())
	require.Equal(t, expectedMap, qdrant.ValueMapToMap(fields))
}

func TestNewValue_Errors(t *testing.T) {
	inputs := []any{
		struct{}{},
		string([]byte{0xff}),
		uint64(math.MaxUint64),
		[]uint64{math.MaxUint64},
		map[string]uint64{"overflow": math.MaxUint64},
	}
	for _, input := range inputs {
		_, err := qdrant.NewValue(input)
		require.Error(t, err, "%T", input)
	}
}
