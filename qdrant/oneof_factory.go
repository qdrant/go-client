// This file contains helper functions to create oneof gRPC fields.
// It is very verbose to create them using the struct literal syntax.
// So we do it for the users and export these helper functions.

package qdrant

import "google.golang.org/protobuf/types/known/timestamppb"

// Creates a *VectorsConfig instance from *VectorParams.
func NewVectorsConfig(params *VectorParams) *VectorsConfig {
	return &VectorsConfig{
		Config: &VectorsConfig_Params{
			Params: params,
		},
	}
}

// Creates a *VectorsConfig instance from a map of *VectorParams.
func NewVectorsConfigMap(paramsMap map[string]*VectorParams) *VectorsConfig {
	return &VectorsConfig{
		Config: &VectorsConfig_ParamsMap{
			ParamsMap: &VectorParamsMap{
				Map: paramsMap,
			},
		},
	}
}

// Creates a *SparseVectorConfig instance from a map of *SparseVectorParams.
func NewSparseVectorsConfig(paramsMap map[string]*SparseVectorParams) *SparseVectorConfig {
	return &SparseVectorConfig{
		Map: paramsMap,
	}
}

// Creates a *VectorsConfigDiff instance from *VectorParamsDiff.
func NewVectorsConfigDiff(params *VectorParamsDiff) *VectorsConfigDiff {
	return &VectorsConfigDiff{
		Config: &VectorsConfigDiff_Params{
			Params: params,
		},
	}
}

// Creates a *VectorsConfigDiff instance from a map of *VectorParamsDiff.
func NewVectorsConfigDiffMap(paramsMap map[string]*VectorParamsDiff) *VectorsConfigDiff {
	return &VectorsConfigDiff{
		Config: &VectorsConfigDiff_ParamsMap{
			ParamsMap: &VectorParamsDiffMap{
				Map: paramsMap,
			},
		},
	}
}

// Creates a *QuantizationConfig instance from *ScalarQuantization.
func NewQuantizationScalar(scalar *ScalarQuantization) *QuantizationConfig {
	return &QuantizationConfig{
		Quantization: &QuantizationConfig_Scalar{
			Scalar: scalar,
		},
	}
}

// Creates a *QuantizationConfig instance from *BinaryQuantization.
func NewQuantizationBinary(binary *BinaryQuantization) *QuantizationConfig {
	return &QuantizationConfig{
		Quantization: &QuantizationConfig_Binary{
			Binary: binary,
		},
	}
}

// Creates a *QuantizationConfig instance from *ProductQuantization.
func NewQuantizationProduct(product *ProductQuantization) *QuantizationConfig {
	return &QuantizationConfig{
		Quantization: &QuantizationConfig_Product{
			Product: product,
		},
	}
}

// Creates a *QuantizationConfigDiff instance from *ScalarQuantization.
func NewQuantizationDiffScalar(scalar *ScalarQuantization) *QuantizationConfigDiff {
	return &QuantizationConfigDiff{
		Quantization: &QuantizationConfigDiff_Scalar{
			Scalar: scalar,
		},
	}
}

// Creates a *QuantizationConfigDiff instance from *BinaryQuantization.
func NewQuantizationDiffBinary(binary *BinaryQuantization) *QuantizationConfigDiff {
	return &QuantizationConfigDiff{
		Quantization: &QuantizationConfigDiff_Binary{
			Binary: binary,
		},
	}
}

// Creates a *QuantizationConfigDiff instance from *ProductQuantization.
func NewQuantizationDiffProduct(product *ProductQuantization) *QuantizationConfigDiff {
	return &QuantizationConfigDiff{
		Quantization: &QuantizationConfigDiff_Product{
			Product: product,
		},
	}
}

// Creates a *QuantizationConfigDiff instance with quantization disabled.
func NewQuantizationDiffDisabled() *QuantizationConfigDiff {
	return &QuantizationConfigDiff{
		Quantization: &QuantizationConfigDiff_Disabled{
			Disabled: &Disabled{},
		},
	}
}

