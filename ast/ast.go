package ast

import (
	"bytes"
	"fmt"
	"io"

	"github.com/athoune/go-magic/model"
)

type Runner struct {
	Files model.Files
}

func NewRunner(files model.Files) *Runner {
	return &Runner{
		Files: files,
	}
}

func (r Runner) Magic(target io.ReadSeeker) (string, error) {
	output := bytes.NewBufferString("")
	for _, file := range r.Files {
		fmt.Println(file.Names)
		for _, t := range file.Tests {
			test := NewTestResult(t, file.Names)
			_, err := test.Test(target, output)
			if err != nil {
				return "", err
			}
		}
	}
	return output.String(), nil
}
