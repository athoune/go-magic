package parse

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/athoune/go-magic/model"
)

var OPERATIONS = []byte("=><&^~")

var spaces_re *regexp.Regexp
var dynamic_value_re *regexp.Regexp

var value_dynamic_idx int
var type_dynamic_idx int
var operator_dynamic_idx int
var arg_dynamic_idx int

func init() {
	spaces_re = regexp.MustCompile(`\s+`)

	dynamic_value_re = regexp.MustCompile(`(?<value>(0x)?[0-9a-f]+)((?<separator>[.,])(?<type>[bBcCeEfFgGhHiIlLmsSqQ]))?((?<operator>[+\-*])(?<arg>.*))?`)
	value_dynamic_idx = dynamic_value_re.SubexpIndex("value")
	type_dynamic_idx = dynamic_value_re.SubexpIndex("type")
	operator_dynamic_idx = dynamic_value_re.SubexpIndex("operator")
	arg_dynamic_idx = dynamic_value_re.SubexpIndex("arg")
}

func ParseFolder(path string) (model.Files, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make(model.Files, 0)
	for _, e := range entries {
		f, err := os.Open(path + "/" + e.Name())
		if err != nil {
			return nil, err
		}
		file := model.NewFile()
		file.Name = e.Name()
		_, err = Parse(f, file)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

/*
Parse reads r and compile tests described. Parsed tests are stored in 'file' argument.
*/
func Parse(r io.Reader, file *model.File) (int, error) {
	scanner := bufio.NewScanner(r)
	var slugs []string
	var err error
	testsParsing := NewTestsParsing()
	n_line := 0
	for scanner.Scan() {
		n_line += 1
		line := scanner.Text()
		if err = scanner.Err(); err != nil {
			return n_line, err
		}
		if len(line) == 0 { // empty
			continue
		}
		if line[0] == '#' { // comment
			continue
		}
		if strings.HasPrefix(line, "!:") {
			slugs = spaces_re.Split(line[2:], -1)
			testsParsing.AppendAction(&model.Action{
				Name: slugs[0],
				Arg:  slugs[1],
			})
			continue
		}
		test := model.NewTest()
		test.Line = n_line - 1
		test.File = file.Name
		err = ParseLine(test, line)
		if err != nil {
			return n_line, err
		}
		if test.Type.Name == "name" {
			file.Names[test.Compare.RawExpected] = test
		}
		testsParsing.AppendTest(test)
	}
	file.Tests = testsParsing.Tests
	return n_line, nil
}
