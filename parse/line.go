package parse

import (
	"fmt"

	"github.com/athoune/go-magic/model"
)

type step_test int

const (
	STEP_LEVEL = step_test(iota)
	STEP_OFFSET
	STEP_COMPARE
	STEP_MESSAGE
)

// ParseLine parse the complete line
func ParseLine(test *model.Test, line string) error {
	test.Raw = line

	// offset
	poz := 0
	end := notSpace(line)
	err := ParseOffset(test.Offset, line[:end])
	if err != nil {
		return err
	}
	poz += end

	// type
	poz += space(line[poz:])
	end = notSpace(line[poz:])
	test.Type, err = ParseType(line[poz : poz+end])
	if err != nil {
		return err
	}
	poz += end

	//compare
	poz += space(line[poz:])
	var size int
	test.Compare, size, err = ParseCompare(line[poz:], test.Type)
	if err != nil {
		return fmt.Errorf("error in line [%v]: %v", line, err)
	}
	poz += size

	//message
	poz += space(line[poz:])
	test.Message = &model.Message{
		Value: line[poz:],
	}
	model.SetTemplateBooleans(test.Message)

	return nil
}
