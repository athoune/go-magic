package unpack

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestUnpackSignedByte(t *testing.T) {
	l, v, err := ReadByte(bytes.NewBuffer([]byte{12, 14, 3}), model.SIGNED)
	assert.NoError(t, err)
	assert.Equal(t, 1, l)
	assert.Equal(t, int64(12), v)
}
func TestUnpackSignedShort(t *testing.T) {
	l, v, err := ReadShort(bytes.NewBuffer([]byte{12, 14, 3}), binary.BigEndian)
	assert.NoError(t, err)
	assert.Equal(t, 2, l)
	assert.Equal(t, int64(12*256+14), v)

	l, v, err = ReadShort(bytes.NewBuffer([]byte{12, 14, 3}), binary.LittleEndian)
	assert.NoError(t, err)
	assert.Equal(t, 2, l)
	assert.Equal(t, int64(14*256+12), v)
}

func TestUnpackSignedLong(t *testing.T) {
	l, v, err := ReadLong(bytes.NewBuffer([]byte{12, 14, 3, 27, 53, 254}), binary.BigEndian)
	assert.NoError(t, err)
	assert.Equal(t, 4, l)
	assert.Equal(t, int64(12*16777216+14*65536+3*256+27), v)

	l, v, err = ReadLong(bytes.NewBuffer([]byte{12, 14, 3, 27, 53, 254}), binary.LittleEndian)
	assert.NoError(t, err)
	assert.Equal(t, 4, l)
	assert.Equal(t, int64(27*16777216+3*65536+14*256+12), v)
}

func TestUnpackSignedQuad(t *testing.T) {
	vv := []byte{12, 14, 3, 27, 53, 254, 7, 9, 27, 128}
	for _, bo := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		l, v, err := ReadQuad(bytes.NewBuffer(vv), bo)
		assert.NoError(t, err)
		assert.Equal(t, 8, l)
		assert.Equal(t, bin2int(vv[:8], bo), v, bo)
	}
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
