package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-joe/joe"
)

const (
	Wordle      string = "wordle"
	Dordle      string = "dordle"
	Quordle     string = "quordle"
	Octordle    string = "octordle"
	Worldle     string = "worldle"
	Tradle      string = "tradle"
	Explordle   string = "explordle"
	Connections string = "connections"
	Strands     string = "strands"
)

type UserScore struct {
	Games map[string]GameScore `json:"games"`
}

type GameScore struct {
	Score float64 `json:"score"`
	Games int     `json:"games"`
}

const wordleBrainKey string = "wordle.scores"

var emojiTranslate map[string]string = map[string]string{
	":one:":        "1",
	":two:":        "2",
	":three:":      "3",
	":four:":       "4",
	":five:":       "5",
	":six:":        "6",
	":seven:":      "7",
	":eight:":      "8",
	":nine:":       "9",
	":keycap_ten:": "10",
	":clock11:":    "11",
	":clock12:":    "12",
	":clock1:":     "13",
}

var attempts map[string]int = map[string]int{
	Wordle:      6,
	Dordle:      7,
	Quordle:     9,
	Octordle:    13,
	Worldle:     6,
	Tradle:      6,
	Explordle:   7,
	Connections: 4,
	Strands:     0,
}

func computeAverage(scores []string) (float64, error) {
	var total float64
	for _, score := range scores {
		s, err := strconv.ParseFloat(score, 32)
		if err != nil {
			return 0.0, err
		}
		total += s
	}
	length := float64(len(scores))
	average := math.Round(total/length*100) / 100
	return average, nil
}

func processGame(user, game string, scores []string) error {
	// Get the scores
	var allScores map[string]UserScore = make(map[string]UserScore)
	_, err := Edi.Store.Get(wordleBrainKey, &allScores)
	if err != nil {
		return err
	}

	var userScores UserScore
	if v, ok := allScores[user]; ok {
		userScores = v
	}

	// Get the average
	average, err := computeAverage(scores)
	if err != nil {
		return err
	}
	if len(userScores.Games) == 0 {
		userScores.Games = make(map[string]GameScore)
	}
	// Add the average to the scores
	if _, ok := userScores.Games[game]; !ok {
		userScores.Games[game] = GameScore{}
	}
	g := userScores.Games[game]
	g.Games++
	g.Score += average
	userScores.Games[game] = g

	// Now store it all back into memory
	allScores[user] = userScores
	Edi.Store.Set(wordleBrainKey, allScores)
	return nil
}

func scoreConnections(game string) (int, error) {
	lines := strings.Split(game, "\n")
	misses := 0
	for _, line := range lines {
		r := regexp.MustCompile(`:large_([a-z]+)_square:`)
		colors := r.FindAllStringSubmatch(line, -1)
		var color string
		match := true
		for i, v := range colors {
			if i == 0 {
				color = v[1]
			} else if color != v[1] {
				match = false
				break
			}
		}
		if !match {
			misses++
		}
	}
	return misses, nil
}

func scoreStrands(game string) (int, error) {
	// Penalize hints more
	const hintScore int = 2

	lines := strings.Split(game, "\n")
	hints := 0
	count := 0
	var spanagram int

	for _, line := range lines {
		r := regexp.MustCompile(`:(\w+):`)
		symbols := r.FindAllStringSubmatch(line, -1)
		for _, v := range symbols {
			fmt.Println("Symbol " + v[1])
			if v[1] == "bulb" {
				hints++
			} else if v[1] == "large_yellow_circle" {
				spanagram = count
			} else {
				count++
			}
		}
	}
	score := (hints * hintScore) + spanagram
	fmt.Printf("Hints: %d, span: %d, count: %d\n", hints, spanagram, count)
	return score, nil
}

// Wordle\s\d+\s(.+)/\d
func WordleScore(msg joe.Message) error {
	const game string = Wordle

	user := msg.AuthorID
	score := strings.TrimSpace(msg.Matches[0])
	if score == "X" {
		score = strconv.Itoa(attempts[game] + 1)
	}
	scores := []string{score}
	err := processGame(user, game, scores)

	return err
}

