package unpack

import (
	"fmt"
	"strconv"

	"github.com/athoune/go-magic/model"
)

func BuildValueFromString(typ *model.Type, txt string) (*model.Value, error) {
	var err error
	v := &model.Value{
		Clue: typ.Clue,
	}
	switch typ.Clue {

	case model.TYPE_CLUE_STRING:
		v.StringValue = txt
		return v, nil

	case model.TYPE_CLUE_FLOAT:
		var size int
		switch typ.Root {
		case "float":
			size = 32
		case "double":
			size = 64
		}
		v.FloatValue, err = strconv.ParseFloat(txt, size)
		return v, err

	case model.TYPE_CLUE_UINT:
		switch typ.Root {
		case "byte":
			v.UIntValue, err = strconv.ParseUint(txt, 0, 8)
		case "short":
			v.UIntValue, err = strconv.ParseUint(txt, 0, 16)
		case "long":
			v.UIntValue, err = strconv.ParseUint(txt, 0, 32)
		case "quad":
			v.UIntValue, err = strconv.ParseUint(txt, 0, 64)
		default:
			return nil, fmt.Errorf("unknown type for an unsigned integer : %s", typ.Root)
		}
		return v, err

	case model.TYPE_CLUE_INT:
		switch typ.Root {
		case "byte":
			v.IntValue, err = strconv.ParseInt(txt, 0, 8)
		case "short":
			v.IntValue, err = strconv.ParseInt(txt, 0, 16)
		case "long":
			v.IntValue, err = strconv.ParseInt(txt, 0, 32)
		case "quad":
			v.IntValue, err = strconv.ParseInt(txt, 0, 64)
		default:
			return nil, fmt.Errorf("unknown type for an integer : %s", typ.Root)
		}
		return v, err

	default:
		return nil, fmt.Errorf("unknown type: %v %v", typ.Clue, typ.Name)
	}
}
