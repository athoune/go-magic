package parse

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestParseType(t *testing.T) {
	type_, err := ParseType(`belong&0xfe00f0f0`)
	assert.NoError(t, err)
	assert.Equal(t, "belong", type_.Name)
	assert.Equal(t, byte('&'), type_.FilterOperator)
	assert.Equal(t, uint64(0xfe00f0f0), type_.FilterBinaryArgument)

	type_, err = ParseType(`pstring/HJ`)
	assert.NoError(t, err)
	assert.Equal(t, "pstring", type_.Root)
	assert.Equal(t, "pstring", type_.Name)
	assert.Equal(t, byte('/'), type_.FilterOperator)
	assert.Equal(t, "HJ", type_.FilterStringArgument)
}

func TestParseStringOptions(t *testing.T) {
	typ := &model.Type{}
	parseOptions(typ, "bob/cC")
	assert.Equal(t, model.STRING_OPTIONS_NONE, typ.StringOptions)
	stringOptions, err := parseStringOptions(typ.FilterStringArgument)
	assert.NoError(t, err)
	assert.Equal(t, model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER,
		stringOptions&model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER)
	assert.Equal(t, model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER,
		stringOptions&model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER)
}

func TestParseSearchOptions(t *testing.T) {
	for _, fixture := range []struct {
		line          string
		searchRange   int
		stringOptions model.StringOptions
	}{
		{"search/1/t", 1, model.STRING_OPTIONS_TEXT_FILE},
		{"search/42/cC", 42,
			model.STRING_OPTIONS_CASE_INSENSITIVE_LOWER |
				model.STRING_OPTIONS_CASE_INSENSITIVE_UPPER},
		{"search/727", 727, model.STRING_OPTIONS_NONE},
		{"search/210965/s", 210965, model.REGEX_OPTIONS_OFFSET_START},
		{"search/8192", 8192, model.STRING_OPTIONS_NONE},
	} {
		typ := &model.Type{}
		parseOptions(typ, fixture.line)
		searchRange, stringOptions, err := parseSearchOptions(typ.FilterStringArgument)
		assert.NoError(t, err)
		assert.Equal(t, fixture.searchRange, searchRange)
		assert.Equal(t, fixture.stringOptions, stringOptions)
	}
}
