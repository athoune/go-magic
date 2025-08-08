package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpaces(t *testing.T) {
	i := spaces("  	plop")
	assert.Equal(t, 3, i)
	i = spaces("	 ")
	assert.Equal(t, 2, i)
	i = spaces("")
	assert.Equal(t, 0, i)
}

func TestNotSpaces(t *testing.T) {
	i := notSpace("beuha ")
	assert.Equal(t, 5, i)
	i = notSpace("beuha")
	assert.Equal(t, 5, i)
}

func TestParseLine(t *testing.T) {
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

		l := NewTest()
		err := ParseLine(l, fixture.line)
		assert.NoError(t, err, fixture.line)
		assert.Equal(t, fixture.type_, l.Type.Name)
	}

}
