package unpack

import (
	"encoding/binary"
	"io"

	"github.com/athoune/go-magic/model"
)

func ReadByte(r io.Reader, signed model.SIGN) (int, int64, error) {
	buff := make([]byte, 1)
	_, err := r.Read(buff)
	if err != nil {
		return 1, 0, err
	}
	// [FIXME] handle Endianness
	if signed == model.SIGNED {
		return 1, int64(buff[0]), nil
	}
	return 1, int64(buff[0]), nil
}

func ReadUShort(r io.Reader, bo binary.ByteOrder) (int, uint64, error) {
	buff := make([]byte, 2)
	_, err := r.Read(buff)
	if err != nil {
		return 1, 0, err
	}
	return 2, uint64(bo.Uint16(buff)), nil
}
func ReadShort(r io.Reader, bo binary.ByteOrder) (int, int64, error) {
	length, value, err := ReadUShort(r, bo)
	return length, int64(value), err
}

func ReadULong(r io.Reader, bo binary.ByteOrder) (int, uint64, error) {
	buff := make([]byte, 4)
	_, err := r.Read(buff)
	if err != nil {
		return 1, 0, err
	}
	return 4, uint64(bo.Uint32(buff)), nil
}

func ReadLong(r io.Reader, bo binary.ByteOrder) (int, int64, error) {
	length, value, err := ReadULong(r, bo)
	return length, int64(value), err
}

func ReadUQuad(r io.Reader, bo binary.ByteOrder) (int, uint64, error) {
	buff := make([]byte, 8)
	_, err := r.Read(buff)
	if err != nil {
		return 1, 0, err
	}
	return 8, uint64(bo.Uint64(buff)), nil
}

func ReadQuad(r io.Reader, bo binary.ByteOrder) (int, int64, error) {
	length, value, err := ReadUQuad(r, bo)
	return length, int64(value), err
}
