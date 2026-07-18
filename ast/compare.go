package ast

import (
	"fmt"

	"github.com/athoune/go-magic/model"
	"github.com/athoune/go-magic/operation"
)

func CompareValues(current, expected *model.Value, comparator byte) (bool, error) {
	switch current.Family {
	case model.TYPE_FAMILY_UINT:
		return operation.CompareNumbers(current.UIntValue, expected.UIntValue, comparator)
	case model.TYPE_FAMILY_INT:
		return operation.CompareNumbers(current.IntValue, expected.IntValue, comparator)
	case model.TYPE_FAMILY_FLOAT:
		return operation.CompareNumbers(current.FloatValue, expected.FloatValue, comparator)
	case model.TYPE_FAMILY_STRING:
		return operation.CompareString(current.StringValue, expected.StringValue, comparator, model.STRING_OPTIONS_NONE)
	default:
		return false, fmt.Errorf("unknown type family: %v", current.Family)
	}
}
