package model

type File struct {
	Name   string
	Header string
	Tests  Tests
	Names  map[string]*Test
}

func NewFile() *File {
	return &File{
		Names: make(map[string]*Test),
	}
}

type Files []*File
