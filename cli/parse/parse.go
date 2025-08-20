package main

import (
	"os"

	"github.com/athoune/go-magic/parse"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	wd, _ := os.Getwd()
	path := os.Args[1]
	logger.Info("",
		zap.String("Current directory", wd),
		zap.String("Path", path))
	//path := "../../file/magic/Magdir"
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		f, err := os.Open(path + "/" + e.Name())
		if err != nil {
			panic(err)
		}
		tests, _, err := parse.Parse(f)
		if err != nil {
			panic(err)
		}
		logger.Info("Magic",
			zap.String("File", e.Name()),
			zap.Any("Tests", tests))
	}
}