// Dordle\s#\d+\s(.+)/\d
func DordleScore(msg joe.Message) error {
	const game string = Dordle

	user := msg.AuthorID
	score := msg.Matches[0]
	scores := strings.Split(score, "&amp;")
	for i, v := range scores {
		if v == "X" {
			scores[i] = strconv.Itoa(attempts[game] + 1)
		}
	}
	err := processGame(user, game, scores)
	return err
}

// Quordle\\s\\d+\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)
func QuordleScore(msg joe.Message) error {
	const game string = Quordle

	user := msg.AuthorID
	scores := msg.Matches
	for i, v := range scores {
		if v == ":large_red_square" {
			scores[i] = strconv.Itoa(attempts[game] + 1)
		} else {
			scores[i] = emojiTranslate[v]
		}
	}
	err := processGame(user, game, scores)
	return err
}

// Octordle\\s\\d+\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)\\s+(:.+:)(:.+:)
func OctordleScore(msg joe.Message) error {
	const game string = Octordle

	user := msg.AuthorID
	scores := msg.Matches
	for i, v := range scores {
		if v == ":large_red_square:" {
			scores[i] = strconv.Itoa(attempts[game] + 1)
		} else {
			scores[i] = emojiTranslate[v]
		}
	}
	err := processGame(user, game, scores)
	return err
}

func WorldleScore(msg joe.Message) error {
	const game string = Worldle

	user := msg.AuthorID
	score := strings.TrimSpace(msg.Matches[0])
	if score == "X" {
		score = strconv.Itoa(attempts[game] + 1)
	}
	scores := []string{score}
	err := processGame(user, game, scores)

	return err
}

func TradleScore(msg joe.Message) error {
	const game string = Tradle

	user := msg.AuthorID
	score := strings.TrimSpace(msg.Matches[0])
	if score == "X" {
		score = strconv.Itoa(attempts[game] + 1)
	}
	scores := []string{score}
	err := processGame(user, game, scores)

	return err
}

func ExplordleScore(msg joe.Message) error {
	const game string = Explordle

	user := msg.AuthorID
	score := strings.TrimSpace(msg.Matches[0])
	if score == "x" {
		score = strconv.Itoa(attempts[game] + 1)
	}
	scores := []string{score}
	err := processGame(user, game, scores)

	return err
}

func ConnectionScore(msg joe.Message) error {
	const game string = Connections

	user := msg.AuthorID

	gamestr := msg.Matches[0]
	misses, err := scoreConnections(gamestr)
	if err != nil {
		return err
	}
	scores := []string{fmt.Sprintf("%d", misses)}
	err = processGame(user, game, scores)
	return err
}

func StrandsScore(msg joe.Message) error {
	const game string = Strands

	user := msg.AuthorID

	gamestr := msg.Matches[0]
	score, err := scoreStrands(gamestr)
	if err != nil {
		return err
	}
	scores := []string{fmt.Sprintf("%d", score)}
	err = processGame(user, game, scores)
	return err
}

func WordleStats(msg joe.Message) error {
	user := msg.AuthorID
	var allScores map[string]UserScore = make(map[string]UserScore)
	_, err := Edi.Store.Get(wordleBrainKey, &allScores)
	if err != nil {
		return err
	}
	var gameOrder = []string{
		Wordle,
		Dordle,
		Quordle,
		Octordle,
		Worldle,
		Tradle,
		Explordle,
		Connections,
		Strands,
	}
	if v, ok := allScores[user]; ok {
		var output []string
		for _, game := range gameOrder {
			if gameData, ok := v.Games[game]; ok {
				attempt := attempts[game]
				avg := gameData.Score / float64(gameData.Games)
				gamestr := fmt.Sprintf("%s: %.3f / %d [%d games]", game, avg, attempt, gameData.Games)
				output = append(output, gamestr)
			}
		}
		msg.Respond(strings.Join(output, "\n"))
	} else {
		msg.Respond("No stats yet")
	}
	return nil
}
