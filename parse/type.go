package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/athoune/go-magic/model"
)

var RE_STRING_OPTIONS *regexp.Regexp
var RE_INT *regexp.Regexp

func init() {
	RE_STRING_OPTIONS = regexp.MustCompile(`(\d+)(/?)`)
	RE_INT = regexp.MustCompile(`(0x[0-9a-fA-F]+)|(\d+)`)
}

type UnknownStringOption struct {
	Options string
}

func (u *UnknownStringOption) Error() string {
	return fmt.Sprintf("unknown String options: %v", u.Options)
}

/*
"search/b/100" =>

	Name="search"
	FilterOperator='/'
	FilterStringArgument="b/100"
*/
func splitTypeAndModifiers(line string) (string, byte, string) {
	for _, o := range []byte("/%&+-^*|") {
		i := strings.IndexByte(line, o)
		if i != -1 {
			return line[:i], o, line[i+1:]
		}
	}
	return line, 0, ""
}

func ParseType(line string) (*model.Type, error) {
	t := &model.Type{}
	var err error
	t.Name, t.FilterOperator, t.FilterStringArgument = splitTypeAndModifiers(line)

	var ok bool
	var signed bool
	t.ByteOrder, signed, t.Root = model.ByteOrderAndSigned(t.Name)
	if t.TypeFamily, ok = model.Types[t.Root]; !ok {
		return nil, fmt.Errorf("unknown root name type: %v", t.Name)
	}
	if t.TypeFamily == model.TYPE_FAMILY_INT && !signed {
		t.TypeFamily = model.TYPE_FAMILY_UINT
	}
	switch t.Root {
	case "indirect":
	/*
		indirect        Starting at the given offset, consult the magic database again.  The offset of the indirect magic is by default absolute in the file, but one can specify /r to indicate that the offset is relative from the
		                beginning of the entry.
	*/
	case "default":
	case "use":
	case "name":
	case "clear":
	case "offset":
	case "der":
		/*
		 */
	case "pstring":
		t.StringOptions, t.StringIntOption, err = parseOptions(t.FilterStringArgument)
		if err != nil {
			u, ok := err.(*UnknownStringOption)
			if !ok {
				return nil, err
			}
			err = pstringModifier(t, u.Options)
			if err != nil {
				return nil, err
			}
		}
	case "string", "string16", "search":
		/*
		 */
		t.StringOptions, t.StringIntOption, err = parseOptions(t.FilterStringArgument)
		if err != nil {
			return nil, err
		}
	case "regex":
		var raw string
		t.StringIntOption, raw, err = splitSearchStringOptions(t.FilterStringArgument)
		if err != nil {
			return nil, err
		}
		options, err := readRegexOptions(raw)
		if err != nil {
			u, ok := err.(*UnknownStringOption)
			if ok {
				opt, _, err := parseOptions(u.Options)
				if err != nil {
					return nil, err
				}
				t.StringOptions |= opt
			} else {
				return nil, err
			}
		}
		t.StringOptions |= options
	case "msdostime":
	case "msdosdate", "date", "qdate", "qldate", "ldate", "qwdate":
	case "guid":
	case "byte", "short", "long", "quad", "float", "double", "4", "8":
		if t.FilterOperator != 0 {
			t.FilterBinaryArgument, err = strconv.ParseUint(t.FilterStringArgument, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("%s with line : %s", err.Error(), line)
			}
		}
	default:
		return nil, fmt.Errorf("unknown parse type name: %s", t.Root)
	}

	return t, nil
}

func pstringModifier(typ_ *model.Type, modifiers string) error {
	for _, modifier := range modifiers {
		switch modifier {
		default:
			typ_.StringIntOption = 1
		case 'B':
			typ_.StringIntOption = 1
		case 'H':
			typ_.StringIntOption = 2
			typ_.ByteOrder = model.BIG_ENDIAN
		case 'h':
			typ_.StringIntOption = 2
			typ_.ByteOrder = model.LITTLE_ENDIAN
		case 'L':
			typ_.StringIntOption = 4
			typ_.ByteOrder = model.BIG_ENDIAN
		case 'l':
			typ_.StringIntOption = 4
			typ_.ByteOrder = model.LITTLE_ENDIAN
		case 'J':
			typ_.StringOptions = model.PSTRING_OPTIONS_SIZE_INCLUDE
		}
	}
	return nil
}

// splitSearchStringOptions => intValue, opts
func splitSearchStringOptions(stringOptionsRaw string) (int, string, error) {
	if stringOptionsRaw == "" {
		return 0, "", nil
	}
	slugs := strings.Split(stringOptionsRaw, "/")
	switch len(slugs) {
	case 2:
		var rank, opt int
		if unicode.IsNumber(rune(slugs[1][0])) {
			rank = 1
			opt = 0
		} else {
			rank = 0
			opt = 1
		}
		intValue, err := strconv.ParseInt(slugs[rank], 0, 64)
		if err != nil {
			return -1, "", err
		}
		return int(intValue), slugs[opt], nil
	case 1:
		loc := RE_INT.FindStringIndex(slugs[0])
		if loc == nil {
			return 0, slugs[0], nil
		}
		intValue, err := strconv.ParseInt(slugs[0][loc[0]:loc[1]], 0, 64)
		if err != nil {
			return -1, "", err
		}
		if len(slugs[0]) == loc[1]-1 {
			return int(intValue), "", nil
		}
		if loc[0] == 0 {
			return int(intValue), slugs[0][loc[1]:], nil
		}
		return int(intValue), slugs[0][:loc[0]], nil
	default:
		return -1, "", fmt.Errorf("this StringOptions has too many elements: %v", slugs)
	}
}

func readRegexOptions(stringOptionsRaw string) (model.StringOptions, error) {
	stringOptions := model.STRING_OPTIONS_NONE
	var unknown strings.Builder
	for _, option := range stringOptionsRaw {
		switch option {
		case 'c':
			stringOptions |= model.REGEX_OPTIONS_CASE_INSENSITIVE
		case 's':
			stringOptions |= model.REGEX_OPTIONS_OFFSET_START
		case 'l':
			stringOptions |= model.REGEX_OPTIONS_LINES
		default:
			unknown.WriteRune(option)
		}
	}
	if unknown.Len() == 0 {
		return stringOptions, nil
	}
	return stringOptions, &UnknownStringOption{unknown.String()}
}

func readStringOptions(stringOptionsRaw string) (model.StringOptions, error) {
	stringOptions := model.STRING_OPTIONS_NONE
	var unknown strings.Builder

	for _, option := range stringOptionsRaw {
		// WwcCtbTf
		switch option {
		default:
			unknown.WriteRune(option)
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
		case 's': // yes, it happens in the MagDir
			stringOptions |= model.REGEX_OPTIONS_OFFSET_START
		case 'w':
			/* [FIXME]
			>>0	string/w	#!\ 	a
			*/
		}
	}
	if unknown.Len() == 0 {
		return stringOptions, nil
	}
	return stringOptions, &UnknownStringOption{unknown.String()}
}

func parseOptions(stringOptionsRaw string) (model.StringOptions, int, error) {
	intValue, rawOptions, err := splitSearchStringOptions(stringOptionsRaw)
	if err != nil {
		return model.STRING_OPTIONS_NONE, 0, err
	}
	stringOptions, err := readStringOptions(rawOptions)
	if err != nil {
		return model.STRING_OPTIONS_NONE, 0, err
	}
	return stringOptions, intValue, nil
}
