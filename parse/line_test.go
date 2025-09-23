package parse

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	for _, fixture := range []struct {
		line    string
		Compare *model.Compare
	}{
		/*
			{`>>12	indirect/r	x`, &model.Compare{

			}},*/
		{`0	bequad	0xB7D800203749DA11`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0xB7D800203749DA11",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: int64(-5199405631432697327),
			},
			Type: &model.Type{
				TypeFamily: model.TYPE_FAMILY_INT,
				Name:       "bequad",
				Root:       "quad",
				ByteOrder:  model.BIG_ENDIAN,
			},
		}},
		{`>>16	belong&0xfe00f0f0	0x3030`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0x3030",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0x3030,
			},
			Type: &model.Type{
				TypeFamily:           model.TYPE_FAMILY_INT,
				Name:                 "belong",
				Root:                 "long",
				FilterOperator:       '&',
				ByteOrder:            model.BIG_ENDIAN,
				FilterBinaryArgument: 0xfe00f0f0,
				FilterStringArgument: "0xfe00f0f0",
			},
		}},
		{`0	lelong		0xc3cbc6c5	RISC OS Chunk data`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0xc3cbc6c5",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0xc3cbc6c5,
			},
			Type: &model.Type{
				Name:       "lelong",
				Root:       "long",
				TypeFamily: model.TYPE_FAMILY_INT,
				ByteOrder:  model.LITTLE_ENDIAN,
			}},
		},
		{`>>>>>>0	ubyte			< 10	Infocom (Z-machine %d`, &model.Compare{
			Comparator:  '<',
			RawExpected: "10",
			Expected: &model.Value{
				Family:    model.TYPE_FAMILY_UINT,
				UIntValue: 10,
			},
			Type: &model.Type{
				Name:       "ubyte",
				Root:       "ubyte",
				TypeFamily: model.TYPE_FAMILY_UINT,
			},
		}},
		{`0	string 		Draw		RISC OS Draw file data`, &model.Compare{
			Comparator:  '=',
			RawExpected: "Draw",
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "Draw",
			},
			Type: &model.Type{
				Name:       "string",
				Root:       "string",
				TypeFamily: model.TYPE_FAMILY_STRING,
			},
		}},
		{`>12	leshort	!1	%d patterns`, &model.Compare{
			Not:         true,
			Comparator:  '=',
			RawExpected: "1",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 1,
			},
			Type: &model.Type{
				Name:       "leshort",
				Root:       "short",
				ByteOrder:  model.LITTLE_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_INT,
			},
		}},
		{`0	string	\x02\x01\x13\x13\x10\x14\x12\x0e`, &model.Compare{
			Comparator:  '=',
			RawExpected: `\x02\x01\x13\x13\x10\x14\x12\x0e`,
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "\x02\x01\x13\x13\x10\x14\x12\x0e",
			},
			Type: &model.Type{
				Name:       "string",
				Root:       "string",
				TypeFamily: model.TYPE_FAMILY_STRING,
			},
		}},
		{`>9	belong  !0x0A0D1A00	game data, CORRUPTED`, &model.Compare{
			Not:         true,
			Comparator:  '=',
			RawExpected: "0x0A0D1A00",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0x0A0D1A00,
			},
			Type: &model.Type{
				Name:       "belong",
				Root:       "long",
				ByteOrder:  model.BIG_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_INT,
			},
		}},
		{`>>>>&1	string		x		"%s"`, &model.Compare{
			X:        true,
			Expected: nil,
			Type: &model.Type{
				Name:       "string",
				Root:       "string",
				TypeFamily: model.TYPE_FAMILY_STRING,
			},
		}},
		{`0	search/8192	(input,`, &model.Compare{
			Comparator:  '=',
			RawExpected: `(input,`,
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "(input,",
			},
			Type: &model.Type{
				Name:                 "search",
				Root:                 "search",
				TypeFamily:           model.TYPE_FAMILY_STRING,
				FilterOperator:       '/',
				FilterStringArgument: "8192",
				SearchRange:          8192,
			},
		}},
		{`>>>>>>&8	ubelong%44100	0`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0",
			Expected: &model.Value{
				Family:    model.TYPE_FAMILY_UINT,
				UIntValue: 0,
			},
			Type: &model.Type{
				Name:                 "ubelong",
				Root:                 "ulong",
				ByteOrder:            model.BIG_ENDIAN,
				TypeFamily:           model.TYPE_FAMILY_UINT,
				FilterOperator:       '%',
				FilterBinaryArgument: 44100,
				FilterStringArgument: "44100",
			},
		}},
		{`>8		ubyte/4		=0		CHN: 4`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0",
			Expected: &model.Value{
				Family:    model.TYPE_FAMILY_UINT,
				UIntValue: 0,
			},
			Type: &model.Type{
				Name:                 "ubyte",
				Root:                 "ubyte",
				TypeFamily:           model.TYPE_FAMILY_UINT,
				FilterOperator:       '/',
				FilterBinaryArgument: 4,
				FilterStringArgument: "4",
			},
		}},
		{`>>&(0x04)	lelong	>0	\b, with %d reference sequences`, &model.Compare{
			Comparator:  '>',
			RawExpected: "0",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0,
			},
			Type: &model.Type{
				Name:       "lelong",
				Root:       "long",
				ByteOrder:  model.LITTLE_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_INT},
		}},
		{`>>>>(0x3C.b+0x0FF)	string	Invalid\ partition\ table		english`, &model.Compare{
			Comparator:  '=',
			RawExpected: `Invalid\ partition\ table`,
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "Invalid partition table",
			},
			Type: &model.Type{
				Name:       "string",
				Root:       "string",
				TypeFamily: model.TYPE_FAMILY_STRING},
		}},
		{`0	string		AES`, &model.Compare{
			Comparator:  '=',
			RawExpected: "AES",
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "AES",
			},
			Type: &model.Type{
				Name:       "string",
				Root:       "string",
				TypeFamily: model.TYPE_FAMILY_STRING},
		}},
		{`>>88	belong	& 1			\b, valid`, &model.Compare{
			Comparator:  '&',
			RawExpected: "1",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 1,
			},
			Type: &model.Type{
				Name:       "belong",
				Root:       "long",
				ByteOrder:  model.BIG_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_INT},
		}},
		{`0 belong 0x736C6821   Allegro datafile (packed)`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0x736C6821",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0x736C6821,
			},
			Type: &model.Type{
				Name:       "belong",
				Root:       "long",
				ByteOrder:  model.BIG_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_INT},
		}},
		{`>(4.L+28)	beshort+1	>0	\b, %u type`, &model.Compare{
			Comparator:  '>',
			RawExpected: "0",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0,
			},
			Type: &model.Type{
				Name:                 "beshort",
				Root:                 "short",
				ByteOrder:            model.BIG_ENDIAN,
				TypeFamily:           model.TYPE_FAMILY_INT,
				FilterOperator:       '+',
				FilterBinaryArgument: 1,
				FilterStringArgument: "1",
			},
		}},
		{`0	belong&0xffffe000	0x76ff2000 CDC Codec archive data`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0x76ff2000",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0x76ff2000,
			},
			Type: &model.Type{
				Name:                 "belong",
				Root:                 "long",
				ByteOrder:            model.BIG_ENDIAN,
				TypeFamily:           model.TYPE_FAMILY_INT,
				FilterOperator:       '&',
				FilterBinaryArgument: 0xffffe000,
				FilterStringArgument: "0xffffe000",
			},
		}},
		{`>2	string	\x2\x4	Xpack DiskImage archive data`, &model.Compare{
			Comparator:  '=',
			RawExpected: `\x2\x4`,
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: string([]byte{0x2, 0x4}),
			},
			Type: &model.Type{
				Name:       "string",
				Root:       "string",
				TypeFamily: model.TYPE_FAMILY_STRING},
		}},
		{`>0x1D5		ubequad		0x2f30313233343536	configuration of Tasmota firmware (ESP8266)`,
			&model.Compare{
				Comparator:  '=',
				RawExpected: "0x2f30313233343536",
				Expected: &model.Value{
					Family:    model.TYPE_FAMILY_UINT,
					UIntValue: 0x2f30313233343536,
				},
				Type: &model.Type{
					Name:       "ubequad",
					Root:       "uquad",
					ByteOrder:  model.BIG_ENDIAN,
					TypeFamily: model.TYPE_FAMILY_UINT,
				},
			}},
		{`>>11		ubyte^0x65	x			\b, version %u`, &model.Compare{
			X: true,
			Type: &model.Type{
				Name:                 "ubyte",
				Root:                 "ubyte",
				TypeFamily:           model.TYPE_FAMILY_UINT,
				FilterOperator:       '^',
				FilterBinaryArgument: 0x65,
				FilterStringArgument: "0x65",
			},
			Expected: nil,
		}},
		{`0	lelong		0x1b031336L	Netboot image,`, &model.Compare{
			Comparator:  '=',
			RawExpected: "0x1b031336L",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 0x1b031336,
			},
			Type: &model.Type{
				Name:       "lelong",
				Root:       "long",
				ByteOrder:  model.LITTLE_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_INT},
		}},
		{`>0x68	lequad		8	\b, UUID=`, &model.Compare{
			Comparator:  '=',
			RawExpected: "8",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 8,
			},
			Type: &model.Type{
				Name:       "lequad",
				Root:       "quad",
				ByteOrder:  model.LITTLE_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_INT,
			},
		}},
		{`>>15	ulelong		!0x00010000h	\b, version %#8.8`, &model.Compare{
			Comparator:  '=',
			Not:         true,
			RawExpected: "0x00010000h",
			Expected: &model.Value{
				Family:    model.TYPE_FAMILY_UINT,
				UIntValue: 0x00010000,
			},
			Type: &model.Type{
				Name:       "ulelong",
				Root:       "ulong",
				ByteOrder:  model.LITTLE_ENDIAN,
				TypeFamily: model.TYPE_FAMILY_UINT,
			},
		}},
		{`>>>>>>(&4.l+(-4))	string		ITOLITLS	\b, Microsoft compiled help format 2.0`,
			&model.Compare{
				Comparator:  '=',
				RawExpected: "ITOLITLS",
				Expected: &model.Value{
					Family:      model.TYPE_FAMILY_STRING,
					StringValue: "ITOLITLS",
				},
				Type: &model.Type{
					Name:       "string",
					Root:       "string",
					TypeFamily: model.TYPE_FAMILY_STRING,
				},
			}},
		{`0	string	zz	MGR bitmap, old format, 1-bit deep, 16-bit aligned`,
			&model.Compare{
				Comparator:  '=',
				RawExpected: "zz",
				Expected: &model.Value{
					Family:      model.TYPE_FAMILY_STRING,
					StringValue: "zz",
				},
				Type: &model.Type{
					Name:       "string",
					Root:       "string",
					TypeFamily: model.TYPE_FAMILY_STRING,
				},
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
	assert.Equal(t, byte('&'), te.Type.FilterOperator)
	assert.Equal(t, uint64(0xffffff00), te.Type.FilterBinaryArgument)
	assert.Equal(t, "belong", te.Type.Name)
}

func TestDisplayable(t *testing.T) {
	te := model.NewTest()
	err := ParseLine(te, `0	belong&0xffffff00	0xffd8ff00	JPEG image data`)
	assert.NoError(t, err)
	assert.Equal(t, "JPEG image data", te.Message.Value)
}
