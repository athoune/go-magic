package ast

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPString(t *testing.T) {
	for _, fixture := range []struct {
		content  []byte
		options  string
		expected string
	}{
		{[]byte{4, 'p', 'l', 'o', 'p'}, "B", "plop"},
		{[]byte{0x0, 0x1c, 'L', 'a', 'z', 'y', ' ', 'd', 'o', 'g', ' ',
			'j', 'u', 'm', 'p', 's', ' ', 'o', 'v', 'e', 'r', ' ',
			't', 'h', 'e', ' ', 'w', 'a', 'l', 'l'}, "H", "Lazy dog jumps over the wall"},
		{[]byte{4, 0, 0, 0, 'p', 'l', 'o', 'p'}, "l", "plop"},
		{[]byte{4 + 1, 'p', 'l', 'o', 'p'}, "BJ", "plop"},
	} {
		r := bytes.NewReader(fixture.content)
		v, err := ReadPstring(r, fixture.options)
		assert.NoError(t, err)
		assert.Equal(t, fixture.expected, v, fixture)
	}
}
