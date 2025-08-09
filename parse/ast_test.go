package parse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	fixture := `# Standard PNG image.
0	string		\x89PNG\x0d\x0a\x1a\x0a\x00\x00\x00\x0DIHDR	PNG image data
!:mime	image/png
!:ext   png
!:strength +10
>16	use		png-ihdr
>33	string		\x00\x00\x00\x08acTL	\b, animated
>>41	ubelong		1			(%d frame
>>41	ubelong		>1			(%d frames
>>45	ubelong		0			\b, infinite repetitions)
>>45	ubelong		1			\b, %d repetition)
>>45	ubelong		>1			\b, %d repetitions)`

	tests, _, err := Parse(strings.NewReader(fixture))
	assert.NoError(t, err)
	assert.Len(t, tests, 1)
	test := tests[0]
	assert.Len(t, test.Actions, 3)
	assert.Len(t, test.SubTests, 2)

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
