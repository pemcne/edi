package main

import (
	"os"

	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter/v2"
	"go.uber.org/zap/zapcore"
)

var Edi *joe.Bot

func loadModules() []joe.Module {
	var modules []joe.Module
	// See if we want to load Slack
	if token, t_ok := os.LookupEnv("SLACK_TOKEN"); t_ok {
		if app_token, a_ok := os.LookupEnv("SLACK_APP_TOKEN"); a_ok {
			adapter := slack.EventsAPIAdapter(
				token,
				slack.WithSocketMode(app_token),
				slack.WithListenPassive(),
				// slack.WithDebug(true),
			)
			modules = append(modules, adapter)
		}
	}
	// Load store
	modules = append(modules, file.Memory("brain.json"))
	// For debugging
	modules = append(modules, joe.WithLogLevel(zapcore.DebugLevel))
	return modules
}

func main() {
	modules := loadModules()
	Edi = joe.New("Edi", modules...)

	// Ping
	Edi.Respond("ping", Pong)

	// Coin
	Edi.Respond("flip a coin", CoinFlip)

	// Points
	Edi.Respond("(\\+|-)\\s*(\\d+) (to|for) (.+)", Points)
	Edi.Respond("score for (.+)", PointsScore)
	Edi.Respond("leaderboard", PointsLeaderboard)
	Edi.Respond("what happened today", Today)

	// Trivia
	Edi.Respond("trivia", TriviaQuestion)

	// Wordle
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
