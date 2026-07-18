package operation

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/athoune/go-magic/model"
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

func compareStringWithFlags(a, b string, opts model.StringOptions) bool {
	for i := 0; i < len(b); i++ {
		ac := a[i]
		bc := b[i]

		if opts&model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER != 0 && unicode.IsLower(rune(ac)) {
			if unicode.ToLower(rune(bc)) != rune(ac) {
				return false
			}
		} else if opts&model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER != 0 && unicode.IsUpper(rune(ac)) {
			if unicode.ToUpper(rune(bc)) != rune(ac) {
				return false
			}
		} else {
			if ac != bc {
				return false
			}
		}
	}
	return true
}

func CompareString(a, b string, operator byte, opts model.StringOptions) (bool, error) {
	switch operator {
	case '=':
		// C semantics: the file buffer is zero-padded to match the pattern length.
		// file_strncmp compares exactly len(b) bytes; missing bytes are treated as 0.
		if len(a) < len(b) {
			a += strings.Repeat("\x00", len(b)-len(a))
		}

		var matched bool
		if opts&model.STRING_OPTIONS_COMPACT_WITH_SPACES != 0 {
			// TODO: implement compact whitespace matching
			matched = a[:len(b)] == b
		} else {
			matched = compareStringWithFlags(a, b, opts)
		}

		if matched && opts&model.STRING_OPTIONS_FULL_WORD != 0 {
			// Check that the next character after the match is a space or end of string
			if len(a) > len(b) && !unicode.IsSpace(rune(a[len(b)])) {
				matched = false
			}
		}
		return matched, nil
	}
	return false, fmt.Errorf("unknown operator: %v", operator)
}
