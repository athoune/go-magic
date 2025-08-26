package model

import "strings"

type Message struct {
	Value         string
	IsTemplate    bool
	IsDisplayable bool
}

func SetTemplateBooleans(msg *Message) {
	msg.IsTemplate = strings.ContainsRune(msg.Value, '&')
	if strings.HasPrefix(msg.Value, "\b") {
		msg.IsDisplayable = true
		msg.Value = msg.Value[2:]
	}
}
