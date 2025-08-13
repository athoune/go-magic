package model

type File struct {
	Name   string
	Header string
	Tests  Tests
}

type Files []*File
