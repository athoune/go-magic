package parse

import (
	"bytes"
	"strconv"
	"unicode"
)

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

func HandleStringEscape(value string) (string, error) {
	poz := 0
	out := &bytes.Buffer{}
	for {
		if poz == len(value) {
			break
		}
		switch {
		case poz+4 <= len(value) && value[poz:poz+2] == `\x`:
			v, err := strconv.ParseInt(value[poz+2:poz+4], 16, 64)
			if err != nil {
				return value, nil // YOLO, file use \x2\x4 in archive#1362
			}
			out.WriteByte(byte(v))
			poz += 4
		case poz+2 <= len(value) && value[poz:poz+2] == `\ `:
			out.WriteByte(' ')
			poz += 2
		default:
			out.WriteByte(value[poz])
			poz++
		}
	}
	return out.String(), nil

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

// space find the first space (\n, \t, something like that) in a string
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

// notSpace find the first non-space character in a string
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
