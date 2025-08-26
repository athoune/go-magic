package unpack

func operation[V int64 | uint64](current, expected V, operation byte) bool {
	switch operation {
	case '=':
		return current == expected
	}
	return false
}

func operationString(current, expected string, operation byte) bool {
	switch operation {
	case '=':
		return current == expected
	}
	return false
}
