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
		compare *Compare
	}{
		{`>>16	belong&0xfe00f0f0	0x3030`, "belong", &Compare{
			Operation: '=',
			IntValue:  0x3030,
			Type:      TYPE_CLUE_INT,
		}},
		{`0	lelong		0xc3cbc6c5	RISC OS Chunk data`, "lelong", &Compare{
			Operation: '=',
			IntValue:  0xc3cbc6c5,
			Type:      TYPE_CLUE_INT,
		}},
		{`>>>>>>0	ubyte			< 10	Infocom (Z-machine %d`, "ubyte", &Compare{
			Operation: '<',
			IntValue:  10,
			Type:      TYPE_CLUE_INT,
		}},
		{`0	string 		Draw		RISC OS Draw file data`, "string", &Compare{
			Operation:   '=',
			StringValue: "Draw",
			Type:        TYPE_CLUE_STRING,
		}},
		{`>12	leshort	!1	%d patterns`, "leshort", &Compare{
			Not:       true,
			Operation: '=',
			IntValue:  1,
			Type:      TYPE_CLUE_INT,
		}},
		{`0	string	\x02\x01\x13\x13\x10\x14\x12\x0e`, "string", &Compare{
			Operation:   '=',
			StringValue: "\x02\x01\x13\x13\x10\x14\x12\x0e",
			Type:        TYPE_CLUE_STRING,
		}},
		{`>9	belong  !0x0A0D1A00	game data, CORRUPTED`, "belong", &Compare{
			Not:       true,
			Operation: '=',
			IntValue:  0x0A0D1A00,
			Type:      TYPE_CLUE_INT,
		}},
		{`>>>>&1	string		x		"%s"`, "string", nil},
		{`0	search/8192	(input,`, "search", &Compare{
			Operation:   '=',
			StringValue: "(input,",
			Type:        TYPE_CLUE_STRING,
		}},
		{`>>>>>>&8	ubelong%44100	0`, "ubelong", &Compare{
			Operation: '=',
			IntValue:  0,
			Type:      TYPE_CLUE_INT,
		}},
		{`>8		ubyte/4		=0		CHN: 4`, "ubyte", &Compare{
			Operation: '=',
			IntValue:  0,
			Type:      TYPE_CLUE_INT,
		}},
		{`>>&(0x04)	lelong	>0	\b, with %d reference sequences`, "lelong", &Compare{
			Operation: '>',
			IntValue:  0,
			Type:      TYPE_CLUE_INT,
		}},
		{`>>>>(0x3C.b+0x0FF)	string	Invalid\ partition\ table		english`, "string", &Compare{
			Operation:   '=',
			StringValue: "Invalid partition table",
			Type:        TYPE_CLUE_STRING,
		}},
		{`0	string		AES`, "string", &Compare{
			Operation:   '=',
			StringValue: "AES",
			Type:        TYPE_CLUE_STRING,
		}},
		{`>>88	belong	& 1			\b, valid`, "belong", &Compare{
			Operation: '&',
			IntValue:  1,
			Type:      TYPE_CLUE_INT,
		}},
		{`0 belong 0x736C6821   Allegro datafile (packed)`, "belong", &Compare{
			Operation: '=',
			IntValue:  0x736C6821,
			Type:      TYPE_CLUE_INT,
		}},
		{`>(4.L+28)	beshort+1	>0	\b, %u type`, "beshort", &Compare{
			Operation: '>',
			IntValue:  0,
			Type:      TYPE_CLUE_INT,
		}},
		{`0	belong&0xffffe000	0x76ff2000 CDC Codec archive data`, "belong", &Compare{
			Operation: '=',
			IntValue:  0x76ff2000,
			Type:      TYPE_CLUE_INT,
		}},
		{`>2	string	\x2\x4	Xpack DiskImage archive data`, "string", &Compare{
			Operation:   '=',
			StringValue: `\x2\x4`,
			Type:        TYPE_CLUE_STRING,
		}},
		{`>0x1D5		ubequad		0x2f30313233343536	configuration of Tasmota firmware (ESP8266)`, "ubequad", &Compare{
			Operation: '=',
			QuadValue: []int64{0x2f303132, 0x33343536},
			Type:      TYPE_CLUE_QUAD,
		}},
		{`>>11		ubyte^0x65	x			\b, version %u`, "ubyte", nil},
		{`0	lelong		0x1b031336L	Netboot image,`, "lelong", &Compare{
			Operation: '=',
			IntValue:  0x1b031336,
			Type:      TYPE_CLUE_INT,
		}},
		{`>0x68	lequad		8	\b, UUID=`, "lequad", &Compare{
			Operation: '=',
			QuadValue: []int64{8, 0},
			Type:      TYPE_CLUE_QUAD,
		}},
		{`>>15	ulelong		!0x00010000h	\b, version %#8.8`, "ulelong", &Compare{
			Operation: '=',
			Not:       true,
			IntValue:  0x00010000,
			Type:      TYPE_CLUE_INT,
		}},
		{`>>>>>>(&4.l+(-4))	string		ITOLITLS	\b, Microsoft compiled help format 2.0`, "ulelong", &Compare{
			Operation:   '=',
			Type:        TYPE_CLUE_STRING,
			StringValue: "ITOLITLS",
		}},
	} {

		l := NewTest()
		err := ParseLine(l, fixture.line)
		assert.NoError(t, err, fixture.line)
		assert.Equal(t, fixture.type_, l.Type.Name, fixture.line)
		assert.Equal(t, fixture.compare, l.Compare, fixture.line)
	}

}

func TestNotSpace(t *testing.T) {
	assert.Equal(t, 4, notSpace("plop"))
	assert.Equal(t, 5, notSpace("beuha aussi"))
}
