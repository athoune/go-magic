package ast

import (
	"fmt"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	r := strings.NewReader(`
# RISC OS Chunk File Format
# From RISC OS Programmer's Reference Manual, Appendix D
# We guess the file type from the type of the first chunk.
0	lelong		0xc3cbc6c5	RISC OS Chunk data
>12	string		OBJ_		\b, AOF object
>12	string		LIB_		\b, ALF library
`)
	tests, _ := Parse(r)
	for _, test := range tests {
		fmt.Println(test)
	}
}
