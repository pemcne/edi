package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter/v2"
	"github.com/pemcne/firestore-memory"

	slackapi "github.com/slack-go/slack"
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

func loadModules(name string) []joe.Module {
	var modules []joe.Module
	// See if we want to load Slack
	if token, t_ok := os.LookupEnv("SLACK_TOKEN"); t_ok {
		if app_token, a_ok := os.LookupEnv("SLACK_APP_TOKEN"); a_ok {
			adapter := slack.EventsAPIAdapter(
				token,
				slack.WithSocketMode(app_token),
				slack.WithListenPassive(),
				slack.WithMessageParams(slackapi.PostMessageParameters{
					Markdown: true,
					AsUser:   true,
				}),
				// slack.WithDebug(true),
			)
			modules = append(modules, adapter)
		}
	}
	// Load store
	if project, p_ok := os.LookupEnv("FIRESTORE_PROJECT"); p_ok {
		memOpt := firestore.WithCollection(name)
		modules = append(modules, firestore.Memory(project, memOpt))
	}

	// For debugging
	// modules = append(modules, joe.WithLogLevel(zapcore.DebugLevel))
	return modules
}

func main() {
	name, ok := os.LookupEnv("BOT_NAME")
	if !ok {
		name = "Edi"
	}
	modules := loadModules(name)
	Edi = joe.New(name, modules...)

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
	Edi.Hear(`Connections\sPuzzle #\d+\s([:\w:\s]+)`, ConnectionScore)
	Edi.Respond("wordle stats", WordleStats)

	// Wordle solver
	err := loadWordleSolveFiles()
	if err != nil {
		Edi.Logger.Error(err.Error())
	}
	Edi.Respond("wordle solve", SolveWordle)

	// Huwordle
	err = loadHuwordleFiles()
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
	Edi.Respond("chess info", ChessInfo)
	if Engine != nil {
		defer Engine.Close()
	}

	err = Edi.Run()
	if err != nil {
		Edi.Logger.Fatal(err.Error())
	}
}
