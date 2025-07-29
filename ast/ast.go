package ast

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var offset_re *regexp.Regexp
var spaces_re *regexp.Regexp
var dynamic_value_re *regexp.Regexp
var test_re *regexp.Regexp

var level_idx int
var offset_idx int
var type_idx int
var compare_idx int
var data_idx int

func init() {
	offset_re = regexp.MustCompile(`(>*)(.+)`)
	dynamic_value_re = regexp.MustCompile(`(\d+).(\w)([\+\-*/%&|^])(.*)`)
	spaces_re = regexp.MustCompile(`\s+`)
	// Use https://regex101.com/ for debugging
	test_re = regexp.MustCompile(`^(?<level>>*)(?<offset>.+?)\s+(?<type>\w+)\s*(?<compare>[!=><&\^~]*[^\t]+)\t*(?<data>.*)`)

	level_idx = test_re.SubexpIndex("level")
	offset_idx = test_re.SubexpIndex("offset")
	type_idx = test_re.SubexpIndex("type")
	compare_idx = test_re.SubexpIndex("compare")
	data_idx = test_re.SubexpIndex("data")
}

func Parse(r io.Reader) ([]*Test, int, error) {
	scanner := bufio.NewScanner(r)
	var slugs []string
	var err error
	var previous *Test
	tests := make([]*Test, 0)
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
		test := &Test{
			SubTests: make([]*Test, 0),
			Actions:  make([]*Action, 0),
			Offset:   &Offset{},
		}
		if previous != nil && strings.HasPrefix(line, "!:") {
			slugs = spaces_re.Split(line[2:], -1)
			previous.Actions = append(previous.Actions, &Action{
				Name: slugs[0],
				Arg:  slugs[1],
			})
			continue
		}
		slugs = test_re.FindStringSubmatch(line)
		if len(slugs) == 0 {
			return nil, n_line, fmt.Errorf("can't parse this line : %s", line)
		}
		fmt.Println("line:", line, "level:", slugs[level_idx], "offset:",
			slugs[offset_idx], "type:", slugs[type_idx], "compare:", slugs[compare_idx],
			"data:", slugs[data_idx])
		t, ok := Types[slugs[type_idx]]
		if !ok {
			return nil, n_line, fmt.Errorf("unknown type : %s", slugs[type_idx])
		}
		test.Type = t
		test.Offset, err = ParseOffset(slugs[level_idx], slugs[offset_idx])
		if err != nil {
			return nil, n_line, err
		}
		test.Compare, err = ParseCompare(slugs[compare_idx], t.Clue_)
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
