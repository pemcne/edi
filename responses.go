package main

import (
	"math/rand"
	"time"

	"github.com/go-joe/joe"
)

var RandomGenerator *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func chooseRandom(choices []string) string {
	randNum := RandomGenerator.Intn(len(choices))
	return choices[randNum]
}

func CommonResponses(b *joe.Bot) {
	b.Hear(`^NO U`, func(msg joe.Message) error {
		msg.Respond("NO U")
		return nil
	})

	b.Hear(`(?i)^hodor`, func(msg joe.Message) error {
		msg.Respond("Hodor!")
		return nil
	})

	b.Hear(`(?i)^good bot`, func(msg joe.Message) error {
		choices := []string{":blush: thanks!", ":smiling_face_with_smiling_eyes_and_hand_covering_mouth: thanks!"}
		msg.Respond(chooseRandom(choices))
		return nil
	})

	b.Hear(`(?i)^bad bot`, func(msg joe.Message) error {
		choices := []string{
			"I'll strive to do better",
			":feelsbadman:",
		}
		msg.Respond(chooseRandom(choices))
		return nil
	})

	b.Respond("(?i)ping", func(msg joe.Message) error {
		msg.Respond("PONG")
		return nil
	})
}
