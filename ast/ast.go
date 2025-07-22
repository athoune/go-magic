package ast

import (
	"bufio"
	"errors"
	"io"
	"log"
	"regexp"
	"strconv"
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

const (
	COMPARE_LESS    = uint8(60)
	COMPARE_GREATER = uint8(62)
	COMPARE_EQUAL   = uint8(61)
	COMPARE_AND     = uint8(38)  // &
	COMPARE_OR      = uint8(94)  // ^
	COMPARE_NEGATED = uint8(126) // ~
	COMPARE_NOT     = uint8(33)  // !
)

func Parse(r io.Reader) ([]*Test, error) {
	scanner := bufio.NewScanner(r)
	var slugs []string
	offset := Offset{}
	var err error
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

		} else {
			slugs = spaces_re.Split(line, -1)
			ooo := offset_re.FindStringSubmatch(slugs[0]) // looking for > prefix
			if len(ooo) == 2 {                            // no
				offset.Level = 0
				offset.Value, err = strconv.ParseInt(ooo[1], 0, 32)
				if err != nil {
					return nil, err
				}
			} else { // additional tests
				offset.Level = len(ooo[1])
				offset.Value, err = strconv.ParseInt(ooo[2], 0, 32)
				if err != nil {
					return nil, err
				}
			}
			tests = append(tests, &Test{
				Offset: offset,
				Type:   slugs[1],
			})
		}
	}
	return tests, nil
}

func ParseOffset(txt string) (*Offset, error) {
	offset := &Offset{}
	var err error
	var value string
	ooo := offset_re.FindStringSubmatch(txt) // looking for > prefix
	if len(ooo) == 2 {                       // no
		offset.Level = 0
		value = ooo[1]
	} else {
		offset.Level = len(ooo[1])
		value = ooo[2]
	}

	if strings.HasPrefix(value, "(") {
		offset.Dynamic = true
		dyn := dynamic_value_re.FindStringSubmatch(value[1 : len(value)-1])
		if len(dyn) != 5 {
			return nil, errors.New("Bad dynamic offset value : " + value)
		}
		offset.DynOffset, err = strconv.ParseInt(dyn[1], 0, 32)
		if err != nil {
			return nil, err
		}
		offset.DynType = dyn[2][0]
		offset.DynAction = dyn[3][0]
		offset.DynArg, err = strconv.ParseInt(dyn[4], 0, 32)
		if err != nil {
			return nil, err
		}
	} else {
		offset.Value, err = strconv.ParseInt(value, 0, 32)
		if err != nil {
			return nil, err
		}
	}
	return offset, nil
}

func ParseTest(txt string, type_ byte) (*Compare, error) {
	if txt[0] == 'x' {
		return nil, nil
	}
	compare := &Compare{}
	var err error
	i := 0
	if txt[i] == '!' {
		compare.Not = true
		i++
	}
	compare.Operation = txt[i]
	not_implicit_equality := false
	for _, a := range []byte("=><&^~") {
		if compare.Operation == a {
			not_implicit_equality = true
			break
		}
	}
	if !not_implicit_equality {
		compare.Operation = '='
	} else {
		i++
	}
	compare.Type = type_
	switch {
	case type_ == 's':
		compare.StringValue = txt[i:]
	case type_ == 'f':
		compare.FloatValue, err = strconv.ParseFloat(txt[i:], 64)
		if err != nil {
			return nil, err
		}
	case type_ == 'i':
		compare.IntValue, err = strconv.ParseInt(txt[i:], 0, 64)
		if err != nil {
			return nil, err
		}
	}
	return compare, nil
}
