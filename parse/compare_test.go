package parse

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestParseCompare(t *testing.T) {
	for _, fixture := range []struct {
		line    string
		type_   *model.Type
		size    int
		compare *model.Compare
	}{
		{"<10", &model.Type{
			TypeFamily: model.TYPE_FAMILY_INT,
			Root:       "long",
		}, 3, &model.Compare{
			Comparator:  COMPARE_LESS,
			RawExpected: "10",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 10,
			},
		}},
		{"< 10", &model.Type{
			TypeFamily: model.TYPE_FAMILY_INT,
			Root:       "long",
		}, 4, &model.Compare{
			Comparator:  COMPARE_LESS,
			RawExpected: "10",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 10,
			},
		}},
		{"0x01000007", &model.Type{
			TypeFamily: model.TYPE_FAMILY_INT,
			Root:       "long",
		}, 10, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: "0x01000007",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 16777223,
			},
		}},
		{"!>10", &model.Type{
			TypeFamily: model.TYPE_FAMILY_INT,
			Root:       "long",
		}, 4, &model.Compare{
			Comparator:  COMPARE_GREATER,
			RawExpected: "10",
			Expected: &model.Value{
				Family:   model.TYPE_FAMILY_INT,
				IntValue: 10,
			},
			Not: true,
		}},
		{"D6E229D3-35DA-11D1-9034-00A0C90349BE", &model.Type{
			TypeFamily: model.TYPE_FAMILY_STRING}, 36, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: "D6E229D3-35DA-11D1-9034-00A0C90349BE",
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "D6E229D3-35DA-11D1-9034-00A0C90349BE",
			},
		}},
		{`Invalid\ partition\ table		english`, &model.Type{
			TypeFamily: model.TYPE_FAMILY_STRING}, 25, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: `Invalid\ partition\ table`,
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "Invalid partition table",
			},
		}},
		{`\x6d\x6a\x70\x32`, &model.Type{
			TypeFamily: model.TYPE_FAMILY_STRING,
			Root:       "string",
		}, 16, &model.Compare{
			Comparator:  COMPARE_EQUAL,
			RawExpected: `\x6d\x6a\x70\x32`,
			Expected: &model.Value{
				Family:      model.TYPE_FAMILY_STRING,
				StringValue: "mjp2",
			},
		}},
	} {
		c, s, err := ParseCompare(fixture.line, fixture.type_)
		assert.NoError(t, err)
		// [FIXME]
		fixture.compare.Type = fixture.type_ // yes, it's cheating
		assert.Equal(t, fixture.compare, c, fixture.line)
		assert.Equal(t, fixture.size, s, fixture.line)
	}
}

func TestParseCompareName(t *testing.T) {
	c, _, err := ParseCompare(`jpeg`, &model.Type{
		Name: "name",
	})
	assert.NoError(t, err)
	assert.Equal(t, "jpeg", c.RawExpected)
}
