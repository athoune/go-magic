package ast

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
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
		if t.test.Offset.Dynamic {
			// Dynamic relative offset: read a value at DynOffset, apply operation, use as absolute offset
			computed, err := t.computeDynamicOffset()
			if err == nil {
				_, _ = t.target.Seek(computed, io.SeekStart)
			}
		} else {
			// Simple relative offset: add to current position
			_, _ = t.target.Seek(t.test.Offset.Value, io.SeekCurrent)
		}
	} else {
		_, _ = t.target.Seek(t.test.Offset.Value, io.SeekStart)
	}
}

// computeDynamicOffset reads a value at Offset.Value + DynOffset, applies DynOperator with DynArg,
// and returns the result as a relative offset.
// C reference: softmagic.c mget() INDIR + do_ops()
func (t *TestResult) computeDynamicOffset() (int64, error) {
	// Save current position
	currentPos, err := t.target.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Seek to the position where we read the indirect value
	// In C: s + offset + off where offset is the magic entry's offset
	readPos := t.test.Offset.Value + t.test.Offset.DynOffset
	_, err = t.target.Seek(readPos, io.SeekStart)
	if err != nil {
		// Restore position on error
		_, _ = t.target.Seek(currentPos, io.SeekStart)
		return 0, err
	}

	// Read the value according to DynType
	var lhs int64
	switch t.test.Offset.DynType {
	case 'b': // byte
		var b byte
		err = binary.Read(t.target, binary.NativeEndian, &b)
		lhs = int64(b)
	case 's': // short, native endian
		var s int16
		err = binary.Read(t.target, binary.NativeEndian, &s)
		lhs = int64(s)
	case 'S': // short, big endian
		var s int16
		err = binary.Read(t.target, binary.BigEndian, &s)
		lhs = int64(s)
	case 'l': // long, native endian
		var l int32
		err = binary.Read(t.target, binary.NativeEndian, &l)
		lhs = int64(l)
	case 'L': // long, big endian
		var l int32
		err = binary.Read(t.target, binary.BigEndian, &l)
		lhs = int64(l)
	case 'q': // quad, native endian
		var q int64
		err = binary.Read(t.target, binary.NativeEndian, &q)
		lhs = q
	case 'Q': // quad, big endian
		var q int64
		err = binary.Read(t.target, binary.BigEndian, &q)
		lhs = q
	default:
		// Restore position on error
		_, _ = t.target.Seek(currentPos, io.SeekStart)
		return 0, fmt.Errorf("unsupported dynamic offset type: %c", t.test.Offset.DynType)
	}

	if err != nil {
		// Restore position on error
		_, _ = t.target.Seek(currentPos, io.SeekStart)
		return 0, err
	}

	// Apply the operation (e.g., + - * / & | %)
	result, err := operation.Operate(lhs, t.test.Offset.DynArg, t.test.Offset.DynOperator)
	if err != nil {
		// Restore position on error
		_, _ = t.target.Seek(currentPos, io.SeekStart)
		return 0, err
	}

	// Restore position before seeking
	_, _ = t.target.Seek(currentPos, io.SeekStart)
	return result, nil
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
	case "search":
		return t.compareSearch()
	case "regex":
		return t.compareRegex()
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
		if t.test.Compare.Comparator == '&' || t.test.Compare.Comparator == '^' {
			ok, err = operation.CompareBits(uint64(value.IntValue),
				uint64(t.test.Compare.Expected.IntValue),
				t.test.Compare.Comparator)
		} else {
			ok, err = operation.CompareNumbers(value.IntValue,
				t.test.Compare.Expected.IntValue,
				t.test.Compare.Comparator)
		}
	case model.TYPE_FAMILY_UINT:
		if t.test.Compare.Comparator == '&' || t.test.Compare.Comparator == '^' {
			ok, err = operation.CompareBits(value.UIntValue,
				t.test.Compare.Expected.UIntValue,
				t.test.Compare.Comparator)
		} else {
			ok, err = operation.CompareNumbers(value.UIntValue,
				t.test.Compare.Expected.UIntValue,
				t.test.Compare.Comparator)
		}
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

// cfmt converts C printf format specifiers to Go-compatible ones.
func cfmt(format string) string {
	// Go's fmt does not support %u; %d handles all integer types.
	return strings.ReplaceAll(format, "%u", "%d")
}

func (t *TestResult) message() error {
	if !t.test.Message.IsTemplate {
		_, err := t.output.Write([]byte(t.test.Message.Value))
		return err
	}
	// Use the actual value read from the file for formatting; fall back to
	// the type's filter argument only when no value was read (e.g., x test).
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
		_, err := fmt.Fprintf(t.output, cfmt(t.test.Message.Value), arg)
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
		_, err := fmt.Fprintf(t.output, cfmt(t.test.Message.Value), arg)
		return err
	default:
		return fmt.Errorf("unknown type for info: '%v'", t.test.Type.TypeFamily)
	}
}

// compareSearch implements FILE_SEARCH: search for a pattern in a range of bytes.
// C reference: softmagic.c magiccheck() case FILE_SEARCH
func (t *TestResult) compareSearch() (bool, bool, error) {
	// Read the file content from current offset
	// In C, mcopy for FILE_SEARCH sets ms->search.s = s + offset
	// with ms->search.s_len = nbytes - offset
	pos, err := t.target.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, false, err
	}

	// Determine how many bytes to read
	range_ := int(t.test.Type.SearchCount) // str_range
	if range_ == 0 {
		range_ = 4096 // default search range
	}

	buf := make([]byte, range_)
	n, err := io.ReadAtLeast(t.target, buf, 1)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false, false, err
	}
	buf = buf[:n]

	// Search for the pattern
	pattern := t.test.Compare.Expected.StringValue
	idx := bytes.Index(buf, []byte(pattern))
	if idx == -1 {
		return false, false, nil
	}

	// Found: seek to the match position
	_, _ = t.target.Seek(pos+int64(idx), io.SeekStart)
	t.readValue = &model.Value{StringValue: string(buf[idx:])}
	return true, false, nil
}

// compareRegex implements FILE_REGEX: match a regular expression against file content.
// C reference: softmagic.c magiccheck() case FILE_REGEX
func (t *TestResult) compareRegex() (bool, bool, error) {
	// Read the file content from current offset
	pos, err := t.target.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, false, err
	}

	// Determine how many bytes to read
	range_ := int(t.test.Type.SearchCount) // str_range
	if range_ == 0 {
		range_ = 8192 // default regex search range
	}

	buf := make([]byte, range_)
	n, err := io.ReadAtLeast(t.target, buf, 1)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false, false, err
	}
	buf = buf[:n]

	// Compile and match the regex
	pattern := t.test.Compare.Expected.StringValue
	flags := ""
	if t.test.Type.StringOptions&model.REGEX_OPTIONS_CASE_INSENSITIVE != 0 {
		flags = "(?i)"
	}
	re, err := regexp.Compile(flags + pattern)
	if err != nil {
		return false, false, fmt.Errorf("regex compile: %w", err)
	}

	loc := re.FindIndex(buf)
	if loc == nil {
		return false, false, nil
	}

	// Found: seek to the match start
	_, _ = t.target.Seek(pos+int64(loc[0]), io.SeekStart)
	t.readValue = &model.Value{StringValue: string(buf[loc[0]:loc[1]])}
	return true, false, nil
}
