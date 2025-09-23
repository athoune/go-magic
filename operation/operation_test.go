package operation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperate(t *testing.T) {
	for i, fixture := range []struct {
		operator byte
		result   int
	}{
		{'+', 3},
		{'-', -1},
		{'*', 2},
		{'/', 0},
		{'&', 0},
		{'|', 3},
		{'%', 1},
		{' ', 1},
	} {
		v, err := Operate(1, 2, fixture.operator)
		assert.NoError(t, err)
		assert.Equal(t, fixture.result, v, i)
		uv, err := Operate(uint64(1), uint64(2), fixture.operator)
		assert.NoError(t, err)
		if fixture.operator == '-' {
			assert.Equal(t, uint64(1<<64-1), uv, i)
		} else {
			assert.Equal(t, uint64(fixture.result), uv, i)
		}
	}
}
