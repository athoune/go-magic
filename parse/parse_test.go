package parse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestParseOffset(t *testing.T) {
	o := &model.Offset{}
	err := ParseOffset(o, ">>>4")
	assert.NoError(t, err)
	assert.Equal(t, 3, o.Level)
	assert.Equal(t, int64(4), o.Value)
	assert.False(t, o.Dynamic)

	o = &model.Offset{}
	err = ParseOffset(o, ">>(4.s*512)")
	assert.NoError(t, err)
	assert.Equal(t, 2, o.Level)
	assert.True(t, o.Dynamic)
	assert.Equal(t, int64(4), o.DynOffset)
	assert.Equal(t, uint8('s'), o.DynType)
	assert.Equal(t, uint8('*'), o.DynOperator)
	assert.Equal(t, int64(512), o.DynArg)

	o = &model.Offset{}
	err = ParseOffset(o, ">>>>>>(&4.l+(-4))")
	assert.NoError(t, err)
	assert.Equal(t, 6, o.Level)
	assert.True(t, o.Dynamic)
	assert.Equal(t, int64(4), o.DynOffset)
	assert.Equal(t, uint8('l'), o.DynType)
	assert.Equal(t, uint8('+'), o.DynOperator)
	assert.Equal(t, int64(-4), o.DynArg)
}

func TestParseDynamicOffset(t *testing.T) {
	o := &model.Offset{}
	err := ParseDynamicOffset(o, "4.s*512")
	assert.NoError(t, err)
	assert.True(t, o.Dynamic)
	assert.Equal(t, int64(4), o.DynOffset)
	assert.Equal(t, uint8('*'), o.DynOperator)
	assert.Equal(t, int64(512), o.DynArg)
}

func TestParseCompare(t *testing.T) {
	for _, fixture := range []struct {
		line    string
		type_   *model.Type
		size    int
		compare *model.Compare
	}{
		{"<10", &model.Type{
			Clue_:  model.TYPE_CLUE_INT,
			Root:   "long",
			Signed: true,
		}, 3, &model.Compare{
			Comparator:  COMPARE_LESS,
			RawExpected: "10",
		}},
		{"< 10", &model.Type{
			Clue_:  model.TYPE_CLUE_INT,
			Root:   "long",
			Signed: true,
		}, 4, &model.Compare{
			Comparator:  COMPARE_LESS,
			RawExpected: "10",
		}},
		{"0x01000007", &model.Type{
			Clue_:  model.TYPE_CLUE_INT,
			Root:   "long",
			Signed: true,
		}, 10, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: "0x01000007",
		}},
		{"!>10", &model.Type{
			Clue_:  model.TYPE_CLUE_INT,
			Root:   "long",
			Signed: true,
		}, 4, &model.Compare{
			Comparator:  COMPARE_GREATER,
			RawExpected: "10",
			Not:         true,
		}},
		{"D6E229D3-35DA-11D1-9034-00A0C90349BE", &model.Type{
			Clue_: model.TYPE_CLUE_STRING}, 36, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: "D6E229D3-35DA-11D1-9034-00A0C90349BE",
		}},
		{`Invalid\ partition\ table		english`, &model.Type{
			Clue_: model.TYPE_CLUE_STRING}, 25, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: "Invalid partition table",
		}},
		{`\x6d\x6a\x70\x32`, &model.Type{
			Clue_: model.TYPE_CLUE_STRING,
			Root:  "string",
		}, 16, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: "mjp2",
		}},
	} {
		c, s, err := ParseCompare(fixture.line, fixture.type_)
		assert.NoError(t, err)
		// [FIXME]
		fixture.compare.Type = fixture.type_ // yes, it's cheating
		assert.Equal(t, fixture.compare, c, fixture.line)
		assert.Equal(t, fixture.size, s, fixture.line)
	}
}

func TestParseCompareName(t *testing.T) {
	c, _, err := ParseCompare(`jpeg`, &model.Type{
		Name: "name",
	})
	assert.NoError(t, err)
	assert.Equal(t, "jpeg", c.RawExpected)
}

