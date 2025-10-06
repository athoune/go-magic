package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndiannessSigned(t *testing.T) {
	e, signed, typ := ByteOrderAndSigned("ubyte")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.False(t, signed)
	assert.Equal(t, typ, "byte")
	e, signed, typ = ByteOrderAndSigned("ubelong")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.False(t, signed)
	assert.Equal(t, typ, "long")
	e, signed, typ = ByteOrderAndSigned("string")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.False(t, signed)
	assert.Equal(t, typ, "string")
	e, signed, typ = ByteOrderAndSigned("beshort")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.True(t, signed)
	assert.Equal(t, typ, "short")
}
