.PHONY: parse

parse:
	go run cli/parse/parse.go file/magic/Magdir

test:
	go test -v -cover ./parse
	go test -v -cover ./ast

govulncheck:
	govulncheck parse/*.go
	govulncheck model/*.go
