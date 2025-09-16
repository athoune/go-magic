package parse

import (
	"strconv"
	"strings"

	"github.com/athoune/go-magic/model"
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
	switch type_.Clue {
	case model.TYPE_CLUE_STRING:
		compare.RawExpected, err = HandleStringEscape(value)
		if err != nil {
			return nil, poz + end, err
		}
	case model.TYPE_CLUE_INT:
		compare.RawExpected = value
		if strings.HasSuffix(compare.RawExpected, "h") ||
			strings.HasSuffix(compare.RawExpected, "L") {
			/*
				[FIXME]
				What the hell are this letters ?!
				 >>15	ulelong		!0x00010000h	\b, version %#8.8
				 0	lelong		0x1b031336L	Netboot image,
			*/
			compare.BinaryExpected, err = strconv.ParseUint(
				compare.RawExpected[:len(compare.RawExpected)-1], 0, 64)
			if err != nil {
				return nil, poz + end, err
			}
		} else {
			compare.BinaryExpected, err = strconv.ParseUint(compare.RawExpected, 0, 64)
			if err != nil {
				return nil, poz + end, err
			}
		}
	default:
		compare.RawExpected = value
	}
	//
	return compare, poz + end, nil
}
