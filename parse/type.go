package parse

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/athoune/go-magic/model"
)

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
	var ok bool
	t.ByteOrder, t.Root = model.ByteOrderAndSigned(t.Name)
	if t.TypeFamily, ok = model.Types[t.Root]; !ok {
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

// parseSearchOptions returns 'range' and 'options' (aka 'string options')
func parseSearchOptions(args string) (int, model.StringOptions, error) {
	if args == "" {
		return 0, 0, nil
	}
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
