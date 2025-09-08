package model

import "strings"

type Message struct {
	Value      string
	IsTemplate bool
}

func SetTemplateBooleans(msg *Message) {
	msg.IsTemplate = strings.ContainsRune(msg.Value, '%')
	msg.Value = strings.TrimPrefix(msg.Value, `\b`)
}
