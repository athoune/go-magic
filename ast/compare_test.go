package compare

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	ok, err := CompareValues(&model.Value{
		Family:   model.TYPE_CLUE_INT,
		IntValue: 12,
	}, &model.Value{
		Family:   model.TYPE_CLUE_INT,
		IntValue: 42,
	}, '<')
	assert.NoError(t, err)
	assert.True(t, ok)
}
