package unpack

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestModelByteOrderToBinaryByteOrder(t *testing.T) {
	assert.Equal(t, binary.LittleEndian, ModelByteOrderToBinaryByteOrder(model.LITTLE_ENDIAN))
	assert.Equal(t, binary.BigEndian, ModelByteOrderToBinaryByteOrder(model.BIG_ENDIAN))
	assert.Equal(t, binary.NativeEndian, ModelByteOrderToBinaryByteOrder(model.NATIVE_ENDIAN))
}

func TestUnpackSignedByte(t *testing.T) {
	typ := &model.Type{
		Signed: true,
		Root:   "byte",
	}
	v, l, err := ReadToValue(typ, bytes.NewReader([]byte{12, 14, 3}))
	assert.NoError(t, err)
	assert.Equal(t, 1, l)
	assert.Equal(t, int64(12), v.IntValue)
}
func TestUnpackSignedShort(t *testing.T) {
	typ := &model.Type{
		Endianness: model.LITTLE_ENDIAN,
		Root:       "short",
		Signed:     true,
	}
	v, l, err := ReadToValue(typ, bytes.NewReader([]byte{12, 14, 3}))
	assert.NoError(t, err)
	assert.Equal(t, 2, l)
	assert.Equal(t, int64(12+14*256), v.IntValue)

	typ.Endianness = model.BIG_ENDIAN

	v, l, err = ReadToValue(typ, bytes.NewReader([]byte{12, 14, 3}))
	assert.NoError(t, err)
	assert.Equal(t, 2, l)
	assert.Equal(t, int64(12+14*256), v.IntValue)
}

func TestUnpackSignedLong(t *testing.T) {
	typ := &model.Type{
		Root:       "long",
		Endianness: model.LITTLE_ENDIAN,
		Signed:     true,
	}
	v, l, err := ReadToValue(typ, bytes.NewReader([]byte{12, 14, 3, 27, 53, 254}))
	assert.NoError(t, err)
	assert.Equal(t, 4, l)
	assert.Equal(t, int64(12+14*256+3*256*256+27*256*256*256), v.IntValue)

	typ.Endianness = model.BIG_ENDIAN
	v, l, err = ReadToValue(typ, bytes.NewReader([]byte{12, 14, 3, 27, 53, 254}))
	assert.NoError(t, err)
	assert.Equal(t, 4, l)
	assert.Equal(t, int64(27*256*256*256+3*256*256+14*256+12), v.IntValue)
}

func TestUnpackSignedQuad(t *testing.T) {
	typ := &model.Type{
		Endianness: model.LITTLE_ENDIAN,
		Root:       "quad",
		Signed:     true,
	}
	vv := []byte{12, 14, 3, 27, 53, 254, 7, 9, 27, 128}
	v, l, err := ReadToValue(typ, bytes.NewReader(vv))
	assert.NoError(t, err)
	assert.Equal(t, 8, l)
	assert.Equal(t, bin2int(vv, modelByteOrderToBinaryByteOrder(typ.Endianness)), v.IntValue)
}

func bin2int(bytes_ []byte, bo binary.ByteOrder) int64 {
	v := 0
	if bo == binary.LittleEndian {
		for exp, b := range bytes_ {
			v += int(b) * intPow(256, exp)
		}
	} else {
		for i, b := range bytes_ {
			ii := len(bytes_) - i - 1
			v += int(b) * intPow(256, ii)
		}
	}
	return int64(v)
}

func intPow(base, exp int) int {
	result := 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}
	return result
}
