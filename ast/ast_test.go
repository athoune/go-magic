package ast

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/athoune/go-magic/model"
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
	assert.NotNil(t, file.Names)
	assert.True(t, len(file.Names) > 0)
	f_test, err := os.Open("../fixtures/kitty.jpg")
	assert.NoError(t, err)
	jfif := false
	for _, tt := range file.Tests {
		output := &strings.Builder{}
		test := NewTestResult(tt, file.Names)
		ok, err := test.Test(f_test, output)
		assert.NoError(t, err, file.Names)
		msg := output.String()
		jfif = jfif || strings.Contains(msg, "JFIF")
		if ok {
			fmt.Println("mime:", test.Mime, "ext:", test.Ext,
				"apple:", test.Apple, "strength:", test.Strength)
		}
	}
	assert.True(t, jfif)
	//assert.False(t, true)
}
