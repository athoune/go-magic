package parse

import (
	"testing"

	"github.com/athoune/go-magic/model"
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
		Compare *model.Compare
	}{
		{`>>16	belong&0xfe00f0f0	0x3030`, &model.Compare{
			Operation: '=',
			IntValue:  0x3030,
			Type: &model.Type{
				Clue_:    model.TYPE_CLUE_INT,
				Name:     "belong",
				Operator: '&',
				Arg:      "0xfe00f0f0",
			},
		}},
		{`0	lelong		0xc3cbc6c5	RISC OS Chunk data`, &model.Compare{
			Operation: '=',
			IntValue:  0xc3cbc6c5,
			Type: &model.Type{
				Name:  "lelong",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`>>>>>>0	ubyte			< 10	Infocom (Z-machine %d`, &model.Compare{
			Operation: '<',
			IntValue:  10,
			Type: &model.Type{Name: "ubyte",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`0	string 		Draw		RISC OS Draw file data`, &model.Compare{
			Operation:   '=',
			StringValue: "Draw",
			Type: &model.Type{
				Name:  "string",
				Clue_: model.TYPE_CLUE_STRING},
		}},
		{`>12	leshort	!1	%d patterns`, &model.Compare{
			Not:       true,
			Operation: '=',
			IntValue:  1,
			Type: &model.Type{
				Name:  "leshort",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`0	string	\x02\x01\x13\x13\x10\x14\x12\x0e`, &model.Compare{
			Operation:   '=',
			StringValue: "\x02\x01\x13\x13\x10\x14\x12\x0e",
			Type: &model.Type{
				Name:  "string",
				Clue_: model.TYPE_CLUE_STRING},
		}},
		{`>9	belong  !0x0A0D1A00	game data, CORRUPTED`, &model.Compare{
			Not:       true,
			Operation: '=',
			IntValue:  0x0A0D1A00,
			Type: &model.Type{
				Name:  "belong",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`>>>>&1	string		x		"%s"`, &model.Compare{
			X: true,
			Type: &model.Type{
				Name:  "string",
				Clue_: model.TYPE_CLUE_STRING,
			},
		}},
		{`0	search/8192	(input,`, &model.Compare{
			Operation:   '=',
			StringValue: "(input,",
			Type: &model.Type{
				Name:     "search",
				Clue_:    model.TYPE_CLUE_STRING,
				Operator: '/',
				Arg:      "8192",
			},
		}},
		{`>>>>>>&8	ubelong%44100	0`, &model.Compare{
			Operation: '=',
			IntValue:  0,
			Type: &model.Type{
				Name:     "ubelong",
				Clue_:    model.TYPE_CLUE_INT,
				Operator: '%',
				Arg:      "44100",
			},
		}},
		{`>8		ubyte/4		=0		CHN: 4`, &model.Compare{
			Operation: '=',
			IntValue:  0,
			Type: &model.Type{
				Name:     "ubyte",
				Clue_:    model.TYPE_CLUE_INT,
				Operator: '/',
				Arg:      "4",
			},
		}},
		{`>>&(0x04)	lelong	>0	\b, with %d reference sequences`, &model.Compare{
			Operation: '>',
			IntValue:  0,
			Type: &model.Type{
				Name:  "lelong",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`>>>>(0x3C.b+0x0FF)	string	Invalid\ partition\ table		english`, &model.Compare{
			Operation:   '=',
			StringValue: "Invalid partition table",
			Type: &model.Type{
				Name:  "string",
				Clue_: model.TYPE_CLUE_STRING},
		}},
		{`0	string		AES`, &model.Compare{
			Operation:   '=',
			StringValue: "AES",
			Type: &model.Type{
				Name:  "string",
				Clue_: model.TYPE_CLUE_STRING},
		}},
		{`>>88	belong	& 1			\b, valid`, &model.Compare{
			Operation: '&',
			IntValue:  1,
			Type: &model.Type{
				Name:  "belong",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`0 belong 0x736C6821   Allegro datafile (packed)`, &model.Compare{
			Operation: '=',
			IntValue:  0x736C6821,
			Type: &model.Type{
				Name:  "belong",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`>(4.L+28)	beshort+1	>0	\b, %u type`, &model.Compare{
			Operation: '>',
			IntValue:  0,
			Type: &model.Type{
				Name:     "beshort",
				Clue_:    model.TYPE_CLUE_INT,
				Operator: '+',
				Arg:      "1",
			},
		}},
		{`0	belong&0xffffe000	0x76ff2000 CDC Codec archive data`, &model.Compare{
			Operation: '=',
			IntValue:  0x76ff2000,
			Type: &model.Type{
				Name:     "belong",
				Clue_:    model.TYPE_CLUE_INT,
				Operator: '&',
				Arg:      "0xffffe000",
			},
		}},
		{`>2	string	\x2\x4	Xpack DiskImage archive data`, &model.Compare{
			Operation:   '=',
			StringValue: `\x2\x4`,
			Type: &model.Type{
				Name:  "string",
				Clue_: model.TYPE_CLUE_STRING},
		}},
		{`>0x1D5		ubequad		0x2f30313233343536	configuration of Tasmota firmware (ESP8266)`,
			&model.Compare{
				Operation: '=',
				QuadValue: []int64{0x2f303132, 0x33343536},
				Type: &model.Type{
					Name:  "ubequad",
					Clue_: model.TYPE_CLUE_QUAD},
			}},
		{`>>11		ubyte^0x65	x			\b, version %u`, &model.Compare{
			X: true,
			Type: &model.Type{
				Name:     "ubyte",
				Clue_:    model.TYPE_CLUE_INT,
				Operator: '^',
				Arg:      "0x65",
			},
		}},
		{`0	lelong		0x1b031336L	Netboot image,`, &model.Compare{
			Operation: '=',
			IntValue:  0x1b031336,
			Type: &model.Type{
				Name:  "lelong",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`>0x68	lequad		8	\b, UUID=`, &model.Compare{
			Operation: '=',
			QuadValue: []int64{8, 0},
			Type: &model.Type{
				Name:  "lequad",
				Clue_: model.TYPE_CLUE_QUAD},
		}},
		{`>>15	ulelong		!0x00010000h	\b, version %#8.8`, &model.Compare{
			Operation: '=',
			Not:       true,
			IntValue:  0x00010000,
			Type: &model.Type{
				Name:  "ulelong",
				Clue_: model.TYPE_CLUE_INT},
		}},
		{`>>>>>>(&4.l+(-4))	string		ITOLITLS	\b, Microsoft compiled help format 2.0`,
			&model.Compare{
				Operation:   '=',
				StringValue: "ITOLITLS",
				Type: &model.Type{
					Name:  "string",
					Clue_: model.TYPE_CLUE_STRING},
			}},
		{`0	string	zz	MGR bitmap, old format, 1-bit deep, 16-bit aligned`,
			&model.Compare{
				Operation:   '=',
				StringValue: "zz",
				Type: &model.Type{
					Name:  "string",
					Clue_: model.TYPE_CLUE_STRING},
			}},
	} {

		l := model.NewTest()
		err := ParseLine(l, fixture.line)
		assert.NoError(t, err, fixture.line)
		assert.NotNil(t, l.Compare)
		assert.Equal(t, fixture.Compare, l.Compare, fixture.line)
	}
}

func TestNotSpace(t *testing.T) {
	assert.Equal(t, 4, notSpace("plop"))
	assert.Equal(t, 5, notSpace("beuha aussi"))
}

func TestTypeFilter(t *testing.T) {
	te := model.NewTest()
	err := ParseLine(te, `0	belong&0xffffff00	0xffd8ff00	JPEG image data`)
	assert.NoError(t, err)
	assert.Equal(t, byte('&'), te.Type.Operator)
	assert.Equal(t, "0xffffff00", te.Type.Arg)
	assert.Equal(t, "belong", te.Type.Name)
}
