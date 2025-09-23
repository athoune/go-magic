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
		if strings.HasPrefix(txt, "-") {
			v.IntValue, err = strconv.ParseInt(txt, 0, 64)
		} else {
			// [FIXME] ParseInt does handle well larges negatives numbers like 0xB7D800203749DA11
			var uValue uint64
			uValue, err = strconv.ParseUint(txt, 0, 64)
			v.IntValue = int64(uValue)
		}
		return v, err

	default:
		return nil, fmt.Errorf("unknown type: %v %v", typ.TypeFamily, typ.Name)
	}
}
