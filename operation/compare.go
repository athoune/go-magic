package operation

import (
	"fmt"
	"strings"
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
		// C semantics: the file buffer is zero-padded to match the pattern length.
		// file_strncmp compares exactly len(b) bytes; missing bytes are treated as 0.
		if len(a) < len(b) {
			a += strings.Repeat("\x00", len(b)-len(a))
		}
		return a[:len(b)] == b, nil
	}
	return false, fmt.Errorf("unknown operator: %v", operator)
}
