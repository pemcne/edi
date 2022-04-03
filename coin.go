package main

import (
	"fmt"
	"math/rand"

	"github.com/go-joe/joe"
)

func CoinFlip(msg joe.Message) error {
	sides := []string{"heads", "tails"}
	choice := sides[rand.Intn(len(sides)-1)]
	msg.Respond(fmt.Sprintf("It's %s!", choice))
	return nil
}
