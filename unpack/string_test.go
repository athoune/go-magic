package unpack

import (
	"bytes"
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

func TestPString(t *testing.T) {
	// PSTRING with 1-byte prefix, little endian
	// Content: [5, 'h', 'e', 'l', 'l', 'o']
	buf := bytes.NewBuffer(append([]byte{5}, []byte("hello")...))
	typ := &model.Type{
		Root:           "pstring",
		TypeFamily:     model.TYPE_FAMILY_STRING,
		StringIntOption: 1,
		ByteOrder:      model.LITTLE_ENDIAN,
	}
	v, n, err := ReadValue(typ, buf)
	assert.NoError(t, err)
	assert.Equal(t, 6, n) // 1 prefix + 5 content
	assert.Equal(t, "hello", v.StringValue)

	// PSTRING with 2-byte prefix, little endian
	buf = bytes.NewBuffer([]byte{5, 0, 'h', 'e', 'l', 'l', 'o'})
	typ = &model.Type{
		Root:           "pstring",
		TypeFamily:     model.TYPE_FAMILY_STRING,
		StringIntOption: 2,
		ByteOrder:      model.LITTLE_ENDIAN,
	}
	v, n, err = ReadValue(typ, buf)
	assert.NoError(t, err)
	assert.Equal(t, 7, n) // 2 prefix + 5 content
	assert.Equal(t, "hello", v.StringValue)

	// PSTRING with 2-byte prefix, big endian
	buf = bytes.NewBuffer([]byte{0, 5, 'h', 'e', 'l', 'l', 'o'})
	typ = &model.Type{
		Root:           "pstring",
		TypeFamily:     model.TYPE_FAMILY_STRING,
		StringIntOption: 2,
		ByteOrder:      model.BIG_ENDIAN,
	}
	v, n, err = ReadValue(typ, buf)
	assert.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, "hello", v.StringValue)

	// PSTRING with /J (length includes itself)
	// Length = 6 (includes 1-byte prefix), so content = 5 bytes
	buf = bytes.NewBuffer(append([]byte{6}, []byte("hello")...))
	typ = &model.Type{
		Root:           "pstring",
		TypeFamily:     model.TYPE_FAMILY_STRING,
		StringIntOption: 1,
		ByteOrder:      model.LITTLE_ENDIAN,
		StringOptions:  model.PSTRING_OPTIONS_SIZE_INCLUDE,
	}
	v, n, err = ReadValue(typ, buf)
	assert.NoError(t, err)
	assert.Equal(t, 6, n)
	assert.Equal(t, "hello", v.StringValue)

	// Empty pstring (length 0)
	buf = bytes.NewBuffer([]byte{0, 'x', 'x', 'x'})
	typ = &model.Type{
		Root:           "pstring",
		TypeFamily:     model.TYPE_FAMILY_STRING,
		StringIntOption: 1,
		ByteOrder:      model.LITTLE_ENDIAN,
	}
	v, n, err = ReadValue(typ, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n) // only the prefix byte
	assert.Equal(t, "", v.StringValue)
}
