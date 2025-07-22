package ast

type Test struct {
	Offset   *Offset
	Type     string
	Compare  *Compare
	SubTests []*Test
}
