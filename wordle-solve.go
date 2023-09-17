package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-joe/joe"
)

//go:embed wordle-data/answers.json
var wordleAnswers string
var answers []string

var weights [26]float64

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
}

func getSolution(solution *Solution) error {
	_, err := Edi.Store.Get(wordleSolveBrainKey, solution)
	if err != nil {
		return err
	}
	if solution.Solution == "" {
		now := time.Now()
		current := now.Format("2006-01-02")
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
	for _, word := range answers {
		score := 0.0
		seen := make([]bool, 26)
		yellows := make([]rune, len(state.present))
		copy(yellows, state.present)
		for idx, letter := range word {
			letterScore := 0.0
			history := state.history[idx]
			letterIdx := int(letter - 'a')
			for _, histLetter := range history {
				if histLetter.letter == letter {
					if histLetter.green {
						letterScore += weights[letterIdx] * 10.0
						break
					} else if histLetter.yellow {
						letterScore += weights[letterIdx] * -5.0
						break
					}
				}
			}
			// This is the fall through if no history found, normal weight
			if letterScore == 0.0 {
				present := inArray(yellows, letter)
				if present != -1 {
					letterScore += weights[letterIdx] * 2.0
					yellows = arrayRemove(yellows, present)
				}
				missed := inArray(state.absent, letter)
				multi := 1.0
				if missed != -1 {
					multi = -10.0
				} else if seen[letterIdx] {
					multi = -0.5
				}
				letterScore += weights[letterIdx] * multi
			}
			seen[letterIdx] = true
			score += letterScore
		}
		// Hard mode check, make sure you're using all info
		if len(yellows) > 0 {
			score = 0
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

func printGame(state *WordState, game int) string {
	out := fmt.Sprintf("Wordle %d %d/6*\n", game, len(state.history[0]))
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

func SolveWordle(msg joe.Message) error {
	err := json.Unmarshal([]byte(wordleAnswers), &answers)
	if err != nil {
		return err
	}

	// Predetermined by calculating the frequency of every letter in the answer list
	weights = [26]float64{0.4204863825696852, 0.11753925024030784, 0.1493015059275875, 0.179689202178789, 0.4560621595642354, 0.07815123357898135, 0.11993591797500823, 0.13152835629605908, 0.2778436398590192, 0.02149951938481256, 0.1107657801986544, 0.24555270746555496, 0.1463024671579624, 0.2183627042614546, 0.30796859980775115, 0.14750720922781183, 0.008330663248958666, 0.30723165652034456, 0.4692950977250802, 0.24038128804870088, 0.187526433835309, 0.05173021467478375, 0.08050304389618715, 0.02148349887856456, 0.15716116629285495, 0.030362063441204717}

	solution := &Solution{}
	err = getSolution(solution)
	if err != nil {
		return err
	}
	state := WordState{}
	state.history = make(map[int][]SolveLetterState)
	for i := 0; i < 5; i++ {
		state.history[i] = make([]SolveLetterState, 0)
	}
	guesses := []string{}
	for i := 0; i < 6; i++ {
		guess := nextGuess(&state)
		guesses = append(guesses, guess)
		correct := scoreGuess(guess, solution.Solution, &state)
		if correct {
			break
		}
	}
	msg.Respond(printGame(&state, solution.GameId))
	outGuesses, err := json.Marshal(guesses)
	if err != nil {
		return err
	}
	Edi.Logger.Info("Wordle solve guesses: " + string(outGuesses))
	return nil
}