func TestParseType(t *testing.T) {
	type_, err := ParseType(`belong&0xfe00f0f0`)
	assert.NoError(t, err)
	assert.Equal(t, "belong", type_.Name)
	assert.Equal(t, byte('&'), type_.FilterOperator)
	assert.Equal(t, uint64(0xfe00f0f0), type_.FilterBinaryArgument)

	type_, err = ParseType(`pstring/HJ`)
	assert.NoError(t, err)
	assert.Equal(t, "pstring", type_.Root)
	assert.Equal(t, "pstring", type_.Name)
	assert.Equal(t, byte('/'), type_.FilterOperator)
	assert.Equal(t, "HJ", type_.FilterStringArgument)
}

func TestParse(t *testing.T) {
	fixture := `# Standard PNG image.
0	string		\x89PNG\x0d\x0a\x1a\x0a\x00\x00\x00\x0DIHDR	PNG image data
!:mime	image/png
!:ext   png
!:strength +10
>16	use		png-ihdr
>33	string		\x00\x00\x00\x08acTL	\b, animated
>>41	ubelong		1			(%d frame
>>41	ubelong		>1			(%d frames
>>45	ubelong		0			\b, infinite repetitions)
>>45	ubelong		1			\b, %d repetition)
>>45	ubelong		>1			\b, %d repetitions)
`

	file := model.NewFile()
	file.Name = "images"
	_, err := Parse(strings.NewReader(fixture), file)
	assert.NoError(t, err)
	assert.Len(t, file.Tests, 1)
	test := file.Tests[0]
	assert.Len(t, test.Actions, 3)
	assert.Len(t, test.SubTests, 2)
	assert.Equal(t, "images", test.File)
	assert.Equal(t, 1, test.Line)
	assert.Equal(t, `>33	string		\x00\x00\x00\x08acTL	\b, animated`, test.SubTests[1].Raw)
	assert.Len(t, test.SubTests[1].SubTests, 5)
}
func TestRead(t *testing.T) {
	r := strings.NewReader(`
# RISC OS Chunk File Format
# From RISC OS Programmer's Reference Manual, Appendix D
# We guess the file type from the type of the first chunk.
0	lelong		0xc3cbc6c5	RISC OS Chunk data
>12	string		OBJ_		\b, AOF object
>12	string		LIB_		\b, ALF library
0	string		Draw		RISC OS Draw file data
`)
	file := model.NewFile()
	file.Name = "risc"
	_, err := Parse(r, file)
	assert.NoError(t, err)
	for _, test := range file.Tests {
		fmt.Println(test)
	}
}

func TestParseStringOptions(t *testing.T) {
	typ := &model.Type{}
	parseOptions(typ, "bob/cC")
	assert.Equal(t, model.STRING_OPTIONS_NONE, typ.StringOptions)
	stringOptions, err := parseStringOptions(typ.FilterStringArgument)
	assert.NoError(t, err)
	assert.Equal(t, model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER,
		stringOptions&model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER)
	assert.Equal(t, model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER,
		stringOptions&model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER)
}

func TestParseSearchOptions(t *testing.T) {
	for _, fixture := range []struct {
		line          string
		searchRange   int
		stringOptions model.StringOptions
	}{
		{"search/1/t", 1, model.STRING_OPTIONS_TEXT_FILE},
		{"search/42/cC", 42,
			model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER |
				model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER},
		{"search/727", 727, model.STRING_OPTIONS_NONE},
		{"search/210965/s", 210965, model.REGEX_OPTIONS_OFFSET_START},
	} {
		typ := &model.Type{}
		parseOptions(typ, fixture.line)
		searchRange, stringOptions, err := parseSearchOptions(typ.FilterStringArgument)
		assert.NoError(t, err)
		assert.Equal(t, fixture.searchRange, searchRange)
		assert.Equal(t, fixture.stringOptions, stringOptions)
	}
}
