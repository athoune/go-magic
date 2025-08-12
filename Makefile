parse:
	go run cli/parse/parse.go file/magic/Magdir

test:
	go test -v -cover ./parse
