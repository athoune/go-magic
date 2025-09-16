package parse

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestParseOffset(t *testing.T) {
	o := &model.Offset{}
	err := ParseOffset(o, ">>>4")
	assert.NoError(t, err)
	assert.Equal(t, 3, o.Level)
	assert.Equal(t, int64(4), o.Value)
	assert.False(t, o.Dynamic)

	o = &model.Offset{}
	err = ParseOffset(o, ">>(4.s*512)")
	assert.NoError(t, err)
	assert.Equal(t, 2, o.Level)
	assert.True(t, o.Dynamic)
	assert.Equal(t, int64(4), o.DynOffset)
	assert.Equal(t, uint8('s'), o.DynType)
	assert.Equal(t, uint8('*'), o.DynOperator)
	assert.Equal(t, int64(512), o.DynArg)

	o = &model.Offset{}
	err = ParseOffset(o, ">>>>>>(&4.l+(-4))")
	assert.NoError(t, err)
	assert.Equal(t, 6, o.Level)
	assert.True(t, o.Dynamic)
	assert.Equal(t, int64(4), o.DynOffset)
	assert.Equal(t, uint8('l'), o.DynType)
	assert.Equal(t, uint8('+'), o.DynOperator)
	assert.Equal(t, int64(-4), o.DynArg)
}

func TestParseDynamicOffset(t *testing.T) {
	o := &model.Offset{}
	err := ParseDynamicOffset(o, "4.s*512")
	assert.NoError(t, err)
	assert.True(t, o.Dynamic)
	assert.Equal(t, int64(4), o.DynOffset)
	assert.Equal(t, uint8('*'), o.DynOperator)
	assert.Equal(t, int64(512), o.DynArg)
}
