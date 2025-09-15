package parse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

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
	spaces_re = regexp.MustCompile(`\s+`)

	dynamic_value_re = regexp.MustCompile(`(?<value>(0x)?[0-9a-f]+)((?<separator>[.,])(?<type>[bBcCeEfFgGhHiIlLmsSqQ]))?((?<operator>[+\-*])(?<arg>.*))?`)
	value_dynamic_idx = dynamic_value_re.SubexpIndex("value")
	type_dynamic_idx = dynamic_value_re.SubexpIndex("type")
	operator_dynamic_idx = dynamic_value_re.SubexpIndex("operator")
	arg_dynamic_idx = dynamic_value_re.SubexpIndex("arg")
}

func ParseFolder(path string) (model.Files, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make(model.Files, 0)
	for _, e := range entries {
		f, err := os.Open(path + "/" + e.Name())
		if err != nil {
			return nil, err
		}
		file := model.NewFile()
		file.Name = e.Name()
		_, err = Parse(f, file)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

/*
Parse reads r and compile tests described. Parsed tests are stored in 'file' argument.
*/
func Parse(r io.Reader, file *model.File) (int, error) {
	scanner := bufio.NewScanner(r)
	var slugs []string
	var err error
	testsParsing := NewTestsParsing()
	n_line := 0
	for scanner.Scan() {
		n_line += 1
		line := scanner.Text()
		if err = scanner.Err(); err != nil {
			return n_line, err
		}
		if len(line) == 0 { // empty
			continue
		}
		if line[0] == '#' { // comment
			continue
		}
		if strings.HasPrefix(line, "!:") {
			slugs = spaces_re.Split(line[2:], -1)
			testsParsing.AppendAction(&model.Action{
				Name: slugs[0],
				Arg:  slugs[1],
			})
			continue
		}
		test := model.NewTest()
		test.Line = n_line - 1
		test.File = file.Name
		err = ParseLine(test, line)
		if err != nil {
			return n_line, err
		}
		if test.Type.Name == "name" {
			file.Names[test.Compare.RawExpected] = test
		}
		testsParsing.AppendTest(test)
	}
	file.Tests = testsParsing.Tests
	return n_line, nil
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
	if type_.Clue_ == model.TYPE_CLUE_STRING {
		compare.RawExpected, err = HandleStringEscape(value)
		if err != nil {
			return nil, poz + end, err
		}
	} else {
		compare.RawExpected = value
	}
	//
	return compare, poz + end, nil
}

func parseOptions(typ *model.Type, line string) {
	for _, o := range []byte("/%&+-^*|") {
		i := strings.IndexByte(line, o)
		if i != -1 {
			typ.Name = line[:i]
			typ.FilterOperator = o
			typ.FilterStringArgument = line[i+1:]
			break
		}
	}
}

func ParseType(line string) (*model.Type, error) {
	t := &model.Type{}
	var err error

	parseOptions(t, line)
	if t.Name == "" {
		t.Name = line
	}
	t.Signed, t.ByteOrder, t.Root = model.ByteOrderAndSigned(t.Name)
	var ok bool
	if t.Clue_, ok = model.Types[t.Name]; !ok {
		return nil, fmt.Errorf("unknown type [%v]", t.Name)
	}
	switch t.Root {
	case "indirect":
	/*
		indirect        Starting at the given offset, consult the magic database again.  The offset of the indirect magic is by default absolute in the file, but one can specify /r to indicate that the offset is relative from the
		                beginning of the entry.
	*/
	case "pstring":
		/*
			boring, the options are handled in the ast
		*/
	case "string":
		/*
			The string type specification can be optionally followed by /[WwcCtbTf]*.
		*/
		t.StringOptions, err = parseStringOptions(t.FilterStringArgument)
		if err != nil {
			return nil, err
		}
	case "regex":
	case "search":
		t.SearchRange, t.StringOptions, err = parseSearchOptions(t.FilterStringArgument)
		if err != nil {
			return nil, err
		}
	default: // all integers
		if t.FilterOperator != 0 {
			t.FilterBinaryArgument, err = strconv.ParseUint(t.FilterStringArgument, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("%s with line : %s", err.Error(), line)
			}
		}
	}

	return t, nil
}

func parseStringOptions(stringOptionsRaw string) (model.StringOptions, error) {
	stringOptions := model.STRING_OPTIONS_NONE
	for _, option := range stringOptionsRaw {
		// WwcCtbTf
		switch option {
		default:
			return 0, fmt.Errorf("unknown String options: %v", option)
		case 'W':
			stringOptions |= model.STRING_OPTIONS_COMPACT_WITH_SPACES
		case 'f':
			stringOptions |= model.STRING_OPTIONS_FULL_WORD
		case 'c':
			stringOptions |= model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER
		case 'C':
			stringOptions |= model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER
		case 't':
			stringOptions |= model.STRING_OPTIONS_TEXT_FILE
		case 'b':
			stringOptions |= model.STRING_OPTIONS_BINARY_FILE
		case 'T':
			stringOptions |= model.STRING_OPTIONS_TRIMMED
		case '4':
			/*
				0	string/4	MOC3	Live2D Cubism MOC3
				[FIXME] What does it means ? Just read 4 characters ?
				No info in the man, nor the NBF file
			*/
		case 's': // yes, it happens in the MagDir
			stringOptions |= model.REGEX_OPTIONS_OFFSET_START
		}
	}
	return stringOptions, nil
}

func parseSearchOptions(args string) (int, model.StringOptions, error) {
	s := strings.IndexRune(args, '/')
	if s == 0 {
		return 0, model.STRING_OPTIONS_NONE,
			fmt.Errorf("search options parsing error, can't start with a '/', 'range' is mandatory : %s", args)
	}
	if s == -1 { // no string options
		r, err := strconv.ParseInt(args, 0, 64)
		if err != nil {
			return 0, model.STRING_OPTIONS_NONE, err
		}
		return int(r), model.STRING_OPTIONS_NONE, err
	} // string options
	nums := new(bytes.Buffer)
	opts := new(bytes.Buffer)
	for _, l := range args[s+1:] {
		if unicode.IsNumber(l) {
			nums.WriteRune(l)
			// numbers should be consecutive, but YOLO
		} else {
			opts.WriteRune(l)
		}
	}
	if nums.Len() > 0 {
		strconv.ParseInt(nums.String(), 0, 64)
	}
	so, err := parseStringOptions(opts.String())
	if err != nil {
		return 0, model.STRING_OPTIONS_NONE, err
	}
	r, err := strconv.ParseInt(args[:s], 0, 64)
	if err != nil {
		return 0, model.STRING_OPTIONS_NONE, err
	}
	return int(r), so, nil
}
