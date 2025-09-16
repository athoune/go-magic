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
