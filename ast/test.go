package ast

type Test struct {
	Offset   *Offset
	Type     *Type
	Compare  *Compare
	Message  string
	SubTests []*Test
	Actions  []*Action
}
