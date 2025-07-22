package ast

type File struct {
	Name   string
	Header string
	Tests  []Test
}
