package parse

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/athoune/go-magic/model"
)

var spaces_re *regexp.Regexp
var dynamic_value_re *regexp.Regexp

var value_dynamic_idx int
var type_dynamic_idx int
var operator_dynamic_idx int
var arg_dynamic_idx int

func init() {
	dynamic_value_re = regexp.MustCompile(`(?<value>(0x)?[0-9a-f]+)((?<separator>[.,])(?<type>[bBcCeEfFgGhHiIlLmsSqQ]))?((?<operator>[+\-*])(?<arg>.*))?`)
	spaces_re = regexp.MustCompile(`\s+`)

	value_dynamic_idx = dynamic_value_re.SubexpIndex("value")
	type_dynamic_idx = dynamic_value_re.SubexpIndex("type")
	operator_dynamic_idx = dynamic_value_re.SubexpIndex("operator")
	arg_dynamic_idx = dynamic_value_re.SubexpIndex("arg")

}

func Parse(r io.Reader) ([]*model.Test, int, error) {
	scanner := bufio.NewScanner(r)
	var slugs []string
	var err error
	var previous *model.Test
	tests := make([]*model.Test, 0)
	n_line := 0
	for scanner.Scan() {
		n_line += 1
		line := scanner.Text()
		if err = scanner.Err(); err != nil {
			return nil, n_line, err
		}
		if len(line) == 0 { // empty
			continue
		}
		if line[0] == '#' { // comment
			continue
		}
		test := model.NewTest()
		if previous != nil && strings.HasPrefix(line, "!:") {
			slugs = spaces_re.Split(line[2:], -1)
			previous.Actions = append(previous.Actions, &model.Action{
				Name: slugs[0],
				Arg:  slugs[1],
			})
			continue
		}
		err = ParseLine(test, line)
		if err != nil {
			return nil, n_line, err
		}
		if previous != nil && test.Offset.Level >= previous.Offset.Level {

		} else {
			tests = append(tests, test)
		}
		previous = test // FIXME
	}
	return tests, n_line, nil
}
