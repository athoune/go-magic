package ast

import (
	"errors"
	"strconv"
	"strings"
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

type Offset struct {
	Level     int
	Value     int64
	Dynamic   bool
	DynOffset int64
	DynType   byte
	DynAction byte
	DynArg    int64
}

type Compare struct {
	Operation   byte // = > < & ^ ~
	Not         bool // !
	StringValue string
	FloatValue  float64
	IntValue    int64
	Type        Clue
}

func ParseOffset(txt string) (*Offset, error) {
	offset := &Offset{}
	var err error
	var value string
	ooo := offset_re.FindStringSubmatch(txt) // looking for > prefix
	if len(ooo) == 2 {                       // no
		offset.Level = 0
		value = ooo[1]
	} else {
		offset.Level = len(ooo[1])
		value = ooo[2]
	}

	if strings.HasPrefix(value, "(") {
		offset.Dynamic = true
		dyn := dynamic_value_re.FindStringSubmatch(value[1 : len(value)-1])
		if len(dyn) != 5 {
			return nil, errors.New("Bad dynamic offset value : " + value)
		}
		offset.DynOffset, err = strconv.ParseInt(dyn[1], 0, 32)
		if err != nil {
			return nil, err
		}
		offset.DynType = dyn[2][0]
		offset.DynAction = dyn[3][0]
		offset.DynArg, err = strconv.ParseInt(dyn[4], 0, 32)
		if err != nil {
			return nil, err
		}
	} else {
		offset.Value, err = strconv.ParseInt(value, 0, 32)
		if err != nil {
			return nil, err
		}
	}
	return offset, nil
}

func ParseCompare(txt string, clue Clue) (*Compare, error) {
	if txt[0] == 'x' {
		return nil, nil
	}
	compare := &Compare{}
	var err error
	i := 0
	if txt[i] == '!' {
		compare.Not = true
		i++
	}
	compare.Operation = txt[i]
	not_implicit_equality := false
	for _, a := range []byte("=><&^~") {
		if compare.Operation == a {
			not_implicit_equality = true
			break
		}
	}
	if !not_implicit_equality {
		compare.Operation = '='
	} else {
		i++
	}
	compare.Type = clue
	switch {
	case clue == TYPE_CLUE_STRING:
		compare.StringValue = txt[i:]
	case clue == TYPE_CLUE_FLOAT:
		compare.FloatValue, err = strconv.ParseFloat(txt[i:], 64)
		if err != nil {
			return nil, err
		}
	case clue == TYPE_CLUE_INT:
		compare.IntValue, err = strconv.ParseInt(txt[i:], 0, 64)
		if err != nil {
			return nil, err
		}
	}
	return compare, nil
}
