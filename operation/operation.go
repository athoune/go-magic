package operation

import (
	"fmt"
)

func Operate[v int | int64 | uint64](a, b v, operator byte) (v, error) {
	switch operator {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		return a / b, nil
	case '&':
		return a & b, nil
	case '|':
		return a | b, nil
	case '%':
		return a % b, nil
	case ' ': // Do nothing
		return a, nil
	}
	return v(0), fmt.Errorf("unknown operator: %v", operator)
}

func OperateString(a, arg string, operator byte) (string, error) {
	return a, nil // [FIXME]
}
