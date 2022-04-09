package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-joe/joe"
)

const (
	Wordle string = "wordle"
	Dordle string = "dordle"
)

type UserScore struct {
	Games map[string]GameScore `json:"games"`
}

type GameScore struct {
	Score float64 `json:"score"`
	Games int     `json:"games"`
}

const wordleBrainKey string = "wordle.scores"

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

	// Add the average to the scores
	if _, ok := userScores.Games[game]; !ok {
		userScores.Games = make(map[string]GameScore)
		userScores.Games[game] = GameScore{}
	}
	g := userScores.Games[game]
	g.Games++
	g.Score += average
	userScores.Games[game] = g

	// Now store it all back into memory
	allScores[user] = userScores
	Edi.Store.Set(wordleBrainKey, allScores)
	fmt.Println(allScores)
	return nil
}

// .*(Wordle|Dordle)\s\W?\d+\s(.+)\/\d
func WordleScore(msg joe.Message) error {
	const game string = Wordle
	const attempt string = "7"
	user := msg.AuthorID

	score := strings.TrimSpace(msg.Matches[0])
	if score == "X" {
		score = attempt
	}
	scores := []string{score}
	err := processGame(user, game, scores)

	return err
}