// Creates a *PayloadIndexParams instance from *KeywordIndexParams.
// This is an alias for NewPayloadIndexParamsKeyword().
func NewPayloadIndexParams(params *KeywordIndexParams) *PayloadIndexParams {
	return NewPayloadIndexParamsKeyword(params)
}

// Creates a *PayloadIndexParams instance from *KeywordIndexParams.
func NewPayloadIndexParamsKeyword(params *KeywordIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_KeywordIndexParams{
			KeywordIndexParams: params,
		},
	}
}

// Creates a *PayloadIndexParams instance from *IntegerIndexParams.
func NewPayloadIndexParamsInt(params *IntegerIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_IntegerIndexParams{
			IntegerIndexParams: params,
		},
	}
}

// Creates a *PayloadIndexParams instance from *FloatIndexParams.
func NewPayloadIndexParamsFloat(params *FloatIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_FloatIndexParams{
			FloatIndexParams: params,
		},
	}
}

// Creates a *PayloadIndexParams instance from *GeoIndexParams.
func NewPayloadIndexParamsGeo(params *GeoIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_GeoIndexParams{
			GeoIndexParams: params,
		},
	}
}

// Creates a *PayloadIndexParams instance from *TextIndexParams.
func NewPayloadIndexParamsText(params *TextIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_TextIndexParams{
			TextIndexParams: params,
		},
	}
}

// Creates a *PayloadIndexParams instance from *BoolIndexParams.
func NewPayloadIndexParamsBool(params *BoolIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_BoolIndexParams{
			BoolIndexParams: params,
		},
	}
}

// Creates a *PayloadIndexParams instance from *DatetimeIndexParams.
func NewPayloadIndexParamsDatetime(params *DatetimeIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_DatetimeIndexParams{
			DatetimeIndexParams: params,
		},
	}
}

// Creates a *PayloadIndexParams instance from *UuidIndexParams.
func NewPayloadIndexParamsUUID(params *UuidIndexParams) *PayloadIndexParams {
	return &PayloadIndexParams{
		IndexParams: &PayloadIndexParams_UuidIndexParams{
			UuidIndexParams: params,
		},
	}
}

// Creates an *AliasOperations instance to create an alias.
// This is an alias for NewAliasCreate().
func NewAlias(aliasName, collectionName string) *AliasOperations {
	return NewAliasCreate(aliasName, collectionName)
}

// Creates an *AliasOperations instance to create an alias.
func NewAliasCreate(aliasName, collectionName string) *AliasOperations {
	return &AliasOperations{
		Action: &AliasOperations_CreateAlias{
			CreateAlias: &CreateAlias{
				AliasName:      aliasName,
				CollectionName: collectionName,
			},
		},
	}
}

// Creates an *AliasOperations instance to rename an alias.
func NewAliasRename(oldAliasName, newAliasName string) *AliasOperations {
	return &AliasOperations{
		Action: &AliasOperations_RenameAlias{
			RenameAlias: &RenameAlias{
				OldAliasName: oldAliasName,
				NewAliasName: newAliasName,
			},
		},
	}
}

// Creates an *AliasOperations instance to delete an alias.
func NewAliasDelete(aliasName string) *AliasOperations {
	return &AliasOperations{
		Action: &AliasOperations_DeleteAlias{
			DeleteAlias: &DeleteAlias{
				AliasName: aliasName,
			},
		},
	}
}

// Creates a *ShardKey instance from a string.
// This is an alias for NewShardKeyKeyword().
func NewShardKey(key string) *ShardKey {
	return NewShardKeyKeyword(key)
}

// Creates a *ShardKey instance from a string.
func NewShardKeyKeyword(key string) *ShardKey {
	return &ShardKey{
		Key: &ShardKey_Keyword{
			Keyword: key,
		},
	}
}

// Creates a *ShardKey instance from a uint64.
func NewShardKeyNum(key uint64) *ShardKey {
	return &ShardKey{
		Key: &ShardKey_Number{
			Number: key,
		},
	}
}

