package main

import (
	"os"

	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter/v2"
	"go.uber.org/zap/zapcore"
)

var Edi joe.Bot

func main() {
	adapter := slack.EventsAPIAdapter(
		os.Getenv("SLACK_TOKEN"),
		slack.WithSocketMode(os.Getenv("SLACK_APP_TOKEN")),
		slack.WithListenPassive(),
		// slack.WithDebug(true),
	)
	Edi := joe.New("Edi", joe.WithLogLevel(zapcore.DebugLevel), adapter)

	Edi.Respond("ping", Pong)
	Edi.Respond("flip a coin", CoinFlip)

	err := Edi.Run()
	if err != nil {
		Edi.Logger.Fatal(err.Error())
	}
}
