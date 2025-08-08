package ast

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/athoune/go-magic/model"
)

func TestParse(t *testing.T) {

	tests, line, err := Parse(strings.NewReader(`>>16	belong&0xfe00f0f0	0x3030`))
	assert.Nil(t, err)
	assert.Equal(t, 1, line)
	assert.NotNil(t, tests)
	assert.Len(t, tests, 1)
	test := tests[0]
	assert.Equal(t, 2, test.Offset.Level)
	assert.Equal(t, int64(16), test.Offset.Value)
	assert.Equal(t, model.TYPE_CLUE_INT, test.Type.Clue_)
	assert.Equal(t, "belong", test.Type.Name)
	assert.Equal(t, COMPARE_EQUAL, test.Compare.Operation)
	assert.False(t, test.Compare.Not)

}
func TestRead(t *testing.T) {
	r := strings.NewReader(`
# RISC OS Chunk File Format
# From RISC OS Programmer's Reference Manual, Appendix D
# We guess the file type from the type of the first chunk.
0	lelong		0xc3cbc6c5	RISC OS Chunk data
>12	string		OBJ_		\b, AOF object
>12	string		LIB_		\b, ALF library
0	string		Draw		RISC OS Draw file data
`)
	tests, _, _ := Parse(r)
	for _, test := range tests {
		fmt.Println(test)
	}
}
