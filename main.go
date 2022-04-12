package main

import (
	"os"

	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter/v2"
	"go.uber.org/zap/zapcore"
)

var Edi *joe.Bot

func main() {
	adapter := slack.EventsAPIAdapter(
		os.Getenv("SLACK_TOKEN"),
		slack.WithSocketMode(os.Getenv("SLACK_APP_TOKEN")),
		slack.WithListenPassive(),
		// slack.WithDebug(true),
	)
	Edi = joe.New("Edi", joe.WithLogLevel(zapcore.DebugLevel), file.Memory("brain.json"), adapter)

	Edi.Respond("ping", Pong)
	Edi.Respond("flip a coin", CoinFlip)
	Edi.Respond("(\\+|-)\\s*(\\d+) (to|for) (.+)", Points)
	Edi.Respond("score for (.+)", PointsScore)
	Edi.Respond("leaderboard", PointsLeaderboard)
	Edi.Respond("what happened today", Today)

	Edi.Hear("Wordle\\s\\d+\\s(.+)/\\d", WordleScore)
	Edi.Hear("Dordle\\s#\\d+\\s(.+)/\\d", DordleScore)
	Edi.Hear("Quordle\\s\\d+\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)", QuordleScore)
	Edi.Hear("Octordle\\s#\\d+\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)", OctordleScore)
	Edi.Hear("Worldle\\s#\\d+\\s(.+)/\\d", WorldleScore)

	err := Edi.Run()
	if err != nil {
		Edi.Logger.Fatal(err.Error())
	}
}
