package ast

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestFilterInt(t *testing.T) {
	v, err := filterInt(300, &model.Type{
		Name:     "belong",
		Operator: '&',
		Arg:      `0xffffff00`,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(256), v)
}
