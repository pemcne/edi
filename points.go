package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-joe/joe"
)

const pointsStoreKey string = "points"

type Point struct {
	Id    string `json:"id"`
	Value int    `json:"value"`
}

type PointList []Point

func (p PointList) Len() int {
	return len(p)
}

func (p PointList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PointList) Less(i, j int) bool {
	return p[i].Value > p[j].Value
}

var allPoints map[string]Point

func Points(msg joe.Message) error {
	Edi.Logger.Info("Getting from store")
	ok, err := Edi.Store.Get(pointsStoreKey, &allPoints)
	if err != nil {
		return err
	}
	if !ok {
		Edi.Logger.Info("setting defaults")
		allPoints = make(map[string]Point)
	}

	symbol := msg.Matches[0]
	amount, err := strconv.Atoi(msg.Matches[1])
	if err != nil {
		return err
	}
	key := strings.TrimSpace(msg.Matches[3])
	Edi.Logger.Info("Got the key: " + key)

	p, ok := allPoints[key]
	if !ok {
		p = Point{
			Id:    key,
			Value: 0,
		}
	}
	if symbol == "+" {
		p.Value += amount
	} else {
		p.Value -= amount
	}
	allPoints[key] = p
	Edi.Store.Set(pointsStoreKey, allPoints)

	msg.Respond("%s%d to %s = %d", symbol, amount, p.Id, p.Value)

	return nil
}

func PointsScore(msg joe.Message) error {
	_, err := Edi.Store.Get(pointsStoreKey, &allPoints)
	if err != nil {
		return err
	}
	key := strings.TrimSpace(msg.Matches[0])
	if val, ok := allPoints[key]; ok {
		msg.Respond("Points for %s is %d", val.Id, val.Value)
	}
	return nil
}

func PointsLeaderboard(msg joe.Message) error {
	_, err := Edi.Store.Get(pointsStoreKey, &allPoints)
	if err != nil {
		return err
	}
	sortlist := make(PointList, len(allPoints))
	i := 0
	for _, v := range allPoints {
		sortlist[i] = v
		i++
	}

	sort.Sort(sortlist)
	output := ""
	counter := 0
	for _, k := range sortlist {
		output += fmt.Sprintf("%s: %d", k.Id, k.Value)
		counter++
		if counter > 2 {
			break
		} else {
			output += "\n"
		}
	}
	msg.Respond(output)

	return nil
}
