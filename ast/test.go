package ast

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/athoune/go-magic/model"
	"github.com/athoune/go-magic/operation"
	"github.com/athoune/go-magic/unpack"
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
	readValue *model.Value // actual value read from file during compare
}

func NewTestResult(test *model.Test, named map[string]*model.Test, output io.Writer) *TestResult {
	return &TestResult{
		test:   test,
		named:  named,
		output: output,
	}
}

func (t *TestResult) Test(target io.ReadSeeker) (bool, error) {
	if t.test.Offset.Level == 0 {
		_, err := target.Seek(0, io.SeekStart)
		if err != nil {
			return false, err
		}
	}
	t.target = target
	t.offset() // Move the cursor
	c, _, err := t.compare()
	if err != nil {
		return false, err
	}
	if t.test.Type.Name == "use" {
		namedTest, ok := t.named[t.test.Compare.RawExpected]
		if !ok {
			return false, fmt.Errorf("unknown name: '%s'", t.test.Compare.RawExpected)
		}
		_, err = NewTestResult(namedTest, t.named, t.output).Test(target)
		if err != nil {
			return false, err
		}
	}
	if c {
		err = t.message()
		if err != nil {
			return false, err
		}
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
		t.Strength, err = operation.Operate(t.Strength, value, operator)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("action: unknown action '%v'", a.Name)
	}
	return nil
}

// offset seeks the target to the absolute offset
func (t *TestResult) offset() {
	if t.test.Offset.Relative {
		// TODO
	} else {
		t.target.Seek(t.test.Offset.Value, io.SeekStart)
	}
}

// compare return compare result, is special and error
func (t *TestResult) compare() (bool, bool, error) {
	switch t.test.Type.Name { // All the special cases
	case "name":
		// Like a LABEL
		return true, true, nil
	case "use":
		// It's just a GOTO
		return false, true, nil
	case "search": // [FIXME]
		return false, true, errors.ErrUnsupported
	case "regex": // [FIXME]
		return false, true, errors.ErrUnsupported
	}

	value, _, err := unpack.ReadValue(t.test.Type, t.target)
	if err != nil {
		return false, false, err
	}
	t.readValue = value
	if t.test.Compare.X {
		/*err = t.message()
		if err != nil {
			return false, false, err
		}*/
		return true, false, nil
	}

	FilterValue(t.test.Type, value)

	var ok bool
	switch t.test.Type.TypeFamily {
	case model.TYPE_FAMILY_INT:
		ok, err = operation.CompareNumbers(value.IntValue,
			t.test.Compare.Expected.IntValue,
			t.test.Compare.Comparator)
	case model.TYPE_FAMILY_UINT:
		ok, err = operation.CompareNumbers(value.UIntValue,
			t.test.Compare.Expected.UIntValue,
			t.test.Compare.Comparator)
	case model.TYPE_FAMILY_FLOAT:
		ok, err = operation.CompareNumbers(value.FloatValue,
			t.test.Compare.Expected.FloatValue,
			t.test.Compare.Comparator)
	case model.TYPE_FAMILY_STRING:
		ok, err = operation.CompareString(value.StringValue,
			t.test.Compare.Expected.StringValue,
			t.test.Compare.Comparator,
			t.test.Type.StringOptions)
	}

	/*err = t.message()
	if err != nil {
		return false, false, err
	}*/
	return ok, false, err
}

func (t *TestResult) message() error {
	if !t.test.Message.IsTemplate {
		_, err := t.output.Write([]byte(t.test.Message.Value))
		return err
	}
	// [FIXME] is it Type value ?!
	switch t.test.Type.TypeFamily {
	case model.TYPE_FAMILY_INT, model.TYPE_FAMILY_UINT, model.TYPE_FAMILY_FLOAT:
		var arg interface{}
		if t.readValue != nil {
			switch t.test.Type.TypeFamily {
			case model.TYPE_FAMILY_INT:
				arg = t.readValue.IntValue
			case model.TYPE_FAMILY_UINT:
				arg = t.readValue.UIntValue
			case model.TYPE_FAMILY_FLOAT:
				arg = t.readValue.FloatValue
			}
		} else {
			arg = t.test.Type.FilterBinaryArgument
		}
		_, err := fmt.Fprintf(t.output, t.test.Message.Value, arg)
		return err
	case model.TYPE_FAMILY_STRING:
		var arg string
		if t.test.Compare.Comparator == '=' || t.test.Compare.Comparator == '!' {
			// C mprint: for '=' and '!' relations, print the pattern (m->value.s)
			arg = t.test.Compare.Expected.StringValue
		} else {
			// C mprint: for other relations, print the file content (p->s)
			arg = t.readValue.StringValue
			if t.test.Type.StringOptions&model.STRING_OPTIONS_TRIMMED != 0 {
				arg = strings.TrimSpace(arg)
			}
		}
		_, err := fmt.Fprintf(t.output, t.test.Message.Value, arg)
		return err
	default:
		return fmt.Errorf("unknown type for info: '%v'", t.test.Type.TypeFamily)
	}
}
