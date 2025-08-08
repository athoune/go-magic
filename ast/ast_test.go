package ast

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReTest(t *testing.T) {

	for _, fixture := range []struct {
		line    string
		type_   string
		compare string
	}{
		{`>>16	belong&0xfe00f0f0	0x3030`, "belong", "&0xfe00f0f0"},
		{`0	lelong		0xc3cbc6c5	RISC OS Chunk data`, "lelong", "0xc3cbc6c5"},
		{`>>>>>>0	ubyte			< 10	Infocom (Z-machine %d`, "ubyte", "< 10"},
		{`0	string 		Draw		RISC OS Draw file data`, "string", "Draw"},
		{`>12	leshort	!1	%d patterns`, "leshort", "!1"},
		{`0	string	\x02\x01\x13\x13\x10\x14\x12\x0e`, "string", `\x02\x01\x13\x13\x10\x14\x12\x0e`},
		{`>9	belong  !0x0A0D1A00	game data, CORRUPTED`, "belong", "!0x0A0D1A00"},
		{`>>>>&1	string		x		"%s"`, "string", "x"},
		{`0	search/8192	(input,`, "search", "(input,"},
		{`>>>>>>&8	ubelong%44100	0`, "ubelong", "0"},
		{`>8		ubyte/4		=0		CHN: 4`, "ubyte", "=0"},
		{`>>&(0x04)	lelong	>0	\b, with %d reference sequences`, "lelong", ">0"},
		{`>>>>(0x3C.b+0x0FF)	string	Invalid\ partition\ table		english`, "string", `Invalid\ partition\ table`},
		{`0	string		AES`, "string", "AES"},
	} {
		m := test_re.FindStringSubmatch(fixture.line)
		assert.True(t, len(m) > 0, "empty regexp's match")
		fmt.Println(fixture.line)
		for i, mm := range m {
			fmt.Println("#", i, mm)
		}
		assert.Equal(t, fixture.compare, strings.Trim(m[compare_test_idx], "\t "))
		assert.Equal(t, fixture.type_, m[type_test_idx])
	}

}

func TestParse(t *testing.T) {

	tests, line, err := Parse(strings.NewReader(`>>16	belong&0xfe00f0f0	0x3030`))
	assert.Nil(t, err)
	assert.Equal(t, 1, line)
	assert.NotNil(t, tests)
	assert.Len(t, tests, 1)
	test := tests[0]
	assert.Equal(t, 2, test.Offset.Level)
	assert.Equal(t, int64(16), test.Offset.Value)
	assert.Equal(t, TYPE_CLUE_INT, test.Type.Clue_)
	assert.Equal(t, "belong", test.Type.Name)
	assert.Equal(t, COMPARE_AND, test.Compare.Operation)
	assert.False(t, test.Compare.Not)

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
	tests, _, _ := Parse(r)
	for _, test := range tests {
		fmt.Println(test)
	}
}
