package model

import "strings"

const (
	NATIVE_ENDIAN = BYTE_ORDER(iota)
	BIG_ENDIAN
	LITTLE_ENDIAN
	MIDDLE_ENDIAN // crappy old stuff
)

type BYTE_ORDER byte

/*
ByteOrderAndSigned use name convention with txt input and
return signed/unsigned, byte order, and root name
*/
func ByteOrderAndSigned(txt string) (BYTE_ORDER, bool, string) {
	switch {
	case txt == "use": // it's not an unsigned 'se'
		return NATIVE_ENDIAN, false, txt
	case strings.HasPrefix(txt, "ube"):
		return BIG_ENDIAN, false, txt[3:]
	case strings.HasPrefix(txt, "ule"):
		return LITTLE_ENDIAN, false, txt[3:]
	case strings.HasPrefix(txt, "ube"):
		return BIG_ENDIAN, false, txt[3:]
	case strings.HasPrefix(txt, "u"):
		return NATIVE_ENDIAN, false, txt[1:]
	case strings.HasPrefix(txt, "be"):
		return BIG_ENDIAN, true, txt[2:]
	case strings.HasPrefix(txt, "le"):
		return LITTLE_ENDIAN, true, txt[2:]
	case strings.HasPrefix(txt, "me"):
		return MIDDLE_ENDIAN, true, txt[2:]
	default:
		return NATIVE_ENDIAN, true, txt
	}
}
