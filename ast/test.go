package ast

import (
	"encoding/binary"
	"io"

	"github.com/athoune/go-magic/model"
)

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
	t.target = target
	t.offset()
	return "", nil
}

func (t *Test) offset() {
	if t.test.Offset.Relative {
		// TODO
	} else {
		t.target.Seek(t.test.Offset.Value, io.SeekStart)
	}

}

func compare(cmp *model.Compare, target io.ReadSeeker) error {
	switch cmp.Type.Name {
	case "name":
		// Like a LABEL
	case "use":
		// It's just a GOTO
	case "ubelong":
		buff := make([]byte, 4)
		_, err := target.Read(buff)
		if err != nil {
			return err
		}
		binary.BigEndian.Uint32(buff)

	}
	return nil
}
