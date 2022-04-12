package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-joe/joe"
)

type Trivia struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func getTrivia() (Trivia, error) {
	q := Trivia{}
	url := "https://trivia.fyi/random-trivia-questions/"
	resp, err := http.Get(url)
	if err != nil {
		return q, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return q, err
	}
	doc.Find("a.query-title-link").Each(func(i int, sel *goquery.Selection) {
		q.Question = strings.TrimSpace(sel.Text())
		return
	})

	doc.Find("div.su-spoiler-content").Each(func(i int, sel *goquery.Selection) {
		q.Answer = strings.TrimSpace(sel.Text())
	})
	return q, nil
}

func TriviaQuestion(msg joe.Message) error {
	t, err := getTrivia()
	fmt.Println(t)
	return err
}
