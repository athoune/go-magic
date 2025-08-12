package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleStringEscape(t *testing.T) {
	for _, fixture := range []struct {
		raw     string
		escaped string
	}{
		{`\x02\x01\x13\x13\x10\x14\x12\x0e`, "\x02\x01\x13\x13\x10\x14\x12\x0e"},
		{`plop`, "plop"},
		{`\x8aMNG`, "\x8aMNG"},
		{`Beuha\ aussi`, "Beuha aussi"},
	} {
		s, err := HandleStringEscape(fixture.raw)
		assert.NoError(t, err)
		assert.Equal(t, fixture.escaped, s)
	}
}

func TestHandleSpaceEscape(t *testing.T) {
	fixture := `Invalid\ partition\ table  `
	end := HandleSpaceEscape(fixture)
	assert.Equal(t, `Invalid\ partition\ table`, fixture[:end])
	fixture = `Invalid\ partition\ table`
	end = HandleSpaceEscape(fixture)
	assert.Equal(t, `Invalid\ partition\ table`, fixture[:end])
}
