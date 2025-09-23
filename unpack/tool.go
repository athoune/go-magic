package unpack

import (
	"bytes"
	"encoding/binary"
	"regexp"
	"strconv"
	"strings"

	"github.com/athoune/go-magic/model"
)

var brokenHexRe *regexp.Regexp

func init() {
	brokenHexRe = regexp.MustCompile(`^\\x[0-9a-fA-F]{1,2}`)
}

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

func HandleNumberSuffix(raw string) string {
	if strings.HasSuffix(raw, "h") ||
		strings.HasSuffix(raw, "L") {
		/*
			[FIXME]
			What the hell are this letters ?!
			 >>15	ulelong		!0x00010000h	\b, version %#8.8
			 0	lelong		0x1b031336L	Netboot image,
		*/
		return raw[:len(raw)-1]
	}
	return raw
}

func HandleStringEscape(value string) (string, error) {
	poz := 0
	out := &bytes.Buffer{}
	for {
		if poz == len(value) {
			break
		}
		switch {
		case strings.HasPrefix(value[poz:], `\x`):
			var (
				v   int64
				err error
			)
			hex := brokenHexRe.FindString(value[poz:])
			if hex != "" {
				v, err = strconv.ParseInt(hex[2:], 16, 64)
				if err != nil {
					return "", err
				}
				out.WriteByte(byte(v))
				poz += len(hex)
			}
		case poz+2 <= len(value) && value[poz:poz+2] == `\ `:
			out.WriteByte(' ')
			poz += 2
		default:
			out.WriteByte(value[poz])
			poz++
		}
	}
	return out.String(), nil
}

// padding8 add 0 to fill 8 byte
func padding8(bo model.BYTE_ORDER, value []byte) []byte {
	vv := make([]byte, 8)
	padding := 8 - len(value)
	// See https://en.wikipedia.org/wiki/Endianness
	if bo == model.LITTLE_ENDIAN {
		for i, v := range value {
			vv[len(value)-1-i] = v
		}
	} else { // Little endian
		for i, v := range value {
			vv[padding+i] = v
		}
	}
	return vv
}

// BytesAndTypeToUInt64 get 8 byte and endianness, and returns the uint64 value
func BytesAndTypeToUInt64(typ *model.Type, value []byte) (uint64, error) {
	r := bytes.NewReader(padding8(typ.ByteOrder, value))
	var v uint64
	err := binary.Read(r, ModelByteOrderToBinaryByteOrder(typ.ByteOrder), &v)
	return v, err
}
