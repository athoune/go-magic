package ast

import (
	"fmt"
	"io"

	"github.com/athoune/go-magic/model"
)

type Runner struct {
	Files *model.Files
}

func NewRunner(files *model.Files) *Runner {
	return &Runner{
		Files: files,
	}
}

func (r Runner) Magic(target io.ReadSeeker) (string, error) {
	for _, file := range *r.Files {
		for _, t := range file.Tests {
			test := NewTest(t)
			infos, err := test.Test(target)
			if err != nil {
				return "", err
			}
			fmt.Println(infos)
		}
	}
	return "", nil
}
