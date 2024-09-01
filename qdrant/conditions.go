// This file contains constructor functions for creating conditions to be used in filters.
// https://qdrant.tech/documentation/concepts/filtering/#filtering-conditions
package qdrant

// TODO: please add comments to all filters below
func NewMatchKeyword(field, keyword string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_Keyword{Keyword: keyword},
				},
			},
		},
	}
}

func NewMatchText(field, text string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_Text{Text: text},
				},
			},
		},
	}
}

func NewMatchBool(field string, value bool) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_Boolean{Boolean: value},
				},
			},
		},
	}
}

func NewMatchInt(field string, value int64) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_Integer{Integer: value},
				},
			},
		},
	}
}

func NewMatchKeywords(field string, keywords ...string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_Keywords{Keywords: &RepeatedStrings{
						Strings: keywords,
					}},
				},
			},
		},
	}
}

func NewMatchInts(field string, values ...int64) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_Integers{Integers: &RepeatedIntegers{
						Integers: values,
					}},
				},
			},
		},
	}
}

// Same as NewMatchExceptKeywords.
func NewMatchExcept(field string, keywords ...string) *Condition {
	return NewMatchExceptKeywords(field, keywords...)
}

func NewMatchExceptKeywords(field string, keywords ...string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_ExceptKeywords{ExceptKeywords: &RepeatedStrings{
						Strings: keywords,
					}},
				},
			},
		},
	}
}

func NewMatchExceptInts(field string, values ...int64) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				Match: &Match{
					MatchValue: &Match_ExceptIntegers{ExceptIntegers: &RepeatedIntegers{
						Integers: values,
					}},
				},
			},
		},
	}
}

func NewIsNull(field string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_IsNull{
			IsNull: &IsNullCondition{
				Key: field,
			},
		},
	}
}

func NewIsEmpty(field string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_IsEmpty{
			IsEmpty: &IsEmptyCondition{
				Key: field,
			},
		},
	}
}

func NewHasID(ids ...*PointId) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_HasId{
			HasId: &HasIdCondition{
				HasId: ids,
			},
		},
	}
}

func NewNestedCondtion(field string, conditon *Condition) *Condition {
	filter := &Filter{
		Must: []*Condition{conditon},
	}

	return &Condition{
		ConditionOneOf: &Condition_Nested{
			Nested: &NestedCondition{
				Key:    field,
				Filter: filter,
			},
		},
	}
}

func NewNestedFilter(field string, filter *Filter) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Nested{
			Nested: &NestedCondition{
				Key:    field,
				Filter: filter,
			},
		},
	}
}

func NewFilterAsCondition(filter *Filter) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Filter{
			Filter: filter,
		},
	}
}

func NewRange(field string, rangeVal *Range) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key:   field,
				Range: rangeVal,
			},
		},
	}
}

func NewGeoRadius(field string, lat, long float64, radius float32) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				GeoRadius: &GeoRadius{
					Radius: radius,
					Center: &GeoPoint{
						Lat: lat,
						Lon: long,
					},
				},
			},
		},
	}
}

func NewGeoBoundingBox(field string, topLeftLat, topLeftLon, bottomRightLat, bottomRightLon float64) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key: field,
				GeoBoundingBox: &GeoBoundingBox{
					TopLeft: &GeoPoint{
						Lat: topLeftLat,
						Lon: topLeftLon,
					},
					BottomRight: &GeoPoint{
						Lat: bottomRightLat,
						Lon: bottomRightLon,
					},
				},
			},
		},
	}
}

func NewGeoPolygon(field string, exterior *GeoLineString, interior ...*GeoLineString) *Condition {
	geoPolygon := &GeoPolygon{
		Exterior:  exterior,
		Interiors: interior,
	}

	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key:        field,
				GeoPolygon: geoPolygon,
			},
		},
	}
}

func NewValuesCount(field string, valuesCount *ValuesCount) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key:         field,
				ValuesCount: valuesCount,
			},
		},
	}
}

func NewDatetimeRange(field string, dateTimeRange *DatetimeRange) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Field{
			Field: &FieldCondition{
				Key:           field,
				DatetimeRange: dateTimeRange,
			},
		},
	}
}
