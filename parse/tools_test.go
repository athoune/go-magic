package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleSpaceEscape(t *testing.T) {
	fixture := `Invalid\ partition\ table  `
	end := HandleSpaceEscape(fixture)
	assert.Equal(t, `Invalid\ partition\ table`, fixture[:end])
	fixture = `Invalid\ partition\ table`
	end = HandleSpaceEscape(fixture)
	assert.Equal(t, `Invalid\ partition\ table`, fixture[:end])
}

func TestSpaces(t *testing.T) {
	i := space("  	plop")
	assert.Equal(t, 3, i)
	i = space("	 ")
	assert.Equal(t, 2, i)
	i = space("")
	assert.Equal(t, 0, i)
}

func TestNotSpaces(t *testing.T) {
	i := notSpace("beuha ")
	assert.Equal(t, 5, i)
	i = notSpace("beuha")
	assert.Equal(t, 5, i)
}
