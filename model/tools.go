package model

import "strings"

const (
	SIGNED = SIGN(iota)
	UNSIGNED
	LOW_ENDIAN = ENDIANESS(iota)
	BIG_ENDIAN
	NATIVE_ENDIAN
)

type SIGN byte
type ENDIANESS byte

func EndianessSigned(txt string) (SIGN, ENDIANESS, string) {
	switch {
	case strings.HasPrefix(txt, "ube"):
		return UNSIGNED, BIG_ENDIAN, txt[3:]
	case strings.HasPrefix(txt, "ule"):
		return UNSIGNED, LOW_ENDIAN, txt[3:]
	case strings.HasPrefix(txt, "u"):
		return UNSIGNED, NATIVE_ENDIAN, txt[1:]
	default:
		return SIGNED, NATIVE_ENDIAN, txt
	}
}
