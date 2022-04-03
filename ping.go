package main

import "github.com/go-joe/joe"

func Pong(msg joe.Message) error {
	msg.Respond("PONG")
	return nil
}
