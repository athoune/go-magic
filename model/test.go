package model

type Test struct {
	Offset   *Offset
	Type     *Type
	Compare  *Compare
	Message  *Message
	SubTests []*Test
	Actions  []*Action
	File     string
	Line     int
	Raw      string // the unparsed line
}

func NewTest() *Test {
	return &Test{
		Offset:   &Offset{},
		Type:     &Type{},
		Compare:  &Compare{},
		SubTests: make([]*Test, 0),
		Actions:  make([]*Action, 0),
	}
}

type Tests []*Test
