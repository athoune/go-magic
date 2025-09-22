package unpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestUnpackSignedByte(t *testing.T) {
	typ := &model.Type{
		Root: "byte",
	}
	v, l, err := ReadValue(typ, bytes.NewReader([]byte{12, 14, 3}))
	assert.NoError(t, err)
	assert.Equal(t, 1, l)
	assert.Equal(t, int64(12), v.IntValue)
}
func TestUnpackSigned(t *testing.T) {
	vv := []byte{12, 14, 3, 27, 53, 254, 7, 9, 27, 128}
	for i, fixture := range []struct {
		endianness model.BYTE_ORDER
		expected   int64
		root       string
		data       []byte
		size       int
	}{
		// Short
		{
			endianness: model.LITTLE_ENDIAN,
			expected:   int64(12 + 14*256),
			root:       "short",
			data:       []byte{12, 14, 3},
			size:       2,
		},
		{
			endianness: model.BIG_ENDIAN,
			expected:   int64(14 + 12*256),
			root:       "short",
			data:       []byte{12, 14, 3},
			size:       2,
		},
		// long
		{
			endianness: model.LITTLE_ENDIAN,
			root:       "long",
			data:       []byte{12, 14, 3, 27, 53, 254},
			expected:   int64(12 + 14*256 + 3*256*256 + 27*256*256*256),
			size:       4,
		},
		{
			endianness: model.BIG_ENDIAN,
			root:       "long",
			data:       []byte{12, 14, 3, 27, 53, 27},
			expected:   int64(12*256*256*256 + 14*256*256 + 3*256 + 27),
			size:       4,
		},
		// quad
		{
			endianness: model.LITTLE_ENDIAN,
			root:       "quad",
			data:       vv,
			expected:   bin2int(vv[:8], binary.LittleEndian),
			size:       8,
		},
		{
			endianness: model.BIG_ENDIAN,
			root:       "quad",
			data:       vv,
			expected:   bin2int(vv[:8], binary.BigEndian),
			size:       8,
		},
	} {
		typ := &model.Type{
			ByteOrder: fixture.endianness,
			Root:      fixture.root,
		}
		v, l, err := ReadValue(typ, bytes.NewReader(fixture.data))
		assert.NoError(t, err)
		assert.Equal(t, fixture.size, l)
		fmt.Println(i, v.IntValue)
		assert.Equal(t, fixture.expected, v.IntValue, "#%v : %v %v => %v",
			i, fixture.endianness, fixture.data, fixture.expected)
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
