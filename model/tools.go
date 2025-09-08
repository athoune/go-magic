package model

import "strings"

const (
	LITTLE_ENDIAN = BYTE_ORDER(iota)
	BIG_ENDIAN
	NATIVE_ENDIAN
)

type BYTE_ORDER byte

/*
EndiannessSigned use name convention with txt input and
return signed/unsigned, byte order, and root name
*/
func EndiannessSigned(txt string) (bool, BYTE_ORDER, string) {
	for _, n := range []string{"name", "use"} { // something like label/goto
		if txt == n {
			return true, NATIVE_ENDIAN, txt
		}
	}
	switch {
	case strings.HasPrefix(txt, "ube"):
		return false, BIG_ENDIAN, txt[3:]
	case strings.HasPrefix(txt, "ule"):
		return false, LITTLE_ENDIAN, txt[3:]
	case strings.HasPrefix(txt, "u"):
		return false, NATIVE_ENDIAN, txt[1:]
	case strings.HasPrefix(txt, "be"):
		return true, BIG_ENDIAN, txt[2:]
	case strings.HasPrefix(txt, "le"):
		return true, LITTLE_ENDIAN, txt[2:]
	default:
		return true, NATIVE_ENDIAN, txt
	}
}
