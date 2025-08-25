package parse

import "github.com/athoune/go-magic/model"

type TestsParsing struct {
	Tests     model.Tests
	offsets   map[int]*model.Test
	lastLevel int
}

func NewTestsParsing() *TestsParsing {
	return &TestsParsing{
		Tests:   make(model.Tests, 0),
		offsets: make(map[int]*model.Test),
	}
}

func (t *TestsParsing) AppendTest(test *model.Test) {
	t.offsets[test.Offset.Level] = test
	t.lastLevel = test.Offset.Level
	if test.Offset.Level == 0 {
		t.Tests = append(t.Tests, test)
	} else {
		t.offsets[test.Offset.Level-1].SubTests = append(t.offsets[test.Offset.Level-1].SubTests, test)
	}
}

func (t *TestsParsing) AppendAction(action *model.Action) {
	t.offsets[t.lastLevel].Actions = append(t.offsets[t.lastLevel].Actions, action)
}
