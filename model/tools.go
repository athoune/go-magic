package model

import "strings"

const (
	SIGNED = SIGN(iota)
	UNSIGNED
	LITTLE_ENDIAN = BYTE_ORDER(iota)
	BIG_ENDIAN
	NATIVE_ENDIAN
)

type SIGN byte
type BYTE_ORDER byte

func EndianessSigned(txt string) (SIGN, BYTE_ORDER, string) {
	for _, n := range []string{"name", "use"} { // something like label/goto
		if txt == n {
			return SIGNED, NATIVE_ENDIAN, txt
		}
	}
	switch {
	case strings.HasPrefix(txt, "ube"):
		return UNSIGNED, BIG_ENDIAN, txt[3:]
	case strings.HasPrefix(txt, "ule"):
		return UNSIGNED, LITTLE_ENDIAN, txt[3:]
	case strings.HasPrefix(txt, "u"):
		return UNSIGNED, NATIVE_ENDIAN, txt[1:]
	case strings.HasPrefix(txt, "be"):
		return SIGNED, BIG_ENDIAN, txt[2:]
	case strings.HasPrefix(txt, "le"):
		return SIGNED, LITTLE_ENDIAN, txt[2:]
	default:
		return SIGNED, NATIVE_ENDIAN, txt
	}
}
