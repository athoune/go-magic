package ast

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

func ReadPstring(r io.Reader, options string) (string, error) {
	if len(options) > 2 {
		return "", fmt.Errorf("arguments parameter is too long: %v", len(options))
	}
	lengthInclude := false
	li := strings.IndexRune(options, 'J')
	mod := options[0]
	if li != -1 {
		lengthInclude = true
		if li == 0 {
			mod = options[1]
		}
	}
	var length uint64
	var lengthLength int
	var valueRaw []byte
	var bo binary.ByteOrder
	switch {
	case mod == 'B':
		lengthLength = 1
		lengthRaw := make([]byte, lengthLength)
		_, err := r.Read(lengthRaw)
		if err != nil {
			return "", err
		}
		length = uint64(lengthRaw[0])
	case mod == 'H' || mod == 'h':
		lengthLength = 2
		if mod == 'H' {
			bo = binary.BigEndian
		} else { // h
			bo = binary.LittleEndian
		}
		var l uint16
		err := binary.Read(r, bo, &l)
		if err != nil {
			return "", err
		}
		length = uint64(l)
	case mod == 'L' || mod == 'l':
		lengthLength = 4
		if mod == 'L' {
			bo = binary.BigEndian
		} else { // l
			bo = binary.LittleEndian
		}
		var l uint32
		err := binary.Read(io.LimitReader(r, int64(lengthLength)), bo, &l)
		if err != nil {
			return "", err
		}
		length = uint64(l)
	default:
		err := fmt.Errorf("unknown modifiers: %v", options)
		if err != nil {
			return "", err
		}
	}
	if lengthInclude {
		length -= uint64(lengthLength)
	}
	valueRaw = make([]byte, length)
	_, err := r.Read(valueRaw)
	return string(valueRaw), err
}
