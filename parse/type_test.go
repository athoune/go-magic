package parse

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestParseTypeFilter(t *testing.T) {
	for _, fixture := range []struct {
		line     string
		expected *model.Type
	}{
		{`regex/1l`, &model.Type{
			Name:                 "regex",
			Root:                 "regex",
			StringIntOption:      1,
			FilterOperator:       byte('/'),
			FilterStringArgument: "1l",
			StringOptions:        model.REGEX_OPTIONS_LINES,
			TypeFamily:           model.TYPE_FAMILY_STRING,
		}},
		{`belong&0xfe00f0f0`, &model.Type{
			Name:                 "belong",
			Root:                 "long",
			FilterOperator:       byte('&'),
			FilterStringArgument: "0xfe00f0f0",
			FilterBinaryArgument: uint64(0xfe00f0f0),
			TypeFamily:           model.TYPE_FAMILY_INT,
			ByteOrder:            model.BIG_ENDIAN,
		}},
		{`pstring/HJ`, &model.Type{
			Name:                 "pstring",
			Root:                 "pstring",
			FilterOperator:       byte('/'),
			FilterStringArgument: "HJ",
			StringIntOption:      2,
			ByteOrder:            model.BIG_ENDIAN,
			StringOptions:        model.PSTRING_OPTIONS_SIZE_INCLUDE,
			TypeFamily:           model.TYPE_FAMILY_STRING,
		}},
	} {
		typ, err := ParseType(fixture.line)
		assert.NoError(t, err, fixture.line)
		assert.Equal(t, fixture.expected, typ, fixture.line)
	}
}

func TestParseType(t *testing.T) {
	for _, fixture := range []struct {
		line         string
		expectedType *model.Type
	}{
		{`u8`, &model.Type{
			Name:       "u8",
			Root:       "8",
			ByteOrder:  model.NATIVE_ENDIAN,
			TypeFamily: model.TYPE_FAMILY_UINT,
		}},
		{`ulong`, &model.Type{
			Name:       "ulong",
			Root:       "long",
			ByteOrder:  model.NATIVE_ENDIAN,
			TypeFamily: model.TYPE_FAMILY_UINT,
		}},
		{`medate`, &model.Type{
			Name:       "medate",
			Root:       "date",
			ByteOrder:  model.MIDDLE_ENDIAN,
			TypeFamily: model.TYPE_FAMILY_STRING,
		}},
		{`lemsdostime`, &model.Type{
			Name:       "lemsdostime",
			Root:       "msdostime",
			ByteOrder:  model.LITTLE_ENDIAN,
			TypeFamily: model.TYPE_FAMILY_STRING,
		}},
	} {
		typ, err := ParseType(fixture.line)
		assert.NoError(t, err)
		assert.Equal(t, fixture.expectedType, typ, fixture.line)
	}
}

func TestSplitSearchStringOptions(t *testing.T) {
	for _, fixture := range []struct {
		line     string
		intValue int
		opts     model.StringOptions
	}{
		{"T/0x1f", 31, model.STRING_OPTIONS_TRIMMED},
		{"27T", 27, model.STRING_OPTIONS_TRIMMED},
		{"T5", 5, model.STRING_OPTIONS_TRIMMED},
		{"8", 8, model.STRING_OPTIONS_NONE},
		{"", 0, model.STRING_OPTIONS_NONE},
		{"cC", 0, model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER | model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER},
		{"T/42", 42, model.STRING_OPTIONS_TRIMMED},
		{"42/T", 42, model.STRING_OPTIONS_TRIMMED},
	} {
		opts, intValue, err := parseOptions(fixture.line)
		assert.NoError(t, err, fixture.line)
		assert.Equal(t, fixture.intValue, intValue, fixture.line)
		assert.Equal(t, fixture.opts, opts, fixture.line)
	}
}
func TestParseStringOptions(t *testing.T) {
	for _, fixture := range []struct {
		line     string
		options  model.StringOptions
		intValue int
	}{
		{
			"cC",
			model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER | model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER,
			0,
		},
		{
			"8",
			model.STRING_OPTIONS_NONE,
			8,
		},
		{
			"24/Tb",
			model.STRING_OPTIONS_TRIMMED | model.STRING_OPTIONS_BINARY_FILE,
			24,
		},
		{"b/100", model.STRING_OPTIONS_BINARY_FILE, 100},
		{"0x1800/s", model.REGEX_OPTIONS_OFFSET_START, 0x1800},
		{"0x93e4f", model.STRING_OPTIONS_NONE, 0x93e4f},
		{"b5", model.STRING_OPTIONS_BINARY_FILE, 5},
		{"1/t", model.STRING_OPTIONS_TEXT_FILE, 1},
		{"42/cC",
			model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER |
				model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER, 42},
		{"727", model.STRING_OPTIONS_NONE, 727},
		{"210965/s", model.REGEX_OPTIONS_OFFSET_START, 210965},
		{"8192", model.STRING_OPTIONS_NONE, 8192},
	} {
		value, optionsRaw, err := splitSearchStringOptions(fixture.line)
		assert.NoError(t, err)
		assert.Equal(t, fixture.intValue, value)
		options, err := readStringOptions(optionsRaw)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, options&fixture.options)
	}

}
