package ast

import (
	"errors"
	"fmt"
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
	Level       int
	Relative    bool
	Value       int64
	Dynamic     bool
	DynOffset   int64
	DynType     byte
	DynOperator byte
	DynArg      int64
}

type Compare struct {
	Operation   byte // = > < & ^ ~
	Not         bool // !
	StringValue string
	FloatValue  float64
	IntValue    int64
	QuadValue   []int64
	Type        Clue
}

func ParseOffset(offset *Offset, line string) error {
	if line == "" {
		return errors.New("empty value")
	}
	var err error

	offset.Level = 0
	//var err error
	for i := 0; i < len(line); i++ {
		if line[i] != '>' {
			break
		}
		offset.Level++
	}

	poz := offset.Level
	if line[poz] == '&' {
		offset.Relative = true
		poz++
	}
	switch {
	case line[poz] == '(':
		i := strings.IndexByte(line[poz+1:], ')')
		if i == -1 {
			return fmt.Errorf("can't find ')' in %s", line)
		}
		err = ParseDynamicOffset(offset, line[poz+1:poz+i+1])
		if err != nil {
			return err
		}
		poz += i
	default:
		offset.Value, err = strconv.ParseInt(line[poz:], 0, 32)
		if err != nil {
			return fmt.Errorf("can't parse int in [%s] at %v", line, poz)
		}
	}
	return nil
}

func ParseDynamicOffset(offset *Offset, line string) error {
	offset.Dynamic = true
	var err error
	dyn := dynamic_value_re.FindStringSubmatch(line)
	if len(dyn) == 0 {
		return errors.New("Bad dynamic offset value : " + line)
	}
	offset.DynOffset, err = strconv.ParseInt(dyn[value_dynamic_idx], 0, 32)
	if err != nil {
		return fmt.Errorf("%s <= %v", line, err)
	}
	if dyn[type_dynamic_idx] != "" {
		offset.DynType = dyn[type_dynamic_idx][0]
	}
	operator := dyn[operator_dynamic_idx]
	if operator != "" {
		offset.DynOperator = operator[0]
	}
	arg := dyn[arg_dynamic_idx]
	if arg != "" {
		offset.DynArg, err = strconv.ParseInt(arg, 0, 32)
		if err != nil {
			return fmt.Errorf("can't parse int in %v", arg)
		}
	}
	return nil
}

func ParseCompare(line string, clue Clue) (*Compare, error) {
	if len(line) == 0 {
		return nil, errors.New("empty compare value")
	}
	if line[0] == 'x' {
		return nil, nil
	}
	compare := &Compare{
		Type: clue,
	}
	var err error
	poz := 0
	if line[poz] == '!' {
		compare.Not = true
		poz++
	}
	compare.Operation = line[poz]
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
		poz++
	}
	value := strings.TrimLeft(line[poz:], " ")
	switch {
	case clue == TYPE_CLUE_STRING:
		compare.StringValue = value
	case clue == TYPE_CLUE_FLOAT:
		compare.FloatValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse float: %v in [%v]", value, line)
		}
	case clue == TYPE_CLUE_INT:
		compare.IntValue, err = strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse int: %v in [%v]", value, line)
		}
	case clue == TYPE_CLUE_QUAD:
		if value == "0" {
			compare.QuadValue = []int64{0}
			return compare, nil
		}
		if !strings.HasPrefix(value, "0x") {
			return nil, fmt.Errorf("0 or hex format is mandatory for quad:  %s", value)
		}
		l := (len(value) - 2) / 16
		v := make([]int64, l)
		for i := range l {
			v[i], err = strconv.ParseInt(value[i*16:(i+1)*16], 0, 64)
			if err != nil {
				return nil, err
			}
		}
		compare.QuadValue = v
	}
	return compare, nil
}

func ParseType(line string) (*Type, error) {
	t := &Type{}
	for _, o := range []byte("/%&") {
		i := strings.IndexByte(line, o)
		if i != -1 {
			t.Name = line[:i]
			t.Operator = o
			t.Arg = line[i:]
			break
		}
	}
	if t.Name == "" {
		t.Name = line
	}
	t.Clue_ = Types[t.Name].Clue_ // FIXME

	return t, nil
}
