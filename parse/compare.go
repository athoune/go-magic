package parse

import (
	"github.com/athoune/go-magic/model"
	"github.com/athoune/go-magic/unpack"
)

const (
	COMPARE_LESS    = uint8(60)
	COMPARE_GREATER = uint8(62)
	COMPARE_EQUAL   = uint8(61)
	COMPARE_AND     = uint8(38)  // &
	COMPARE_OR      = uint8(94)  // ^
	COMPARE_NEGATED = uint8(126) // ~
	COMPARE_NOT     = uint8(33)  // !
)

// ParseCompare extract the operation, the value (typed) and the new position
func ParseCompare(line string, type_ *model.Type) (*model.Compare, int, error) {
	compare := &model.Compare{
		Type: type_,
	}
	if line[0] == 'x' {
		compare.X = true
		compare.Expected = nil
		compare.Comparator = 0
		return compare, 1, nil
	}
	poz := 0
	if type_.Name == "name" {
		end := notSpace(line)
		compare.RawExpected = line[:end]
		return compare, end, nil
	}
	var err error

	// Not
	if line[poz] == '!' {
		compare.Not = true
		poz++
	}

	// Operation
	compare.Comparator = line[poz]
	if !IsOperation(compare.Comparator) {
		compare.Comparator = '='
	} else {
		poz++
	}
	if line[poz] == ' ' {
		poz++
	}
	end := HandleSpaceEscape(line[poz:])

	// Value
	value := line[poz : poz+end]
	compare.Expected, err = unpack.BuildValueFromString(type_, value)
	if err != nil {
		return nil, poz + end, err
	}
	compare.RawExpected = value
	return compare, poz + end, nil
}
