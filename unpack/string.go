package unpack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/athoune/go-magic/model"
)

func BuildValueFromString(typ *model.Type, txt string) (*model.Value, error) {
	var err error
	v := &model.Value{
		Family: typ.TypeFamily,
	}
	switch typ.TypeFamily {

	case model.TYPE_FAMILY_STRING:
		v.StringValue, err = HandleStringEscape(txt)
		return v, err

	case model.TYPE_FAMILY_FLOAT:
		var size int
		switch typ.Root {
		case "float":
			size = 32
		case "double":
			size = 64
		}
		v.FloatValue, err = strconv.ParseFloat(txt, size)
		return v, err

	case model.TYPE_FAMILY_UINT:
		txt = HandleNumberSuffix(txt)
		v.UIntValue, err = strconv.ParseUint(txt, 0, 64)
		return v, err

	case model.TYPE_FAMILY_INT:
		txt = HandleNumberSuffix(txt)
		negative := false
		if strings.HasPrefix(txt, "-") {
			negative = true
			txt = txt[1:]
		}
		// Always parse as uint64 first, then cast to int64.
		// This matches C behavior where the value is stored as uint64
		// in the VALUETYPE union and sign-extended later by file_signextend().
		// strconv.ParseInt fails for hex values > MaxInt64 (e.g. 0xB7D800203749DA11).
		var uValue uint64
		uValue, err = strconv.ParseUint(txt, 0, 64)
		v.IntValue = int64(uValue)
		if negative {
			v.IntValue = -v.IntValue
		}
		return v, err

	default:
		return nil, fmt.Errorf("unknown type: %v %v", typ.TypeFamily, typ.Name)
	}
}
