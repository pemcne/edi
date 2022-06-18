package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter/v2"
	"go.uber.org/zap/zapcore"
)

var Edi *joe.Bot

var RandomGenerator *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func correctRoom(msg joe.Message, rooms []string) bool {
	for _, i := range rooms {
		if i == msg.Channel {
			return true
		}
	}
	return false
}

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

	CommonResponses(Edi)

	// Coin
	Edi.Respond("flip a coin", CoinFlip)

	// Points
	Edi.Respond("(\\+|-)\\s*(\\d+) (to|for) (.+)", Points)
	Edi.Respond("score for (.+)", PointsScore)
	Edi.Respond("leaderboard", PointsLeaderboard)
	Edi.Respond("what happened today", Today)

	// Trivia
	Edi.Respond("trivia question", TriviaQuestion)
	Edi.Respond("trivia answer", TriviaAnswer)
	Edi.Respond("fuck this question", TriviaGiveUp)
	Edi.Hear(".+", TriviaGuess)

	// Wordle
	Edi.Hear(`Wordle\s\d+\s(.+)/\d`, WordleScore)
	Edi.Hear(`Dordle\s#\d+\s(.+)/\d`, DordleScore)
	Edi.Hear(`Quordle\s\d+\s*(:.+:)(:.+:)\s+(:.+:)(:.+:)`, QuordleScore)
	Edi.Hear(`Octordle\s#\d+\s*(:.+:)(:.+:)\s+(:.+:)(:.+:)\s+(:.+:)(:.+:)\s+(:.+:)(:.+:)`, OctordleScore)
	Edi.Hear(`Worldle\s#\d+\s(.+)/\d`, WorldleScore)
	Edi.Hear(`Tradle\s#\d+\s(.+)/\d`, TradleScore)
	Edi.Hear(`Explordle\s\d+\s(.+)/\d`, ExplordleScore)
	Edi.Respond("wordle stats", WordleStats)

	// Huwordle
	err := loadHuwordleFiles()
	if err != nil {
		Edi.Logger.Error(err.Error())
	}
	Edi.Respond("huwordle new", HuwordleNew)
	Edi.Hear(`^\w+$`, HuwordleGuess)

	// Schedules
	err = cronInit()
	if err != nil {
		Edi.Logger.Error(err.Error())
	}
	Edi.Respond(`schedule (?:new|add)(?: <#(\w+)\|.+>)? "(.*?)" ((?:.|\s)*)$`, ScheduleNew)
	Edi.Respond(`schedule list`, ScheduleList)
	Edi.Respond(`schedule (remove|delete) (\d+)`, ScheduleRemove)

	// Chess
	err = initChess(false)
	if err != nil {
		Edi.Logger.Error(err.Error())
	}
	Edi.Respond("chess new", ChessNew)
	Edi.Respond("chess state", ChessState)
	Edi.Hear(`^(\S+)$`, ChessMove)
	Edi.Respond(`chess elo set (\d+)`, ChessElo)
	defer Engine.Close()

	err = Edi.Run()
	if err != nil {
		Edi.Logger.Fatal(err.Error())
	}
}
