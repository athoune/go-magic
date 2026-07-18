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
		// Like C's file_strncmp: compare only the first len(expected) bytes.
		// The file may contain trailing nulls or extra data after the match.
		if len(b) > len(a) {
			return a == b[:len(a)], nil
		}
		return a[:len(b)] == b, nil
	}
	return false, fmt.Errorf("unknown operator: %v", operator)
}
