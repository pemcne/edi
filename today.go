package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-joe/joe"
)

type TodayCache struct {
	Key   string    `json:"key"`
	Cache OnThisDay `json:"cache"`
}

type OnThisDay struct {
	Wikipedia string           `json:"wikipedia"`
	Date      string           `json:"date"`
	Events    []OnThisDayEvent `json:"events"`
}

type OnThisDayEvent struct {
	Year        string          `json:"year"`
	Description string          `json:"description"`
	Wikipedia   []OnThisDayLink `json:"wikipedia"`
}

type OnThisDayLink struct {
	Title     string `json:"title"`
	Wikipedia string `json:"wikipedia"`
}

func Today(msg joe.Message) error {
	t := time.Now()
	urlDate := fmt.Sprintf("%d/%d", int(t.Month()), t.Day())

	var cache TodayCache
	ok, err := Edi.Store.Get("today.cache", &cache)
	if err != nil {
		return err
	}
	var events OnThisDay
	if !ok || cache.Key != urlDate {
		Edi.Logger.Info("Refreshing today cache data")
		url := fmt.Sprintf("https://byabbe.se/on-this-day/%s/events.json", urlDate)
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(content, &events)
		if err != nil {
			return err
		}

		cache = TodayCache{
			Key:   urlDate,
			Cache: events,
		}
		Edi.Store.Set("today.cache", cache)
	} else {
		Edi.Logger.Info("Pulled data from cache")
		events = cache.Cache
	}

	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	randEvent := events.Events[gen.Intn(len(events.Events)-1)]
	output := fmt.Sprintf("On this day in %s: %s", randEvent.Year, randEvent.Description)
	msg.Respond(output)
	return nil
}
