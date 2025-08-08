package main

import (
	"fmt"
	"os"

	"github.com/athoune/go-magic/ast"
)

func main() {
	fmt.Println(os.Getwd())
	path := os.Args[1]
	//path := "../../file/magic/Magdir"
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		fmt.Println("file :", e.Name(), "\n\n")
		f, err := os.Open(path + "/" + e.Name())
		if err != nil {
			panic(err)
		}
		tests, _, err := ast.Parse(f)
		if err != nil {
			panic(err)
		}
		for _, test := range tests {
			fmt.Println(test)
		}
	}

}