// Creates a *ReadConsistency instance from a ReadConsistencyType.
// This is an alias for NewReadConsistencyType().
func NewReadConsistency(readConsistencyType ReadConsistencyType) *ReadConsistency {
	return NewReadConsistencyType(readConsistencyType)
}

// Creates a *ReadConsistency instance from a ReadConsistencyType.
func NewReadConsistencyType(readConsistencyType ReadConsistencyType) *ReadConsistency {
	return &ReadConsistency{
		Value: &ReadConsistency_Type{
			Type: readConsistencyType,
		},
	}
}

// Creates a *ReadConsistency instance from a factor in uint64.
func NewReadConsistencyFactor(readConsistencyFactor uint64) *ReadConsistency {
	return &ReadConsistency{
		Value: &ReadConsistency_Factor{
			Factor: readConsistencyFactor,
		},
	}
}

// Creates a *PointId instance from a UUID string.
// Same as NewIDUUID().
func NewID(uuid string) *PointId {
	return NewIDUUID(uuid)
}

// Creates a *PointId instance from a UUID string.
func NewIDUUID(uuid string) *PointId {
	return &PointId{
		PointIdOptions: &PointId_Uuid{
			Uuid: uuid,
		},
	}
}

// Creates a *PointId instance from a positive integer.
func NewIDNum(num uint64) *PointId {
	return &PointId{
		PointIdOptions: &PointId_Num{
			Num: num,
		},
	}
}

// Creates a *VectorInput instance for dense vectors.
// This is an alias for NewVectorInputDense().
func NewVectorInput(values ...float32) *VectorInput {
	return NewVectorInputDense(values)
}

// Creates a *VectorInput instance from a *PointId.
func NewVectorInputID(id *PointId) *VectorInput {
	return &VectorInput{
		Variant: &VectorInput_Id{
			Id: id,
		},
	}
}

// Creates a *VectorInput instance for dense vectors.
func NewVectorInputDense(vector []float32) *VectorInput {
	return &VectorInput{
		Variant: &VectorInput_Dense{
			Dense: &DenseVector{
				Data: vector,
			},
		},
	}
}

// Creates a *VectorInput instance for sparse vectors.
func NewVectorInputSparse(indices []uint32, values []float32) *VectorInput {
	return &VectorInput{
		Variant: &VectorInput_Sparse{
			Sparse: &SparseVector{
				Values:  values,
				Indices: indices,
			},
		},
	}
}

// Creates a *VectorInput instance for multi vectors.
func NewVectorInputMulti(vectors [][]float32) *VectorInput {
	var multiVec []*DenseVector
	for _, vector := range vectors {
		multiVec = append(multiVec, &DenseVector{
			Data: vector,
		})
	}
	return &VectorInput{
		Variant: &VectorInput_MultiDense{
			MultiDense: &MultiDenseVector{
				Vectors: multiVec,
			},
		},
	}
}

// Creates a *WithPayloadSelector instance with payload enabled/disabled.
// This is an alias for NewWithPayloadEnable().
func NewWithPayload(enable bool) *WithPayloadSelector {
	return NewWithPayloadEnable(enable)
}

// Creates a *WithPayloadSelector instance with payload enabled/disabled.
func NewWithPayloadEnable(enable bool) *WithPayloadSelector {
	return &WithPayloadSelector{
		SelectorOptions: &WithPayloadSelector_Enable{
			Enable: enable,
		},
	}
}

// Creates a *WithPayloadSelector instance with payload fields included.
func NewWithPayloadInclude(include ...string) *WithPayloadSelector {
	return &WithPayloadSelector{
		SelectorOptions: &WithPayloadSelector_Include{
			Include: &PayloadIncludeSelector{
				Fields: include,
			},
		},
	}
}

// Creates a *WithPayloadSelector instance with payload fields excluded.
func NewWithPayloadExclude(exclude ...string) *WithPayloadSelector {
	return &WithPayloadSelector{
		SelectorOptions: &WithPayloadSelector_Exclude{
			Exclude: &PayloadExcludeSelector{
				Fields: exclude,
			},
		},
	}
}

