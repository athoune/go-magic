package main

import (
	"os"

	"github.com/athoune/go-magic/parse"
	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	logger := zap.L()
	wd, _ := os.Getwd()
	path := os.Args[1]
	logger.Info("",
		zap.String("Current directory", wd),
		zap.String("Path", path))
	//path := "../../file/magic/Magdir"
	_, err := parse.ParseFolder(path)
	if err != nil {
		panic(err)
	}
}
