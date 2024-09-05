// This file contains constructor functions for creating conditions to be used in filters.
// https://qdrant.tech/documentation/concepts/filtering/#filtering-conditions

package qdrant

// Creates a condition that matches an exact keyword in a specified field.
// This is an alias for NewMatchKeyword().
// See: https://qdrant.tech/documentation/concepts/filtering/#match
func NewMatch(field, keyword string) *Condition {
	return NewMatchKeyword(field, keyword)
}

// Creates a condition that matches an exact keyword in a specified field.
// See: https://qdrant.tech/documentation/concepts/filtering/#match
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

// Creates a condition to match a specific substring, token or phrase.
// Exact texts that will match the condition depend on full-text index configuration.
// See: https://qdrant.tech/documentation/concepts/filtering/#full-text-match
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

// Creates a condition that matches a boolean value in a specified field.
// See: https://qdrant.tech/documentation/concepts/filtering/#match
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

// Creates a condition that matches an integer value in a specified field.
// See: https://qdrant.tech/documentation/concepts/filtering/#match
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

// Creates a condition that matches any of the given keywords in a specified field.
// See: https://qdrant.tech/documentation/concepts/filtering/#match-any
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

// Creates a condition that matches any of the given integer values in a specified field.
// See: https://qdrant.tech/documentation/concepts/filtering/#match-any
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

// Creates a condition that matches any value except the given keywords in a specified field.
// This is an alias for NewMatchExceptKeywords.
// See: https://qdrant.tech/documentation/concepts/filtering/#match-except
func NewMatchExcept(field string, keywords ...string) *Condition {
	return NewMatchExceptKeywords(field, keywords...)
}

// Creates a condition that matches any value except the given keywords in a specified field.
// See: https://qdrant.tech/documentation/concepts/filtering/#match-except
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

// Creates a condition that matches any value except the given integer values in a specified field.
// See: https://qdrant.tech/documentation/concepts/filtering/#match-except
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

// Creates a condition that checks if a specified field is null.
// See: https://qdrant.tech/documentation/concepts/filtering/#is-null
func NewIsNull(field string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_IsNull{
			IsNull: &IsNullCondition{
				Key: field,
			},
		},
	}
}

// Creates a condition that checks if a specified field is empty.
// See: https://qdrant.tech/documentation/concepts/filtering/#is-empty
func NewIsEmpty(field string) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_IsEmpty{
			IsEmpty: &IsEmptyCondition{
				Key: field,
			},
		},
	}
}

// Creates a condition that checks if a point has any of the specified IDs.
// See: https://qdrant.tech/documentation/concepts/filtering/#has-id
func NewHasID(ids ...*PointId) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_HasId{
			HasId: &HasIdCondition{
				HasId: ids,
			},
		},
	}
}

// Creates a nested condition for filtering on nested fields.
// See: https://qdrant.tech/documentation/concepts/filtering/#nested
func NewNestedCondition(field string, conditon *Condition) *Condition {
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

// Creates a nested filter for filtering on nested fields.
// See: https://qdrant.tech/documentation/concepts/filtering/#nested
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

// Creates a condition from a filter.
// This is useful for creating complex nested conditions.
func NewFilterAsCondition(filter *Filter) *Condition {
	return &Condition{
		ConditionOneOf: &Condition_Filter{
			Filter: filter,
		},
	}
}

// Creates a range condition for numeric or date fields.
// See: https://qdrant.tech/documentation/concepts/filtering/#range
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

// Creates a geo filter condition that matches points within a specified radius from a center point.
// See: https://qdrant.tech/documentation/concepts/filtering/#geo-radius
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

// Creates a geo filter condition that matches points within a specified bounding box.
// See: https://qdrant.tech/documentation/concepts/filtering/#geo-bounding-box
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

// Creates a geo filter condition that matches points within a specified polygon.
// See: https://qdrant.tech/documentation/concepts/filtering/#geo-polygon
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

// Creates a condition that filters based on the number of values in an array field.
// See: https://qdrant.tech/documentation/concepts/filtering/#values-count
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

// Creates a condition that filters based on a datetime range.
// See: https://qdrant.tech/documentation/concepts/filtering/#datetime-range
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