// Creates a *WithVectorsSelector instance with vectors enabled/disabled.
// This is an alias for NewWithVectorsEnable().
func NewWithVectors(enable bool) *WithVectorsSelector {
	return NewWithVectorsEnable(enable)
}

// Creates a *WithVectorsSelector instance with vectors enabled/disabled.
func NewWithVectorsEnable(enable bool) *WithVectorsSelector {
	return &WithVectorsSelector{
		SelectorOptions: &WithVectorsSelector_Enable{
			Enable: enable,
		},
	}
}

// Creates a *WithVectorsSelector instance with vectors included.
func NewWithVectorsInclude(names ...string) *WithVectorsSelector {
	return &WithVectorsSelector{
		SelectorOptions: &WithVectorsSelector_Include{
			Include: &VectorsSelector{
				Names: names,
			},
		},
	}
}

// Creates a *Vectors instance for dense vectors.
// This is an alias for NewVectorsDense().
func NewVectors(values ...float32) *Vectors {
	return NewVectorsDense(values)
}

// Creates a *Vectors instance for dense vectors.
func NewVectorsDense(vector []float32) *Vectors {
	return &Vectors{
		VectorsOptions: &Vectors_Vector{
			Vector: NewVectorDense(vector),
		},
	}
}

// Creates a *Vectors instance for sparse vectors.
func NewVectorsSparse(indices []uint32, values []float32) *Vectors {
	return &Vectors{
		VectorsOptions: &Vectors_Vector{
			Vector: NewVectorSparse(indices, values),
		},
	}
}

// Creates a *Vectors instance for multi vectors.
func NewVectorsMulti(vectors [][]float32) *Vectors {
	return &Vectors{
		VectorsOptions: &Vectors_Vector{
			Vector: NewVectorMulti(vectors),
		},
	}
}

// Creates a *Vectors instance for a map of named *Vector.
func NewVectorsMap(vectors map[string]*Vector) *Vectors {
	return &Vectors{
		VectorsOptions: &Vectors_Vectors{
			Vectors: &NamedVectors{
				Vectors: vectors,
			},
		},
	}
}

// Creates a *Vector instance for dense vectors.
// This is an alias for NewVectorDense().
func NewVector(values ...float32) *Vector {
	return NewVectorDense(values)
}

// Creates a *Vector instance for dense vectors.
func NewVectorDense(vector []float32) *Vector {
	return &Vector{
		Data: vector,
	}
}

// Creates a *Vector instance for sparse vectors.
func NewVectorSparse(indices []uint32, values []float32) *Vector {
	return &Vector{
		Data: values,
		Indices: &SparseIndices{
			Data: indices,
		},
	}
}

// Creates a *Vector instance for multi vectors.
func NewVectorMulti(vectors [][]float32) *Vector {
	vectorsCount := uint32(len(vectors))
	var flattenedVec []float32
	for _, vector := range vectors {
		flattenedVec = append(flattenedVec, vector...)
	}
	return &Vector{
		Data:         flattenedVec,
		VectorsCount: &vectorsCount,
	}
}

// Creates a *StartFrom instance for a float value.
func NewStartFromFloat(value float64) *StartFrom {
	return &StartFrom{
		Value: &StartFrom_Float{
			Float: value,
		},
	}
}

// Creates a *StartFrom instance for an integer value.
func NewStartFromInt(value int64) *StartFrom {
	return &StartFrom{
		Value: &StartFrom_Integer{
			Integer: value,
		},
	}
}

// Creates a *StartFrom instance for a timestamp value.
// Parameters:
// seconds: Represents seconds of UTC time since Unix epoch 1970-01-01T00:00:00Z.
// Must be from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
// nanos: Non-negative fractions of a second at nanosecond resolution.
// Negative second values with fractions must still have non-negative
// nanos values that count forward in time.
// Must be from 0 to 999,999,999 inclusive.
func NewStartFromTimestamp(seconds int64, nanos int32) *StartFrom {
	return &StartFrom{
		Value: &StartFrom_Timestamp{
			Timestamp: &timestamppb.Timestamp{
				Seconds: seconds,
				Nanos:   nanos,
			},
		},
	}
}

