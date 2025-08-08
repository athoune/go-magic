package ast

import (
	"fmt"
	"unicode"
)

type step_test int

const (
	STEP_LEVEL = step_test(iota)
	STEP_OFFSET
	STEP_COMPARE
	STEP_MESSAGE
)

type TestLineParser struct {
	test *Test
	step step_test
	poz  int
}

func spaces(line string) int {
	poz := 0
	for i := range line {
		if !unicode.IsSpace(rune(line[i])) {
			break
		}
		poz++
	}
	return poz
}

func notSpace(line string) int {
	// nor CR
	poz := 0
	for i := range line {
		if unicode.IsSpace(rune(line[i])) {
			break
		}
		poz++
	}
	return poz
}

// ParseLine parse the complete line
func ParseLine(test *Test, line string) error {
	poz := 0
	end := notSpace(line)
	err := ParseOffset(test.Offset, line[:end])
	if err != nil {
		return err
	}
	poz += end

	poz += spaces(line[poz:])
	end = notSpace(line[poz:])

	test.Type, err = ParseType(line[poz : poz+end])
	if err != nil {
		return err
	}
	poz += end
	poz += spaces(line[poz:])
	end = notSpace(line[poz:])

	var size int
	test.Compare, size, err = ParseCompare(line[poz:], test.Type.Clue_)
	if err != nil {
		return fmt.Errorf("error in line [%v]: %v", line, err)
	}
	poz += size
	poz += spaces(line[poz:])
	test.Message = line[poz:]
	return nil
}

func Contains(needle byte, haystack string) bool {
	for h := range haystack {
		if h == int(needle) {
			return true
		}
	}
	return false
}
