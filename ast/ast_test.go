package ast

import (
	"fmt"
	"os"
	"testing"

	"github.com/athoune/go-magic/parse"
	"github.com/stretchr/testify/assert"
)

func TestTests(t *testing.T) {
	f_rules, err := os.Open("../file/magic/Magdir/jpeg")
	assert.NoError(t, err)
	testsRaw, _, err := parse.Parse(f_rules, "images")
	assert.NoError(t, err)
	f_rules.Close()
	f_test, err := os.Open("../fixtures/kitty.jpg")
	assert.NoError(t, err)
	for _, tt := range testsRaw {
		test := NewTest(tt)
		msg, err := test.Test(f_test)
		assert.NoError(t, err)
		if msg != "" {
			fmt.Println("message:", msg)
		}
	}
	//assert.False(t, true)
}
