package qdrant_test

import (
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestVectorOutputAccessorsPreferTypedVector(t *testing.T) {
	tests := []struct {
		name       string
		vector     *qdrant.VectorOutput
		wantDense  *qdrant.DenseVector
		wantSparse *qdrant.SparseVector
		wantMulti  *qdrant.MultiDenseVector
	}{
		{
			name: "legacy dense",
			vector: &qdrant.VectorOutput{
				Data: []float32{1, 2},
			},
			wantDense: &qdrant.DenseVector{Data: []float32{1, 2}},
		},
		{
			name: "typed dense takes precedence over legacy data",
			vector: &qdrant.VectorOutput{
				Data:   []float32{9},
				Vector: &qdrant.VectorOutput_Dense{Dense: &qdrant.DenseVector{Data: []float32{1, 2}}},
			},
			wantDense: &qdrant.DenseVector{Data: []float32{1, 2}},
		},
		{
			name: "typed sparse does not expose dense data",
			vector: &qdrant.VectorOutput{
				Data:    []float32{9},
				Indices: &qdrant.SparseIndices{Data: []uint32{9}},
				Vector:  &qdrant.VectorOutput_Sparse{Sparse: &qdrant.SparseVector{Values: []float32{1}, Indices: []uint32{2}}},
			},
			wantSparse: &qdrant.SparseVector{Values: []float32{1}, Indices: []uint32{2}},
		},
		{
			name: "typed multi dense takes precedence over legacy data",
			vector: &qdrant.VectorOutput{
				Data:         []float32{9, 9},
				VectorsCount: func() *uint32 { value := uint32(1); return &value }(),
				Vector:       &qdrant.VectorOutput_MultiDense{MultiDense: &qdrant.MultiDenseVector{Vectors: []*qdrant.DenseVector{{Data: []float32{1, 2}}}}},
			},
			wantMulti: &qdrant.MultiDenseVector{Vectors: []*qdrant.DenseVector{{Data: []float32{1, 2}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantDense, tt.vector.GetDenseVector())
			require.Equal(t, tt.wantSparse, tt.vector.GetSparseVector())
			require.Equal(t, tt.wantMulti, tt.vector.GetMultiVector())
		})
	}
}
