package ast

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/athoune/go-magic/parse"
	"github.com/stretchr/testify/assert"
)

func TestTests(t *testing.T) {
	f_rules, err := os.Open("../file/magic/Magdir/jpeg")
	assert.NoError(t, err)
	file := model.NewFile()
	file.Name = "images"
	_, err = parse.Parse(f_rules, file)
	assert.NoError(t, err)
	f_rules.Close()
	f_test, err := os.Open("../fixtures/kitty.jpg")
	assert.NoError(t, err)
	for _, tt := range file.Tests {
		output := bytes.NewBufferString("")
		test := NewTest(tt)
		ok, err := test.Test(f_test, output)
		assert.NoError(t, err)
		msg := output.String()
		if msg != "" {
			fmt.Println("message:", msg)
		}
		if ok {
			fmt.Println("mime:", test.Mime, "ext:", test.Ext,
				"apple:", test.Apple, "strength:", test.Strength)
		}
	}
	//assert.False(t, true)
}
