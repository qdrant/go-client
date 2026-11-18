package qdrant_test

import (
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/assert"
)

func TestGetDenseVectorDoesNotReinterpretLegacySparseOrMulti(t *testing.T) {
	// Legacy sparse vectors put their values in the deprecated data field and use
	// indices as the discriminator. They must not be returned as a dense vector.
	sparse := &qdrant.VectorOutput{
		Data:    []float32{0.1, 0.2},
		Indices: &qdrant.SparseIndices{Data: []uint32{1, 5}},
	}
	assert.Nil(t, sparse.GetDenseVector())
	assert.NotNil(t, sparse.GetSparseVector())

	// Legacy multi vectors flatten their values into data and use vectors_count as
	// the discriminator. They must not be returned as a dense vector either.
	vectorsCount := uint32(2)
	multi := &qdrant.VectorOutput{
		Data:         []float32{0.1, 0.2, 0.3, 0.4},
		VectorsCount: &vectorsCount,
	}
	assert.Nil(t, multi.GetDenseVector())
	assert.NotNil(t, multi.GetMultiVector())

	// A legacy dense vector has data without either discriminator and must still work.
	dense := &qdrant.VectorOutput{Data: []float32{0.1, 0.2}}
	got := dense.GetDenseVector()
	assert.NotNil(t, got)
	assert.Equal(t, []float32{0.1, 0.2}, got.GetData())

	// The non-deprecated oneof form must keep working.
	oneof := &qdrant.VectorOutput{
		Vector: &qdrant.VectorOutput_Dense{Dense: &qdrant.DenseVector{Data: []float32{0.3, 0.4}}},
	}
	assert.Equal(t, []float32{0.3, 0.4}, oneof.GetDenseVector().GetData())
}
