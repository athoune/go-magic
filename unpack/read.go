package unpack

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/athoune/go-magic/model"
)

const MAXstring = 128 // max len of "string" types, matches C file.h

func ReadValue(typ *model.Type, r io.Reader) (*model.Value, int, error) {
	bo := ModelByteOrderToBinaryByteOrder(typ.ByteOrder)

	v := &model.Value{}
	switch typ.TypeFamily {
	/*
		case model.TYPE_CLUE_STRING:
			v.StringValue = string(buff)
	*/
	case model.TYPE_FAMILY_FLOAT:
		var size int
		var err error
		switch typ.Root {
		case "float":
			size = 4
			var f float32
			err = binary.Read(r, bo, &f)
			v.FloatValue = float64(f)
		case "double":
			size = 8
			err = binary.Read(r, bo, &v.FloatValue)
		}
		return v, size, err

	case model.TYPE_FAMILY_UINT:
		var err error
		var size int
		switch typ.Root {
		case "byte":
			var b byte
			err = binary.Read(r, bo, &b)
			v.UIntValue = uint64(b)
			size = 1
		case "short":
			var s uint16
			err = binary.Read(r, bo, &s)
			v.UIntValue = uint64(s)
			size = 2
		case "long":
			var l uint32
			err = binary.Read(r, bo, &l)
			v.UIntValue = uint64(l)
			size = 4
		case "quad":
			var q uint64
			err = binary.Read(r, bo, &q)
			v.UIntValue = uint64(q)
			size = 8
		default:
			return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
		}
		v.RawBinaryValue = v.UIntValue
		return v, size, err

	case model.TYPE_FAMILY_INT:
		var err error
		var size int
		switch typ.Root {
		case "byte":
			var b byte
			err = binary.Read(r, bo, &b)
			v.IntValue = int64(b)
			size = 1
		case "short":
			var s int16
			err = binary.Read(r, bo, &s)
			v.IntValue = int64(s)
			size = 2
		case "long":
			var l int32
			err = binary.Read(r, bo, &l)
			v.IntValue = int64(l)
			size = 4
		case "quad":
			var q int64
			err = binary.Read(r, bo, &q)
			v.IntValue = int64(q)
			size = 8
		default:
			return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
		}
		v.RawBinaryValue = uint64(v.IntValue)
		return v, size, err
	case model.TYPE_FAMILY_STRING:
		if typ.Root == "pstring" {
			return readPString(typ, r, bo)
		}
		// Read up to MAXstring bytes from the reader, like C's file_getbuffer
		buf := make([]byte, MAXstring)
		n, err := io.ReadAtLeast(r, buf, 1)
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, 0, fmt.Errorf("string read: %w", err)
		}
		v.StringValue = string(buf[:n])
		return v, n, nil
	default:
		return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
	}
}

// readPString reads a Pascal string: a length prefix followed by the string content.
// C reference: softmagic.c mprint() case FILE_PSTRING
func readPString(typ *model.Type, r io.Reader, bo binary.ByteOrder) (*model.Value, int, error) {
	v := &model.Value{Family: model.TYPE_FAMILY_STRING}

	// Read the length prefix (1, 2, or 4 bytes)
	prefixSize := typ.StringIntOption
	if prefixSize == 0 {
		prefixSize = 1 // default: 1-byte length
	}

	prefixBuf := make([]byte, prefixSize)
	if _, err := io.ReadFull(r, prefixBuf); err != nil {
		return nil, 0, fmt.Errorf("pstring prefix read: %w", err)
	}

	// Decode the length from the prefix
	var strLen uint64
	switch prefixSize {
	case 1:
		strLen = uint64(prefixBuf[0])
	case 2:
		strLen = uint64(bo.Uint16(prefixBuf))
	case 4:
		strLen = uint64(bo.Uint32(prefixBuf))
	default:
		return nil, 0, fmt.Errorf("unsupported pstring prefix size: %d", prefixSize)
	}

	// If PSTRING_LENGTH_INCLUDES_ITSELF, subtract the prefix size
	if typ.StringOptions&model.PSTRING_OPTIONS_SIZE_INCLUDE != 0 {
		strLen -= uint64(prefixSize)
	}

	// Cap at MAXstring like C does (sizeof(p->s) - sz)
	maxContent := MAXstring - prefixSize
	if strLen > uint64(maxContent) {
		strLen = uint64(maxContent)
	}

	// Read the string content
	if strLen == 0 {
		v.StringValue = ""
		return v, prefixSize, nil
	}

	contentBuf := make([]byte, strLen)
	n, err := io.ReadFull(r, contentBuf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, 0, fmt.Errorf("pstring content read: %w", err)
	}

	v.StringValue = string(contentBuf[:n])
	return v, prefixSize + n, nil
}
