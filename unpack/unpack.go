package unpack

/*
Unpack reads and unpacks contents.

Content can be read from `string` (from configuration files),
or from `io.Reader` (from binary contents).
*/

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/athoune/go-magic/model"
)

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
	case model.TYPE_CLUE_INT:
		if typ.Signed {
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
		} else {
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
		}
	default:
		return nil, fmt.Errorf("unknown type: %v %v", typ.Clue, typ.Name)
	}
}

func ReadToValue(typ *model.Type, r io.Reader) (*model.Value, int, error) {
	var err error
	bo := ModelByteOrderToBinaryByteOrder(typ.ByteOrder)

	v := &model.Value{}
	switch typ.Clue {
	/*
		case model.TYPE_CLUE_STRING:
			v.StringValue = string(buff)
	*/
	case model.TYPE_CLUE_FLOAT:
		var size int
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
	case model.TYPE_CLUE_INT:
		if typ.Signed {
			switch typ.Root {
			case "byte":
				var b byte
				err = binary.Read(r, bo, &b)
				v.IntValue = int64(b)
				return v, 1, err
			case "short":
				var s int16
				err = binary.Read(r, bo, &s)
				v.IntValue = int64(s)
				return v, 2, err
			case "long":
				var l int32
				err = binary.Read(r, bo, &l)
				v.IntValue = int64(l)
				return v, 4, err
			case "quad":
				var q int64
				err = binary.Read(r, bo, &q)
				v.IntValue = int64(q)
				return v, 8, err
			default:
				return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
			}
		} else {
			switch typ.Root {
			case "byte":
				var b byte
				err = binary.Read(r, bo, &b)
				v.UIntValue = uint64(b)
				return v, 1, err
			case "short":
				var s uint16
				err = binary.Read(r, bo, &s)
				v.UIntValue = uint64(s)
				return v, 2, err
			case "long":
				var l uint32
				err = binary.Read(r, bo, &l)
				v.UIntValue = uint64(l)
				return v, 4, err
			case "quad":
				var q uint64
				err = binary.Read(r, bo, &q)
				v.UIntValue = uint64(q)
				return v, 8, err
			default:
				return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
			}
		}
	default:
		return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
	}
}
