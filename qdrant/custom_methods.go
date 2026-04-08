package qdrant

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func (x *Condition) UnmarshalJSON(data []byte) error {
	var temp struct {
		Match   json.RawMessage `json:"match"`
		Range   json.RawMessage `json:"range"`
		IsEmpty json.RawMessage `json:"isEmpty"`
		GeoRadius  json.RawMessage `json:"geo_radius"`
		HasId   json.RawMessage `json:"hasId"`
		Should  json.RawMessage `json:"should"`
		Must    json.RawMessage `json:"must"`
		MustNot json.RawMessage `json:"must_not"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if len(temp.Match) > 0 || len(temp.Range) > 0 || len(temp.GeoRadius) > 0 || len(temp.IsEmpty) > 0 || len(temp.HasId) > 0{
		var field FieldCondition

		if err := json.Unmarshal(data, &field); err != nil {
			return err
		}

		x.ConditionOneOf = &Condition_Field{
			Field: &field,
		}
	} else if len(temp.Should) > 0 || len(temp.Must) > 0 || len(temp.MustNot) > 0 {
		var innerCondition Filter
		if err := json.Unmarshal(data, &innerCondition); err != nil {
			return err
		}
		x.ConditionOneOf = &Condition_Filter{
			Filter: &innerCondition,
		}
	}

	return nil
}

func (x *Match) createMatchStruct(value any) (isMatch_MatchValue, error) {
	var matchValue isMatch_MatchValue
	rt := reflect.TypeOf(value)
	switch rt.Kind() {
	case reflect.Float64:
		var m Match_Integer
		m.Integer = int64(value.(float64))
		matchValue = &m
	case reflect.String:
		var m Match_Keyword
		m.Keyword = value.(string)
		matchValue = &m
	case reflect.Bool:
		var m Match_Boolean
		m.Boolean = value.(bool)
		matchValue = &m
	default:
		return nil, fmt.Errorf("cannot find match struct")
	}
	return matchValue, nil
}

func (x *Match) UnmarshalJSON(data []byte) error {
	var temp struct {
		Any    any `json:"any"`
		Value  any `json:"value"`
		Except any `json:"except"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Value != nil {
		m, err := x.createMatchStruct(temp.Value)
		if err != nil {
			return err
		}
		x.MatchValue = m
	} else if temp.Except != nil {
		rt := reflect.TypeOf(temp.Except)
		switch rt.Kind() {
		case reflect.Slice:
			except := temp.Except.([]any)
			if len(except) == 0 {
				x.MatchValue = &Match_ExceptIntegers{
					ExceptIntegers: &RepeatedIntegers{},
				}
				break
			}
			itemKind := reflect.TypeOf(except[0]).Kind()
			var m isMatch_MatchValue
			switch itemKind {
			case reflect.Float64:
				m = &Match_ExceptIntegers{
					ExceptIntegers: &RepeatedIntegers{},
				}
			case reflect.String:
				m = &Match_ExceptKeywords{
					ExceptKeywords: &RepeatedStrings{},
				}
			default:
				return fmt.Errorf("only integers and keywords allowed in except %s given", rt.Kind())
			}
			for i := 1; i < len(except); i++ {
				if itemKind != reflect.TypeOf(except[i]).Kind() {
					return fmt.Errorf("all array elements should be of the same type")
				}
			}
			switch itemKind {
			case reflect.Float64:
				items := make([]int64, 0)
				for _, y := range except {
					items = append(items, int64(y.(float64)))
				}
				m.(*Match_ExceptIntegers).ExceptIntegers.Integers = items
			case reflect.String:
				items := make([]string, 0)
				for _, y := range except {
					items = append(items, y.(string))
				}
				m.(*Match_ExceptKeywords).ExceptKeywords.Strings = items
			}
			x.MatchValue = m
		default:
			return fmt.Errorf("cannot find match struct")
		}
	} else if temp.Any != nil {
		rt := reflect.TypeOf(temp.Any)
		switch rt.Kind() {
		case reflect.Slice:
			anyValues := temp.Any.([]any)
			if len(anyValues) == 0 {
				x.MatchValue = &Match_Integers{
					Integers: &RepeatedIntegers{},
				}
				break
			}
			itemKind := reflect.TypeOf(anyValues[0]).Kind()
			var m isMatch_MatchValue
			switch itemKind {
			case reflect.Float64:
				m = &Match_Integers{
					Integers: &RepeatedIntegers{},
				}
			case reflect.String:
				m = &Match_Keywords{
					Keywords: &RepeatedStrings{},
				}
			default:
				return fmt.Errorf("only integers and keywords allowed in any %s given", rt.Kind())
			}
			for i := 1; i < len(anyValues); i++ {
				if itemKind != reflect.TypeOf(anyValues[i]).Kind() {
					return fmt.Errorf("all array elements should be of the same type")
				}
			}
			switch itemKind {
			case reflect.Float64:
				items := make([]int64, 0)
				for _, y := range anyValues {
					items = append(items, int64(y.(float64)))
				}
				m.(*Match_Integers).Integers.Integers = items
			case reflect.String:
				items := make([]string, 0)
				for _, y := range anyValues {
					items = append(items, y.(string))
				}
				m.(*Match_Keywords).Keywords.Strings = items
			}
			x.MatchValue = m
		default:
			return fmt.Errorf("cannot find match struct")
		}
	}

	return nil
}
