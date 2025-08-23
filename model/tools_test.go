package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndianessSigned(t *testing.T) {
	s, e, typ := EndianessSigned("ubyte")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.Equal(t, s, UNSIGNED)
	assert.Equal(t, typ, "byte")
	s, e, typ = EndianessSigned("ubelong")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.Equal(t, s, UNSIGNED)
	assert.Equal(t, typ, "long")
	s, e, typ = EndianessSigned("string")
	assert.Equal(t, e, NATIVE_ENDIAN)
	assert.Equal(t, s, SIGNED)
	assert.Equal(t, typ, "string")
	s, e, typ = EndianessSigned("beshort")
	assert.Equal(t, e, BIG_ENDIAN)
	assert.Equal(t, s, SIGNED)
	assert.Equal(t, typ, "short")
}
