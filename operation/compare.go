package operation

import (
	"fmt"
)

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
