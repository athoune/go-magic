package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndianessSigned(t *testing.T) {
	s, e, typ := ByteOrderAndSigned("ubyte")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.False(t, s)
	assert.Equal(t, typ, "byte")
	s, e, typ = ByteOrderAndSigned("ubelong")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.False(t, s)
	assert.Equal(t, typ, "long")
	s, e, typ = ByteOrderAndSigned("string")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.True(t, s)
	assert.Equal(t, typ, "string")
	s, e, typ = ByteOrderAndSigned("beshort")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.True(t, s)
	assert.Equal(t, typ, "short")
}
