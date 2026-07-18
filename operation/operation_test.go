package operation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareString(t *testing.T) {
	for i, fixture := range []struct {
		file    string
		pattern string
		match   bool
	}{
		{"JFIF\x00", "JFIF\x00", true},        // exact match
		{"JFIF\x00extra", "JFIF\x00", true},   // file longer than pattern
		{"JFI", "JFIF\x00", false},             // file shorter, no zero padding match
		{"JFIF", "JFIF\x00", true},             // file shorter by 1, zero-padded match
		{"hello", "hello", true},                // no null terminator
		{"", "\x00", true},                      // empty file vs single null
	} {
		v, err := CompareString(fixture.file, fixture.pattern, '=')
		assert.NoError(t, err)
		assert.Equal(t, fixture.match, v, "case %d: file=%q pattern=%q", i, fixture.file, fixture.pattern)
	}
}

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
