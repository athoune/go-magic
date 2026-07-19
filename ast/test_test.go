package ast

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/athoune/go-magic/parse"
	"github.com/stretchr/testify/assert"
)

func TestKitty(t *testing.T) {
	f_test, err := os.Open("../fixtures/kitty.jpg")
	assert.NoError(t, err)
	tt := model.NewTest()
	err = parse.ParseLine(tt, `0	belong&0xffffff00	0xffd8ff00	JPEG image data`)
	assert.NoError(t, err)
	output := &strings.Builder{}
	test := NewTestResult(tt, nil, output)
	ok, err := test.Test(f_test)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestKittyWithTests(t *testing.T) {
	f := model.NewFile()
	f.Name = "jpeg"

	_, err := parse.Parse(bytes.NewBufferString(
		`
0	name		jpeg
>6	string		JFIF		\b, JFIF standard
>>11	byte		x		\b %d.
>>12	byte		x		\b%02d
`), f)
	assert.NoError(t, err)
	assert.Len(t, f.Tests, 1)

	test := f.Tests[0]
	output := &strings.Builder{}
	tr := NewTestResult(test, nil, output)

	f_test, err := os.Open("../fixtures/kitty.jpg")
	assert.NoError(t, err)
	ok, err := tr.Test(f_test)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.False(t, tr.test.Message.IsTemplate)
	for i := range 2 {
		sub := tr.test.SubTests[0].SubTests[i]
		assert.Equal(t, uint64(0), sub.Type.FilterBinaryArgument)
		assert.True(t, sub.Message.IsTemplate, sub.Raw)
	}
	assert.Equal(t, ", JFIF standard 1.01", output.String())
}

func TestDynamicOffset(t *testing.T) {
	// Create a binary fixture: byte at offset 0 = 0x04, byte at offset 7 = 0x42
	// Magic: match byte at 0, then use &(0.b+3) to compute offset 4+3=7
	fixture := []byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42}

	f := model.NewFile()
	f.Name = "test"

	_, err := parse.Parse(bytes.NewBufferString(`
0	name	test
>0	byte	x	\b, found
>>&(0.b+3)	byte	x	\b, dynamic=%d
`), f)
	assert.NoError(t, err)
	assert.Len(t, f.Tests, 1)

	test := f.Tests[0]
	output := &strings.Builder{}
	tr := NewTestResult(test, nil, output)

	f_test := bytes.NewReader(fixture)
	ok, err := tr.Test(f_test)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, ", found, dynamic=66", output.String())
}