// Creates a *StartFrom instance for a datetime string in the RFC3339 format.
func NewStartFromDatetime(value string) *StartFrom {
	return &StartFrom{
		Value: &StartFrom_Datetime{
			Datetime: value,
		},
	}
}

// Creates a *TargetVector instance from a *Vector.
// This is an alias for NewTargetVector().
func NewTarget(vector *Vector) *TargetVector {
	return NewTargetVector(vector)
}

// Creates a *TargetVector instance from a *Vector.
func NewTargetVector(vector *Vector) *TargetVector {
	return &TargetVector{
		Target: &TargetVector_Single{
			Single: &VectorExample{
				Example: &VectorExample_Vector{
					Vector: vector,
				},
			},
		},
	}
}

// Creates a *TargetVector instance from a *PointId.
func NewTargetID(id *PointId) *TargetVector {
	return &TargetVector{
		Target: &TargetVector_Single{
			Single: &VectorExample{
				Example: &VectorExample_Id{
					Id: id,
				},
			},
		},
	}
}

// Creates a *Query instance for a nearest query from *VectorInput.
func NewQueryNearest(nearest *VectorInput) *Query {
	return &Query{
		Variant: &Query_Nearest{
			Nearest: nearest,
		},
	}
}

// Creates a *Query instance for a nearest query from dense vectors.
// This is an alias for NewQueryDense().
func NewQuery(values ...float32) *Query {
	return NewQueryDense(values)
}

// Creates a *Query instance for a nearest query from dense vectors.
func NewQueryDense(vector []float32) *Query {
	return NewQueryNearest(NewVectorInputDense(vector))
}

// Creates a *Query instance for a nearest query from sparse vectors.
func NewQuerySparse(indices []uint32, values []float32) *Query {
	return NewQueryNearest(NewVectorInputSparse(indices, values))
}

// Creates a *Query instance for a nearest query from multi vectors.
func NewQueryMulti(vectors [][]float32) *Query {
	return NewQueryNearest(NewVectorInputMulti(vectors))
}

// Creates a *Query instance for a nearest query from *PointId.
func NewQueryID(id *PointId) *Query {
	return NewQueryNearest(NewVectorInputID(id))
}

// Creates a *Query instance for recommend query from *RecommendInput.
func NewQueryRecommend(recommend *RecommendInput) *Query {
	return &Query{
		Variant: &Query_Recommend{
			Recommend: recommend,
		},
	}
}

// Creates a *Query instance for a discover query from *DiscoverInput.
func NewQueryDiscover(discover *DiscoverInput) *Query {
	return &Query{
		Variant: &Query_Discover{
			Discover: discover,
		},
	}
}

// Creates a *Query instance for a context query from *ContextInput.
func NewQueryContext(context *ContextInput) *Query {
	return &Query{
		Variant: &Query_Context{
			Context: context,
		},
	}
}

// Creates a *Query instance for ordering points with *OrderBy.
func NewQueryOrderBy(orderBy *OrderBy) *Query {
	return &Query{
		Variant: &Query_OrderBy{
			OrderBy: orderBy,
		},
	}
}

// Creeates a *Query instance for combining prefetch results with Fusion.
func NewQueryFusion(fusion Fusion) *Query {
	return &Query{
		Variant: &Query_Fusion{
			Fusion: fusion,
		},
	}
}

// Creates a *Query instance for sampling points.
func NewQuerySample(sample Sample) *Query {
	return &Query{
		Variant: &Query_Sample{
			Sample: sample,
		},
	}
}

// Creates a *FacetValue instance from a string.
func NewFacetValue(value string) *FacetValue {
	return &FacetValue{
		Variant: &FacetValue_StringValue{
			StringValue: value,
		},
	}
}

// Creates a *PointsUpdateOperation instance for upserting points.
func NewPointsUpdateUpsert(upsert *PointsUpdateOperation_PointStructList) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_Upsert{
			Upsert: upsert,
		},
	}
}

