package main

import (
	"fmt"
	"os"

	"github.com/athoune/go-magic/ast"
)

func main() {
	fmt.Println(os.Getwd())
	entries, err := os.ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		f, err := os.Open(os.Args[1] + "/" + e.Name())
		if err != nil {
			panic(err)
		}
		tests, err := ast.Parse(f)
		if err != nil {
			panic(err)
		}
		for _, test := range tests {
			fmt.Println(test)
		}
	}

}
