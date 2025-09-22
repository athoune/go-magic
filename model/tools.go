package model

import "strings"

const (
	NATIVE_ENDIAN = BYTE_ORDER(iota)
	BIG_ENDIAN
	LITTLE_ENDIAN
)

type BYTE_ORDER byte

/*
ByteOrderAndSigned use name convention with txt input and
return signed/unsigned, byte order, and root name
*/
func ByteOrderAndSigned(txt string) (BYTE_ORDER, string) {
	for _, n := range []string{"name", "use"} { // something like label/goto
		if txt == n {
			return NATIVE_ENDIAN, txt
		}
	}
	switch {
	case strings.HasPrefix(txt, "ube"):
		return BIG_ENDIAN, "u" + txt[3:]
	case strings.HasPrefix(txt, "ule"):
		return LITTLE_ENDIAN, "u" + txt[3:]
	case strings.HasPrefix(txt, "u"):
		return NATIVE_ENDIAN, txt
	case strings.HasPrefix(txt, "be"):
		return BIG_ENDIAN, txt[2:]
	case strings.HasPrefix(txt, "le"):
		return LITTLE_ENDIAN, txt[2:]
	default:
		return NATIVE_ENDIAN, txt
	}
}
