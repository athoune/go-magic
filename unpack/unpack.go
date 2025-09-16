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

	"github.com/athoune/go-magic/model"
)

func ReadValue(typ *model.Type, r io.Reader) (*model.Value, int, error) {
	bo := ModelByteOrderToBinaryByteOrder(typ.ByteOrder)

	v := &model.Value{}
	switch typ.Clue {
	/*
		case model.TYPE_CLUE_STRING:
			v.StringValue = string(buff)
	*/
	case model.TYPE_CLUE_FLOAT:
		var size int
		var err error
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
		var err error
		var size int
		if typ.Signed {
			switch typ.Root {
			case "byte":
				var b byte
				err = binary.Read(r, bo, &b)
				v.IntValue = int64(b)
				size = 1
			case "short":
				var s int16
				err = binary.Read(r, bo, &s)
				v.IntValue = int64(s)
				size = 2
			case "long":
				var l int32
				err = binary.Read(r, bo, &l)
				v.IntValue = int64(l)
				size = 4
			case "quad":
				var q int64
				err = binary.Read(r, bo, &q)
				v.IntValue = int64(q)
				size = 8
			default:
				return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
			}
		} else {
			switch typ.Root {
			case "byte":
				var b byte
				err = binary.Read(r, bo, &b)
				v.UIntValue = uint64(b)
				size = 1
			case "short":
				var s uint16
				err = binary.Read(r, bo, &s)
				v.UIntValue = uint64(s)
				size = 2
			case "long":
				var l uint32
				err = binary.Read(r, bo, &l)
				v.UIntValue = uint64(l)
				size = 4
			case "quad":
				var q uint64
				err = binary.Read(r, bo, &q)
				v.UIntValue = uint64(q)
				size = 8
			default:
				return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
			}
		}
		return v, size, err
	default:
		return nil, 0, fmt.Errorf("wrong type: %s", typ.Root)
	}
}
