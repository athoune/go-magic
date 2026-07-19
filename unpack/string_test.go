package unpack

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestStringValue(t *testing.T) {
	typ := &model.Type{
		Root:       "long",
		TypeFamily: model.TYPE_FAMILY_INT,
	}
	v, err := BuildValueFromString(typ, "42")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), v.IntValue)

	v, err = BuildValueFromString(typ, "0xF")
	assert.NoError(t, err)
	assert.Equal(t, int64(15), v.IntValue)

	// Negative decimal
	v, err = BuildValueFromString(typ, "-1")
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), v.IntValue)

	// Large hex > MaxInt64: 0xB7D800203749DA11
	// In C, this is stored as uint64 then cast to int64 (two's complement)
	var largeHex uint64 = 0xB7D800203749DA11
	v, err = BuildValueFromString(typ, "0xB7D800203749DA11")
	assert.NoError(t, err)
	assert.Equal(t, int64(largeHex), v.IntValue)
	assert.True(t, v.IntValue < 0, "large hex should be negative in two's complement")

	// Max int64
	v, err = BuildValueFromString(typ, "0x7FFFFFFFFFFFFFFF")
	assert.NoError(t, err)
	assert.Equal(t, int64(0x7FFFFFFFFFFFFFFF), v.IntValue)

	// Min int64 (0x8000000000000000)
	v, err = BuildValueFromString(typ, "0x8000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, int64(-1<<63), v.IntValue)

	// Negative hex with explicit minus
	v, err = BuildValueFromString(typ, "-0x10")
	assert.NoError(t, err)
	assert.Equal(t, int64(-16), v.IntValue)
}
