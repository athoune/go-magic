package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndianessSigned(t *testing.T) {
	e, typ := ByteOrderAndSigned("ubyte")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.Equal(t, typ, "ubyte")
	e, typ = ByteOrderAndSigned("ubelong")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.Equal(t, typ, "ulong")
	e, typ = ByteOrderAndSigned("string")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.Equal(t, typ, "string")
	e, typ = ByteOrderAndSigned("beshort")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.Equal(t, typ, "short")
}
