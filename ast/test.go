package ast

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/athoune/go-magic/model"
)

var SIZES map[string]int

func init() {
	SIZES = map[string]int{
		"byte":  1,
		"short": 2,
		"long":  4,
		"quad":  8,
	}
}

type Test struct {
	test   *model.Test
	target io.ReadSeeker
}

func NewTest(test *model.Test) *Test {
	return &Test{
		test: test,
	}
}

func (t *Test) Test(target io.ReadSeeker) (string, error) {
	_, err := target.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}
	t.target = target
	t.offset()
	c, _, err := t.compare()
	if err != nil {
		return "", err
	}
	if c {
		fmt.Println(t.test.Raw)
	}
	return "", nil
}

// offset seeks the target
func (t *Test) offset() {
	if t.test.Offset.Relative {
		// TODO
	} else {
		t.target.Seek(t.test.Offset.Value, io.SeekCurrent)
	}
}

// compare return compare result, is special and error
func (t *Test) compare() (bool, bool, error) {
	signed, byteOrder, typ := model.EndianessSigned(t.test.Type.Name)

	var bo binary.ByteOrder
	switch byteOrder {
	case model.BIG_ENDIAN:
		bo = binary.BigEndian
	case model.LITTLE_ENDIAN:
		bo = binary.LittleEndian
	default:
		bo = binary.NativeEndian
	}

	size, ok := SIZES[typ]

	if !ok {
		switch typ {
		case "string":
			size = len(t.test.Compare.StringValue)
		case "search":
			size = len(t.test.Compare.StringValue)
		}
	}
	buff := make([]byte, size)
	_, err := t.target.Read(buff)
	if err != nil {
		return false, false, err
	}

	switch typ { // All the special cases
	case "name":
		// Like a LABEL
		return false, true, nil
	case "use":
		// It's just a GOTO
		return false, true, nil
	case "search":
		return false, true, nil
	}

	cmp, ok := TYPE_COMPARATORS[typ]
	if !ok {
		return false, false, fmt.Errorf("unknown type: %v", typ)
	}
	return cmp.compare(buff, bo, signed, t.test.Type, t.test.Compare)
}

func operation[V int64 | uint64](current, expected V, operation byte) bool {
	switch operation {
	case '=':
		return current == expected
	}
	return false
}

func operationString(current, expected string, operation byte) bool {
	switch operation {
	case '=':
		return current == expected
	}
	return false
}