// Creates a *PointsUpdateOperation instance for setting payload.
func NewPointsUpdateSetPayload(setPayload *PointsUpdateOperation_SetPayload) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_SetPayload_{
			SetPayload: setPayload,
		},
	}
}

// Creates a *PointsUpdateOperation instance for overwriting payload.
func NewPointsUpdateOverwritePayload(overwritePayload *PointsUpdateOperation_OverwritePayload) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_OverwritePayload_{
			OverwritePayload: overwritePayload,
		},
	}
}

// Creates a *PointsUpdateOperation instance for deleting payload fields.
func NewPointsUpdateDeletePayload(deletePayload *PointsUpdateOperation_DeletePayload) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_DeletePayload_{
			DeletePayload: deletePayload,
		},
	}
}

// Creates a *PointsUpdateOperation instance for updating vectors.
func NewPointsUpdateUpdateVectors(updateVectors *PointsUpdateOperation_UpdateVectors) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_UpdateVectors_{
			UpdateVectors: updateVectors,
		},
	}
}

// Creates a *PointsUpdateOperation instance for deleting vectors.
func NewPointsUpdateDeleteVectors(deleteVectors *PointsUpdateOperation_DeleteVectors) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_DeleteVectors_{
			DeleteVectors: deleteVectors,
		},
	}
}

// Creates a *PointsUpdateOperation instance for deleting points.
func NewPointsUpdateDeletePoints(deletePoints *PointsUpdateOperation_DeletePoints) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_DeletePoints_{
			DeletePoints: deletePoints,
		},
	}
}

// Creates a *PointsUpdateOperation instance for clearing payload.
func NewPointsUpdateClearPayload(clearPayload *PointsUpdateOperation_ClearPayload) *PointsUpdateOperation {
	return &PointsUpdateOperation{
		Operation: &PointsUpdateOperation_ClearPayload_{
			ClearPayload: clearPayload,
		},
	}
}

// Creates a *OrderValue instance from an integer.
func NewOrderValueInt(value int64) *OrderValue {
	return &OrderValue{
		Variant: &OrderValue_Int{
			Int: value,
		},
	}
}

// Creates a *OrderValue instance from a float.
func NewOrderValueFloat(value float64) *OrderValue {
	return &OrderValue{
		Variant: &OrderValue_Float{
			Float: value,
		},
	}
}

// Creates a *OrderValue instance from an unsigned integer.
func NewGroupIDUnsigned(value uint64) *GroupId {
	return &GroupId{
		Kind: &GroupId_UnsignedValue{
			UnsignedValue: value,
		},
	}
}

// Creates a *GroupId instance from an integer.
func NewGroupIDInt(value int64) *GroupId {
	return &GroupId{
		Kind: &GroupId_IntegerValue{
			IntegerValue: value,
		},
	}
}

// Creates a *GroupId instance from a string.
func NewGroupIDString(value string) *GroupId {
	return &GroupId{
		Kind: &GroupId_StringValue{
			StringValue: value,
		},
	}
}

// Creates a *PointsSelector instance for selecting points by IDs.
// This is an alias for NewPointsSelectorIDs().
func NewPointsSelector(ids ...*PointId) *PointsSelector {
	return NewPointsSelectorIDs(ids)
}

// Creates a *PointsSelector instance for selecting points by filter.
func NewPointsSelectorFilter(filter *Filter) *PointsSelector {
	return &PointsSelector{
		PointsSelectorOneOf: &PointsSelector_Filter{
			Filter: filter,
		},
	}
}

// Creates a *PointsSelector instance for selecting points by IDs.
func NewPointsSelectorIDs(ids []*PointId) *PointsSelector {
	return &PointsSelector{
		PointsSelectorOneOf: &PointsSelector_Points{
			Points: &PointsIdsList{
				Ids: ids,
			},
		},
	}
}

// Creates a pointer to a value of any type.
func PtrOf[T any](t T) *T {
	return &t
}
