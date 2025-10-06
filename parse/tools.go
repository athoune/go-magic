package parse

import (
	"unicode"
)

var OPERATIONS = []byte("=><&^~")

// HandleSpaceEscape line is closed by some spaces, but it handles escape : "\ "
func HandleSpaceEscape(line string) int {
	end := 0
	for {
		ns := notSpace(line[end:])
		if ns == 0 {
			break
		}
		end += ns
		if end+1 < len(line) && line[end-1:end+1] == `\ ` {
			end++
		} else {
			break
		}
	}
	return end
}

func IsOperation(op byte) bool {
	for _, a := range OPERATIONS {
		if op == a {
			return true
		}
	}
	return false
}

func Contains(needle byte, haystack string) bool {
	for h := range haystack {
		if h == int(needle) {
			return true
		}
	}
	return false
}

// space read all spaces, return the position before the first non space character
func space(line string) int {
	poz := 0
	for i := range line {
		if !unicode.IsSpace(rune(line[i])) {
			break
		}
		poz++
	}
	return poz
}

// notSpace the position before the first non space charactert
func notSpace(line string) int {
	// nor CR
	poz := 0
	for i := range line {
		if unicode.IsSpace(rune(line[i])) {
			break
		}
		poz++
	}
	return poz
}
