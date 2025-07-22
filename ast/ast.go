package ast

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

var offset_re *regexp.Regexp
var spaces_re *regexp.Regexp
var dynamic_value_re *regexp.Regexp

func init() {
	offset_re = regexp.MustCompile(`(>*)(.+)`)
	dynamic_value_re = regexp.MustCompile(`(\d+).(\w)([\+\-*/%&|^])(.*)`)
	spaces_re = regexp.MustCompile(("\t+"))
}

func Parse(r io.Reader) ([]*Test, error) {
	scanner := bufio.NewScanner(r)
	var slugs []string
	tests := make([]*Test, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		if len(line) == 0 { // empty
			continue
		}
		if line[0] == '#' { // comment
			continue
		}
		if strings.HasPrefix(line, "!:") {

			continue
		}
		slugs = spaces_re.Split(line, -1)
		offset, err := ParseOffset(slugs[0])
		if err != nil {
			return nil, err
		}
		compare, err := ParseCompare(slugs[1], Types[slugs[1]].Clue_)
		tests = append(tests, &Test{
			Offset:  offset,
			Type:    slugs[1],
			Compare: compare,
		})
	}
	return tests, nil
}
