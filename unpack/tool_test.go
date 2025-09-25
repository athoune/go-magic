package unpack

import (
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
func TestPadding(t *testing.T) {
	vv := padding8(model.BIG_ENDIAN, []byte{0x2a, 0x3})
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0x2a, 0x3}, vv)

	vv = padding8(model.LITTLE_ENDIAN, []byte{0x2a, 0x3})
	assert.Equal(t, []byte{0x3, 0x2a, 0, 0, 0, 0, 0, 0}, vv)
}

func TestHandleStringEscape(t *testing.T) {
	for _, fixture := range []struct {
		raw     string
		escaped string
	}{
		{`\x02\x01\x13\x13\x10\x14\x12\x0e`, "\x02\x01\x13\x13\x10\x14\x12\x0e"},
		{`plop`, "plop"},
		{`\x8aMNG`, "\x8aMNG"},
		{`Beuha\ aussi`, "Beuha aussi"},
		{`\x2\x4`, string([]byte{2, 4})},
	} {
		s, err := HandleStringEscape(fixture.raw)
		assert.NoError(t, err)
		assert.Equal(t, fixture.escaped, s, "%v => %v", fixture.raw, s)
	}
}

func TestMiddleBigEndian(t *testing.T) {
	assert.Equal(t, []byte{0xB7, 0xA0, 0x08, 0x07},
		middleBigEndian([]byte{0xA0, 0xB7, 0x07, 0x08}))
}
