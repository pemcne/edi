package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-joe/joe"
	"github.com/go-joe/joe/reactions"
)

const CORRECT string = ":large_green_square:"
const PRESENT string = ":large_yellow_square:"
const ABSENT string = ":black_large_square:"

var HUWORDLEROOMS = []string{
	"", // terminal
	"C035YQ3UG79", // personal
	"C033N9SPX33", // testing
	"C04GESW4WKS", // rdc
}

const huwordleStoreKey string = "huwordle.word"

//go:embed huwordle-data/answers.json
var rawAnswers string

//go:embed huwordle-data/dictionary.json
var rawDictionary string

var ANSWERS []string
var DICTIONARY []string

type LetterState struct {
	Present []string `json:"present"`
	Absent  []string `json:"absent"`
	Unknown []string `json:"unknown"`
}

type HuwordleState struct {
	Word    string      `json:"word"`
	Guesses int         `json:"guesses"`
	State   []string    `json:"state"`
	Letters LetterState `json:"letterState"`
}

func loadHuwordleFiles() error {
	Edi.Logger.Info("Loading the huwordle files")
	if len(ANSWERS) == 0 || len(DICTIONARY) == 0 {
		err := json.Unmarshal([]byte(rawAnswers), &ANSWERS)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(rawDictionary), &DICTIONARY)
		if err != nil {
			return err
		}
	}
	return nil
}

func printState(state *HuwordleState, won, first bool) string {
	length := len(state.Word)
	output := ""
	if !won {
		if first {
			output += "New word is "
		}
		output += fmt.Sprintf("%d letters: %s", length, strings.Join(state.State, ""))
		if !first {
			presentLetters := strings.Join(state.Letters.Present, "")
			absentLetters := strings.Join(state.Letters.Absent, "")
			output += fmt.Sprintf("\n%s: %s", PRESENT, presentLetters)
			output += fmt.Sprintf("\n%s: %s", ABSENT, absentLetters)
		}
	} else {
		output += fmt.Sprintf("%d total guesses", state.Guesses)
	}
	return output
}

func emoji(word []string) string {
	output := ""
	for _, l := range word {
		output += fmt.Sprintf(":alphabet-yellow-%s:", l)
	}
	return output
}

func arrayIn(arr []string, el string) int {
	index := -1
	for k, v := range arr {
		if v == el {
			index = k
		}
	}
	return index
}

func arrayRemove(arr []string, i int) []string {
	if i == -1 {
		return arr
	} else if i == len(arr)-1 {
		return arr[:len(arr)-1]
	} else {
		copy(arr[i:], arr[i+1:])
		arr[len(arr)-1] = ""
		arr = arr[:len(arr)-1]
		return arr
	}
}

func processGuess(guess string, state *HuwordleState) []string {
	word := state.Word
	wordchars := strings.Split(word, "")
	guesschars := strings.Split(guess, "")
	output := make([]string, len(wordchars))
	letters := state.Letters
	state.Guesses++

	for i, wordLetter := range wordchars {
		guessLetter := guesschars[i]
		if guessLetter == wordLetter {
			output[i] = CORRECT
			guesschars[i] = ""
			wordchars[i] = ""
			// Update letter state
			state.State[i] = emoji([]string{wordLetter})
			letters.Unknown[i] = ""
			if arrayIn(letters.Unknown, guessLetter) == -1 {
				letters.Present = arrayRemove(letters.Present, arrayIn(letters.Present, guessLetter))
			}
		}
	}
	for i := range wordchars {
		guessLetter := guesschars[i]
		if guessLetter != "" {
			index := arrayIn(wordchars, guessLetter)
			if index == -1 {
				// Letter isn't in the word
				output[i] = ABSENT
				inAbsent := arrayIn(letters.Absent, guessLetter)
				inPresent := arrayIn(letters.Present, guessLetter)
				if inAbsent == -1 && inPresent == -1 {
					letters.Absent = append(letters.Absent, guessLetter)
				}
			} else {
				// Letter is in the word but not the right spot
				output[i] = PRESENT
				wordchars[index] = ""
				if arrayIn(letters.Present, guessLetter) == -1 {
					letters.Present = append(letters.Present, guessLetter)
				}
			}
		}
	}
	sort.Strings(letters.Present)
	sort.Strings(letters.Absent)
	state.Letters = letters

	return output
}

func getState() (HuwordleState, error) {
	state := HuwordleState{}
	_, err := Edi.Store.Get(huwordleStoreKey, &state)
	if err != nil {
		return state, err
	}
	return state, nil
}

func setState(state HuwordleState) error {
	err := Edi.Store.Set(huwordleStoreKey, state)
	return err
}

func newWord(msg *joe.Message) error {
	randNum := RandomGenerator.Intn(len(ANSWERS))
	word := ANSWERS[randNum]
	var wordState []string
	for i := 0; i < len(word); i++ {
		wordState = append(wordState, ABSENT)
	}
	state := HuwordleState{
		Word:  word,
		State: wordState,
		Letters: LetterState{
			Unknown: strings.Split(word, ""),
		},
	}
	Edi.Logger.Info("Huwordle word: " + word)
	msg.Respond(printState(&state, false, true))
	return setState(state)
}

func HuwordleNew(msg joe.Message) error {
	if !correctRoom(msg, HUWORDLEROOMS) {
		return nil
	}
	state, err := getState()
	if err != nil {
		return err
	}
	if state.Word != "" {
		msg.Respond("Previous word was: " + state.Word)
	}
	newWord(&msg)

	return nil
}

func HuwordleGuess(msg joe.Message) error {
	if !correctRoom(msg, HUWORDLEROOMS) {
		return nil
	}
	state, err := getState()
	if err != nil {
		return err
	}

	guess := strings.ToLower(strings.TrimSpace(msg.Text))
	won := state.Word == guess
	if !won && (len(guess) != len(state.Word)) {
		return nil
	}
	if arrayIn(DICTIONARY, guess) == -1 {
		msg.React(reactions.Reaction{
			Shortcode: "x",
		})
		return nil
	}
	Edi.Logger.Debug("Huwordle: processing guess of '" + guess + "'")
	results := processGuess(guess, &state)
	guessEmoji := emoji(strings.Split(guess, ""))
	output := fmt.Sprintf("%s\n%s\n", guessEmoji, strings.Join(results, ""))
	output += printState(&state, won, false)
	msg.Respond(output)
	if won {
		return newWord(&msg)
	}
	return setState(state)
}
