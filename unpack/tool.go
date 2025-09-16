package unpack

import (
	"encoding/binary"

	"github.com/athoune/go-magic/model"
)

// ModelByteOrderToBinaryByteOrder convert homebrew model.BYTE_ORDER to standard binary.ByteOrder
func ModelByteOrderToBinaryByteOrder(bo model.BYTE_ORDER) binary.ByteOrder {
	switch bo {
	case model.LITTLE_ENDIAN:
		return binary.LittleEndian
	case model.BIG_ENDIAN:
		return binary.BigEndian
	default:
		return binary.NativeEndian
	}
}
