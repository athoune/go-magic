package ast

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/athoune/go-magic/model"
)

var SIZES map[string]int

func init() {
	SIZES = map[string]int{
		"byte":  1,
		"short": 2,
		"long":  4,
		"quad":  8,
	}
}

type TestResult struct {
	test     *model.Test
	target   io.ReadSeeker
	named    map[string]*model.Test
	output   io.Writer
	Strength int
	Mime     string
	Ext      string
	Apple    string // Apple type for Mac OS 9 and older
}

func NewTestResult(test *model.Test, named map[string]*model.Test, output io.Writer) *TestResult {
	return &TestResult{
		test:   test,
		named:  named,
		output: output,
	}
}

func (t *TestResult) Test(target io.ReadSeeker) (bool, error) {
	_, err := target.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}
	t.target = target
	t.offset()
	c, _, err := t.compare()
	if err != nil {
		return false, err
	}
	if t.test.Type.Name == "use" {
		namedTest, ok := t.named[t.test.Compare.StringValue]
		if !ok {
			return false, fmt.Errorf("unknown name: '%s'", t.test.Compare.StringValue)
		}
		_, err = NewTestResult(namedTest, t.named, t.output).Test(target)
		if err != nil {
			return false, err
		}
	}
	if c {
		for _, action := range t.test.Actions {
			err = t.action(action)
			if err != nil {
				return false, err
			}
		}
		for _, sub := range t.test.SubTests {
			_, err = NewTestResult(sub, t.named, t.output).Test(target)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, nil
}

func (t *TestResult) action(a *model.Action) error {
	switch a.Name {
	case "mime":
		t.Mime = a.Arg
	case "apple":
		t.Apple = a.Arg
	case "ext":
		t.Ext = a.Arg
	case "strength":
		/*
			The operand OP can be: +, -, *, or / and VALUE is a constant between 0 and 255.
			This constant is applied using the specified operand to the currently computed
			default magic strength.
		*/
		operator := a.Arg[0]
		value, err := strconv.Atoi(a.Arg[1:])
		if err != nil {
			return err
		}
		switch operator {
		case '+':
			t.Strength = t.Strength + value
		case '-':
			t.Strength = t.Strength - value
		case '*':
			t.Strength = t.Strength * value
		case '/':
			t.Strength = t.Strength / value
		default:
			return fmt.Errorf("strength action: unknown operator '%v' in '%s'", operator, a.Arg)
		}
	default:
		return fmt.Errorf("action: unknown action '%v'", a.Name)
	}
	return nil
}

// offset seeks the target
func (t *TestResult) offset() {
	if t.test.Offset.Relative {
		// TODO
	} else {
		t.target.Seek(t.test.Offset.Value, io.SeekCurrent)
	}
}

// compare return compare result, is special and error
func (t *TestResult) compare() (bool, bool, error) {
	signed, byteOrder, typ := model.EndianessSigned(t.test.Type.Name)

	var bo binary.ByteOrder
	switch byteOrder {
	case model.BIG_ENDIAN:
		bo = binary.BigEndian
	case model.LITTLE_ENDIAN:
		bo = binary.LittleEndian
	default:
		bo = binary.NativeEndian
	}

	size, ok := SIZES[typ]

	if !ok {
		switch typ {
		case "string":
			size = len(t.test.Compare.StringValue)
		case "search":
			size = len(t.test.Compare.StringValue)
		}
	}
	buff := make([]byte, size)
	_, err := t.target.Read(buff)
	if err != nil {
		return false, false, err
	}

	switch typ { // All the special cases
	case "name":
		// Like a LABEL
		return true, true, nil
	case "use":
		// It's just a GOTO
		return false, true, nil
	case "search":
		return false, true, nil
	}

	cmp, ok := TYPE_COMPARATORS[typ]
	if !ok {
		return false, false, fmt.Errorf("unknown type: %v in `%s`", typ, t.test.Raw)
	}
	return cmp.compare(buff, bo, signed, t.test.Type, t.test.Compare)
}

func operation[V int64 | uint64](current, expected V, operation byte) bool {
	switch operation {
	case '=':
		return current == expected
	}
	return false
}

func (t *TestResult) message() error {
	if t.test.Message.IsDisplayable {
		if !t.test.Message.IsTemplate {
			_, err := t.output.Write([]byte(t.test.Message.Value))
			return err
		}
		switch t.test.Type.Clue_ {
		case model.TYPE_CLUE_INT:
			_, err := fmt.Fprintf(t.output, t.test.Message.Value, t.test.Compare.IntValue)
			return err
		case model.TYPE_CLUE_STRING:
			_, err := fmt.Fprintf(t.output, t.test.Message.Value, t.test.Compare.StringValue)
			return err
		default:
			return fmt.Errorf("unknown type for info: '%v'", t.test.Type.Clue_)
		}
	}
	return false
}
