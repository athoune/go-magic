.PHONY: parse

parse:
	go run cli/parse/parse.go file/magic/Magdir

test:
	go test -v -cover ./ast
	go test -v -cover ./model
	go test -v -cover ./parse
	go test -v -cover ./unpack

govulncheck:
	govulncheck ast/*.go
	govulncheck parse/*.go
	govulncheck model/*.go
	govulncheck unpack/*.go
