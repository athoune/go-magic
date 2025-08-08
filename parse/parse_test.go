package parse

import (
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
		clue    model.Clue
		size    int
		compare *model.Compare
	}{
		{"<10", model.TYPE_CLUE_INT, 3, &model.Compare{
			Operation: COMPARE_LESS,
			IntValue:  int64(10),
		}},
		{"< 10", model.TYPE_CLUE_INT, 4, &model.Compare{
			Operation: COMPARE_LESS,
			IntValue:  int64(10),
		}},
		{"0x01000007", model.TYPE_CLUE_INT, 10, &model.Compare{
			Operation: COMPARE_EQUAL,
			IntValue:  int64(16777223),
		}},
		{"!>10", model.TYPE_CLUE_INT, 4, &model.Compare{
			Operation: COMPARE_GREATER,
			IntValue:  int64(10),
			Not:       true,
		}},
		{"D6E229D3-35DA-11D1-9034-00A0C90349BE", model.TYPE_CLUE_STRING, 36, &model.Compare{
			Operation:   COMPARE_EQUAL,
			StringValue: "D6E229D3-35DA-11D1-9034-00A0C90349BE",
		}},
		{`Invalid\ partition\ table		english`, model.TYPE_CLUE_STRING, 25, &model.Compare{
			Operation:   COMPARE_EQUAL,
			StringValue: "Invalid partition table",
		}},
	} {
		c, s, err := ParseCompare(fixture.line, model.Clue(fixture.clue))
		assert.NoError(t, err)
		fixture.compare.Type = fixture.clue // yes, it's cheating
		assert.Equal(t, fixture.compare, c, fixture.line)
		assert.Equal(t, fixture.size, s, fixture.line)
	}
}

func TestParseType(t *testing.T) {

}

func TestHandleStringEscape(t *testing.T) {
	for _, fixture := range []struct {
		raw     string
		escaped string
	}{
		{`\x02\x01\x13\x13\x10\x14\x12\x0e`, "\x02\x01\x13\x13\x10\x14\x12\x0e"},
		{`plop`, "plop"},
		{`\x8aMNG`, "\x8aMNG"},
		{`Beuha\ aussi`, "Beuha aussi"},
	} {
		s, err := HandleStringEscape(fixture.raw)
		assert.NoError(t, err)
		assert.Equal(t, fixture.escaped, s)
	}
}

func TestHandleSpaceEscape(t *testing.T) {
	fixture := `Invalid\ partition\ table  `
	end := HandleSpaceEscape(fixture)
	assert.Equal(t, `Invalid\ partition\ table`, fixture[:end])
	fixture = `Invalid\ partition\ table`
	end = HandleSpaceEscape(fixture)
	assert.Equal(t, `Invalid\ partition\ table`, fixture[:end])
}
