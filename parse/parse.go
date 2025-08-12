package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
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

var OPERATIONS = []byte("=><&^~")

var spaces_re *regexp.Regexp
var dynamic_value_re *regexp.Regexp

var value_dynamic_idx int
var type_dynamic_idx int
var operator_dynamic_idx int
var arg_dynamic_idx int

func init() {
	dynamic_value_re = regexp.MustCompile(`(?<value>(0x)?[0-9a-f]+)((?<separator>[.,])(?<type>[bBcCeEfFgGhHiIlLmsSqQ]))?((?<operator>[+\-*])(?<arg>.*))?`)
	spaces_re = regexp.MustCompile(`\s+`)

	value_dynamic_idx = dynamic_value_re.SubexpIndex("value")
	type_dynamic_idx = dynamic_value_re.SubexpIndex("type")
	operator_dynamic_idx = dynamic_value_re.SubexpIndex("operator")
	arg_dynamic_idx = dynamic_value_re.SubexpIndex("arg")

}

func Parse(r io.Reader) ([]*model.Test, int, error) {
	scanner := bufio.NewScanner(r)
	var slugs []string
	var err error
	var previous *model.Test
	tests := make([]*model.Test, 0)
	n_line := 0
	for scanner.Scan() {
		n_line += 1
		line := scanner.Text()
		if err = scanner.Err(); err != nil {
			return nil, n_line, err
		}
		if len(line) == 0 { // empty
			continue
		}
		if line[0] == '#' { // comment
			continue
		}
		test := model.NewTest()
		if previous != nil && strings.HasPrefix(line, "!:") {
			slugs = spaces_re.Split(line[2:], -1)
			previous.Actions = append(previous.Actions, &model.Action{
				Name: slugs[0],
				Arg:  slugs[1],
			})
			continue
		}
		err = ParseLine(test, line)
		if err != nil {
			return nil, n_line, err
		}
		if test.Offset.Level == 0 {
			previous = test
		}
		if previous != nil && test.Offset.Level >= previous.Offset.Level {
			previous.SubTests = append(previous.SubTests, test)
		} else {
			previous = test
			tests = append(tests, test)
		}
	}
	return tests, n_line, nil
}
func ParseOffset(offset *model.Offset, line string) error {
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
		i := strings.LastIndexByte(line[poz+1:], ')')
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

func ParseDynamicOffset(offset *model.Offset, line string) error {
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
		start := strings.IndexByte(arg, '(')
		if start != -1 {
			// In msdos#175
			// >>>>>>(&4.l+(-4))	string		ITOLITLS	\b, Microsoft compiled help format 2.0
			end := strings.LastIndexByte(arg[1:], ')')
			if end == -1 {
				return fmt.Errorf("parenthesis mismatch: [%v]", arg)
			}
			arg = arg[1 : len(arg)-1]
		}
		offset.DynArg, err = strconv.ParseInt(arg, 0, 32)
		if err != nil {
			return fmt.Errorf("can't parse int in %v", arg)
		}
	}
	return nil
}

// ParseCompare extract the operation, the value (typed) and the new position
func ParseCompare(line string, clue model.Clue) (*model.Compare, int, error) {
	if line[0] == 'x' {
		return nil, 1, nil
	}
	compare := &model.Compare{
		Type: clue,
	}
	var err error
	poz := 0
	if line[poz] == '!' {
		compare.Not = true
		poz++
	}
	compare.Operation = line[poz]
	if !IsOperation(compare.Operation) {
		compare.Operation = '='
	} else {
		poz++
	}
	if line[poz] == ' ' {
		poz++
	}
	end := HandleSpaceEscape(line[poz:])
	value := line[poz : poz+end]
	if clue == model.TYPE_CLUE_STRING {
		value, _ = HandleStringEscape(value)
	}
	switch {
	case clue == model.TYPE_CLUE_STRING:
		compare.StringValue, err = HandleStringEscape(value)
		if err != nil {
			return nil, poz + end, err
		}
	case clue == model.TYPE_CLUE_FLOAT:
		compare.FloatValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, poz + end, fmt.Errorf("can't parse float: %v in [%v]", value, line)
		}
	case clue == model.TYPE_CLUE_INT:
		// In filesystems#1160 there is :
		// 0	lelong		0x1b031336L	Netboot image,
		// In mail.news#91
		// >>15	ulelong		!0x00010000h	\b, version %#8.8
		if len(value) > 1 {
			for _, s := range []byte("hL") {
				if value[len(value)-1] == byte(s) {
					value = value[:len(value)-1]
					break
				}
			}
		}
		compare.IntValue, err = strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, poz + end, fmt.Errorf("can't parse int: %v in [%v]", value, line)
		}
	case clue == model.TYPE_CLUE_QUAD:
		if value == "0" {
			compare.QuadValue = []int64{0}
			return compare, poz + end, nil
		}
		if strings.HasPrefix(value, "0x") {
			l := (len(value) - 2) / 8
			v := make([]int64, l)
			vv := value[2:]
			for i := 0; i < l; i++ {
				v[i], err = strconv.ParseInt(vv[i*8:(i+1)*8], 16, 64)
				if err != nil {
					return nil, 0, err
				}
			}
			compare.QuadValue = v
		} else {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, poz + end, fmt.Errorf("can't parse int: %v in [%v]", value, line)
			}
			compare.QuadValue = []int64{v, 0}
		}
	default:
		return nil, 0, fmt.Errorf("unknown clue: %v", clue)
	}
	return compare, poz + end, nil
}

func ParseType(line string) (*model.Type, error) {
	t := &model.Type{}
	for _, o := range []byte("/%&+-^*|") {
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
	tt, ok := model.Types[t.Name]
	if !ok {
		return nil, fmt.Errorf("unknown type [%v]", t.Name)
	}
	t.Clue_ = tt.Clue_ // FIXME

	return t, nil
}
