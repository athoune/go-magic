package ast

import (
	"testing"

	"github.com/athoune/go-magic/model"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	for _, fixture := range []struct {
		typ      *model.Type
		value    *model.Value
		expected int64
	}{
		{&model.Type{
			Root:                 "short",
			ByteOrder:            model.LITTLE_ENDIAN,
			FilterOperator:       '&',
			FilterBinaryArgument: 0x0f,
		},
			&model.Value{
				Family:         model.TYPE_FAMILY_INT,
				RawBinaryValue: 0xa3,
				UIntValue:      163,
			},
			int64(0x03),
		},
		{&model.Type{
			Root:                 "long",
			ByteOrder:            model.LITTLE_ENDIAN,
			FilterOperator:       '|',
			FilterBinaryArgument: 0x0001,
		},
			&model.Value{
				Family:         model.TYPE_FAMILY_INT,
				RawBinaryValue: 0x302a,
				UIntValue:      12330,
			},
			int64(0x302b),
		},
	} {
		err := FilterValue(fixture.typ, fixture.value)
		assert.NoError(t, err)
		assert.Equal(t, fixture.expected, fixture.value.IntValue, "Type: %#v Value: %#v", fixture.typ, fixture.value)
	}
}
