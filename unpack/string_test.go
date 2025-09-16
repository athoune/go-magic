package unpack

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestStringValue(t *testing.T) {
	typ := &model.Type{
		Root:   "long",
		Clue:   model.TYPE_CLUE_INT,
		Signed: true,
	}
	v, err := BuildValueFromString(typ, "42")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), v.IntValue)

	v, err = BuildValueFromString(typ, "0xF")
	assert.NoError(t, err)
	assert.Equal(t, int64(15), v.IntValue)
}
