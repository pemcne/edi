package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-joe/joe"
)

//go:embed wordle-data/answers.json
var wordleAnswers string
var answers []string

//go:embed wordle-data/weights.json
var rawWeights string
var weights map[int][26]float64

const wordleSolveBrainKey string = "wordle.solve"

type Best struct {
	Word   string
	Weight float64
}

type WordState struct {
	history map[int][]SolveLetterState
	present []rune
	absent  []rune
}

type SolveLetterState struct {
	green  bool
	yellow bool
	letter rune
}

type Solution struct {
	Solution string `json:"solution"`
	GameId   int    `json:"days_since_launch"`
	Date     string `json:"print_date"`
}

func getSolution(solution *Solution) error {
	now := time.Now()
	current := now.Format("2006-01-02")
	_, err := Edi.Store.Get(wordleSolveBrainKey, solution)
	if err != nil {
		return err
	}
	if solution.Date != current {
		url := fmt.Sprintf("https://www.nytimes.com/svc/wordle/v2/%s.json", current)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(solution)
		if err != nil {
			return err
		}
		Edi.Store.Set(wordleSolveBrainKey, solution)
	}
	return nil
}

func nextGuess(state *WordState) string {
	best := Best{}
	// Loop through every word in the answers
	// TODO: can probably trim this down over time
WORD:
	for _, word := range answers {
		score := 0.0
		seen := make([]bool, 26)
		yellows := make([]rune, len(state.present))
		copy(yellows, state.present)
		// For each letter in the word, compare that to the
		// letter state
		for idx, letter := range word {
			letterScore := 0.0
			history := state.history[idx]
			letterIdx := int(letter - 'a')
			for _, histLetter := range history {
				if histLetter.green {
					if histLetter.letter == letter {
						letterScore += weights[idx][letterIdx] * 10.0
						break
					} else {
						// Hard mode
						// Have a green that doesn't match the current letter, skip entirely
						continue WORD
					}
				}
				if histLetter.letter == letter && histLetter.yellow {
					// Seen this letter as a yellow here, skip
					continue WORD
				}
			}
			// This is the fall through if no history found, normal weight
			if letterScore == 0.0 {
				missed := inArray(state.absent, letter)
				multi := 1.0
				present := inArray(yellows, letter)
				if present != -1 {
					// The letter might be somewhere in here
					letterScore += weights[idx][letterIdx] * 2.0
					yellows = arrayRemove(yellows, present)
				} else if missed != -1 {
					// If the letter is gray, skip
					continue WORD
				} else if seen[letterIdx] {
					// Just a small padding to discourage repeat letters
					multi = -0.25
				}
				letterScore += weights[idx][letterIdx] * multi
			}
			seen[letterIdx] = true
			score += letterScore
		}
		// Add some randomization to the first guess
		if len(state.history[0]) == 0 {
			score = score * rand.Float64()
		}
		// Hard mode check, make sure you're using all info
		if len(yellows) > 0 {
			continue WORD
		}
		if best.Weight < score {
			best.Weight = score
			best.Word = word
		}
	}
	return best.Word
}

func scoreGuess(guess, solution string, state *WordState) bool {
	guessChars := []rune(guess)
	solChars := []rune(solution)
	for idx, letter := range guessChars {
		if letter == solChars[idx] {
			history := SolveLetterState{
				green:  true,
				letter: letter,
			}
			state.history[idx] = append([]SolveLetterState{history}, state.history[idx]...)
			solChars[idx] = '_'
			guessChars[idx] = '_'
		}
	}
	if guess == solution {
		return true
	}
	state.present = make([]rune, 0)
	for idx, letter := range guessChars {
		if letter == '_' {
			continue
		}
		found := inArray(solChars, letter)
		if found != -1 {
			history := SolveLetterState{
				yellow: true,
				letter: letter,
			}
			state.history[idx] = append([]SolveLetterState{history}, state.history[idx]...)
			solChars[found] = '_'
			state.present = append(state.present, letter)
		} else {
			history := SolveLetterState{
				letter: letter,
			}
			// inPresent := inArray(letter, state.present)
			inAbsent := inArray(state.absent, letter)
			state.history[idx] = append([]SolveLetterState{history}, state.history[idx]...)
			if inAbsent == -1 {
				state.absent = append(state.absent, letter)
			}
		}
	}
	return false
}

func printGame(state *WordState, game int, correct bool) string {
	total := "X"
	if correct {
		total = fmt.Sprintf("%d", len(state.history[0]))
	}
	out := fmt.Sprintf("Wordle %d %s/6*\n", game, total)
	for guess := len(state.history[0]) - 1; guess >= 0; guess-- {
		line := ""
		for letter := 0; letter < 5; letter++ {
			hist := state.history[letter][guess]
			if hist.green {
				line += ":large_green_square:"
			} else if hist.yellow {
				line += ":large_yellow_square:"
			} else {
				line += ":black_large_square:"
			}
		}
		line += "\n"
		out += line
	}
	return out
}

func loadWordleSolveFiles() error {
	Edi.Logger.Info("Loading worlde solve files")
	err := json.Unmarshal([]byte(wordleAnswers), &answers)
	if err != nil {
		return err
	}

	// Predetermined by calculating the frequency of every letter in the word list
	// Weights are calculated by taking the frequency of each letter for each
	// slot position
	// eg Freq of the letter 'a' being in the first letter slot
	err = json.Unmarshal([]byte(rawWeights), &weights)
	if err != nil {
		return err
	}
	return nil
}

func SolveWordle(msg joe.Message) error {
	solution := &Solution{}
	err := getSolution(solution)
	if err != nil {
		return err
	}
	state := WordState{}
	state.history = make(map[int][]SolveLetterState)
	for i := 0; i < 5; i++ {
		state.history[i] = make([]SolveLetterState, 0)
	}
	guesses := []string{}
	correct := false
	for i := 0; i < 6; i++ {
		guess := nextGuess(&state)
		guesses = append(guesses, guess)
		correct = scoreGuess(guess, solution.Solution, &state)
		if correct {
			break
		}
	}
	msg.Respond(printGame(&state, solution.GameId, correct))
	outGuesses, err := json.Marshal(guesses)
	if err != nil {
		return err
	}
	Edi.Logger.Info("Wordle solve guesses: " + string(outGuesses))
	return nil
}
