package compare

import (
	"fmt"

	"github.com/athoune/go-magic/model"
)

func CompareValues(current, expected *model.Value, comparator byte) (bool, error) {
	switch current.Family {
	case model.TYPE_CLUE_UINT:
		return CompareNumbers(current.UIntValue, expected.UIntValue, comparator)
	case model.TYPE_CLUE_INT:
		return CompareNumbers(current.IntValue, expected.IntValue, comparator)
	case model.TYPE_CLUE_FLOAT:
		return CompareNumbers(current.FloatValue, expected.FloatValue, comparator)
	case model.TYPE_CLUE_STRING:
		return CompareString(current.StringValue, expected.StringValue, comparator)
	default:
		return false, fmt.Errorf("unknown type family: %v", current.Family)
	}
}

// = > < & ^ ~
func CompareNumbers[K int64 | uint64 | float64](a, b K, comparator byte) (bool, error) {
	switch comparator {
	case '=':
		return a == b, nil
	case '<':
		return a < b, nil
	case '>':
		return a > b, nil
	default:
		return false, fmt.Errorf("unknown comparator: %v", comparator)
	}
}

func CompareString(a, b string, operator byte) (bool, error) {
	switch operator {
	case '=':
		return a == b, nil
	}
	return false, fmt.Errorf("unknown operator: %v", operator)
}
