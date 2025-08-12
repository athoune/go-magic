package model

type Test struct {
	Offset   *Offset
	Type     *Type
	Compare  *Compare
	Message  string
	SubTests []*Test
	Actions  []*Action
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
