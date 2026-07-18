package operation

import (
	"testing"

	"github.com/athoune/go-magic/model"
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
		v, err := CompareString(fixture.file, fixture.pattern, '=', model.STRING_OPTIONS_NONE)
		assert.NoError(t, err)
		assert.Equal(t, fixture.match, v, "case %d: file=%q pattern=%q", i, fixture.file, fixture.pattern)
	}
}

func TestCompareStringFlags(t *testing.T) {
	for i, fixture := range []struct {
		file    string
		pattern string
		opts    model.StringOptions
		match   bool
	}{
		// /c ignore lowercase: file lowercase matches any case in pattern
		{"hello", "HELLO", model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER, true},
		{"HELLO", "hello", model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER, false},
		// /C ignore uppercase: file uppercase matches any case in pattern
		{"HELLO", "hello", model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER, true},
		{"hello", "HELLO", model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER, false},
		// /f full word: match only if followed by space or end of string
		{"hello", "hello", model.STRING_OPTIONS_FULL_WORD, true},
		{"helloworld", "hello", model.STRING_OPTIONS_FULL_WORD, false},
		{"hello world", "hello", model.STRING_OPTIONS_FULL_WORD, true},
		// /w optional whitespace: file spaces are skipped, pattern spaces are required
		{"hello world", "hello world", model.STRING_OPTIONS_COMPACT_OPTIONAL_WHITESPACE, true},
		{"hello  world", "hello world", model.STRING_OPTIONS_COMPACT_OPTIONAL_WHITESPACE, true},
		{"helloworld", "hello world", model.STRING_OPTIONS_COMPACT_OPTIONAL_WHITESPACE, false},
		// /W compact whitespace: file must have space when pattern does
		{"hello world", "hello world", model.STRING_OPTIONS_COMPACT_WITH_SPACES, true},
		{"helloworld", "hello world", model.STRING_OPTIONS_COMPACT_WITH_SPACES, false},
		{"hello  world", "hello world", model.STRING_OPTIONS_COMPACT_WITH_SPACES, false},
	} {
		v, err := CompareString(fixture.file, fixture.pattern, '=', fixture.opts)
		assert.NoError(t, err)
		assert.Equal(t, fixture.match, v, "case %d: opts=%d file=%q pattern=%q", i, fixture.opts, fixture.file, fixture.pattern)
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
