package main

import (
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-joe/joe"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Trivia struct {
	Question  string   `json:"question"`
	RawAnswer string   `json:"answer"`
	Answers   []string `json:"altanswers"`
}

const triviaStoreKey string = "edi.trivia"

func pruneAnswer(answer string) []string {
	// This just compiles a list of possible answers to assist
	// the fuzzy match
	var output []string = []string{answer}
	// Filter out any articles
	articleReg := regexp.MustCompile(`(?i)\b(a|an|the)\b\s+`)
	articlePrune := articleReg.ReplaceAllString(answer, "")
	if articlePrune != answer {
		output = append(output, articlePrune)
	}

	// Filter out any punctuation or special chars
	puncReg := regexp.MustCompile(`(\.|,|&|!|-|")`)
	puncPrune := puncReg.ReplaceAllString(answer, "")
	if puncPrune != answer {
		output = append(output, puncPrune)
	}

	output = append(output, strings.ToLower(answer))
	return output
}

func newTriviaQuestion(t *Trivia) error {
	url := "https://trivia.fyi/random-trivia-questions/"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	t.Question = strings.TrimSpace(doc.Find("a.query-title-link").First().Text())
	t.RawAnswer = strings.TrimSpace(doc.Find("div.su-spoiler-content").First().Text())
	t.Answers = pruneAnswer(t.RawAnswer)
	Edi.Logger.Info("Trivia question: " + t.Question + " - " + t.RawAnswer)
	return nil
}

func askQuestion(msg *joe.Message) error {
	q := Trivia{}
	err := newTriviaQuestion(&q)
	if err != nil {
		return err
	}
	Edi.Store.Set(triviaStoreKey, q)
	msg.Respond(q.Question)
	return nil
}

func answerQuestion(correct bool, msg *joe.Message) error {
	q := Trivia{}
	ok, err := Edi.Store.Get(triviaStoreKey, &q)
	if err != nil {
		return err
	}
	if !ok {
		msg.Respond("Need to ask a question first")
		return nil
	}
	correctStr := ""
	if correct {
		correctStr = "Correct!! "
	}
	msg.Respond("%s%s -- %s", correctStr, q.Question, q.RawAnswer)
	return askQuestion(msg)
}

func TriviaQuestion(msg joe.Message) error {
	err := askQuestion(&msg)
	if err != nil {
		return err
	}
	return err
}

func TriviaAnswer(msg joe.Message) error {
	err := answerQuestion(false, &msg)
	if err != nil {
		return err
	}
	return nil
}

func TriviaGuess(msg joe.Message) error {
	q := Trivia{}
	ok, err := Edi.Store.Get(triviaStoreKey, &q)
	if err != nil {
		return err
	}
	if ok {
		matches := fuzzy.RankFindFold(msg.Text, q.Answers)
		// Set up an 85% rule for fuzzy matching
		acceptDiff := int(math.Round(float64(len(q.RawAnswer)) * 0.85))
		accept := len(q.RawAnswer) - acceptDiff
		if accept < 3 {
			accept = 3
		}
		for _, m := range matches {
			if m.Distance <= accept {
				return answerQuestion(true, &msg)
			}
		}
	}
	return nil
}

func TriviaGiveUp(msg joe.Message) error {
	msg.Respond("Sorry :(")
	return answerQuestion(false, &msg)
}
