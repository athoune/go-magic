package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOffset(t *testing.T) {
	o, err := ParseOffset(">>>4")
	assert.NoError(t, err)
	assert.Equal(t, 3, o.Level)
	assert.Equal(t, int64(4), o.Value)
	assert.False(t, o.Dynamic)

	o, err = ParseOffset(">>(4.s*512)")
	assert.NoError(t, err)
	assert.Equal(t, 2, o.Level)
	assert.True(t, o.Dynamic)
	assert.Equal(t, int64(4), o.DynOffset)
	assert.Equal(t, uint8('s'), o.DynType)
	assert.Equal(t, uint8('*'), o.DynAction)
	assert.Equal(t, int64(512), o.DynArg)
}

func TestParseCompare(t *testing.T) {
	c, err := ParseCompare("0x01000007", 'i')
	assert.NoError(t, err)
	assert.Equal(t, COMPARE_EQUAL, c.Operation)
	assert.Equal(t, int64(16777223), c.IntValue)

	c, err = ParseCompare("<10", 'i')
	assert.NoError(t, err)
	assert.Equal(t, COMPARE_LESS, c.Operation)
	assert.Equal(t, int64(10), c.IntValue)

	c, err = ParseCompare("!>10", 'i')
	assert.NoError(t, err)
	assert.True(t, c.Not)
	assert.Equal(t, COMPARE_GREATER, c.Operation)
	assert.Equal(t, int64(10), c.IntValue)

	c, err = ParseCompare("D6E229D3-35DA-11D1-9034-00A0C90349BE", 's')
	assert.NoError(t, err)
	assert.Equal(t, COMPARE_EQUAL, c.Operation)
	assert.Equal(t, "D6E229D3-35DA-11D1-9034-00A0C90349BE", c.StringValue)
}
