package main

import (
	"strings"

	"github.com/go-joe/joe"
)

func Points(msg joe.Message) error {
	symbol := msg.Matches[0]
	amount := msg.Matches[1]
	key := strings.TrimSpace(msg.Matches[4])

	Edi.Store.Get()
	return nil
}
