package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-joe/joe"
	"github.com/robfig/cron/v3"
)

var schedule *cron.Cron

const scheduleStoreKey = "schedule"

type Schedule struct {
	Pattern string `json:"pattern"`
	Message string `json:"messasge"`
	Channel string `json:"channel"`
}

var scheduleState []Schedule

var allSchedules = make(map[cron.EntryID]Schedule)

func cronInit() error {
	schedule = cron.New()
	schedule.Start()
	return loadSchedules()
}

func loadSchedules() error {
	_, err := Edi.Store.Get(scheduleStoreKey, &scheduleState)
	if err != nil {
		return err
	}
	for _, s := range scheduleState {
		_, _, err := addSchedule(s.Pattern, s.Message, s.Channel)
		if err != nil {
			return err
		}
	}
	return nil
}

func addSchedule(pattern, message, channel string) (cron.EntryID, Schedule, error) {
	cronId, err := schedule.AddFunc(pattern, func() {
		Edi.Adapter.Send(message, channel)
	})
	if err != nil {
		return 0, Schedule{}, err
	}
	s := Schedule{
		Pattern: pattern,
		Message: message,
		Channel: channel,
	}
	allSchedules[cronId] = s
	return cronId, s, nil
}

func removeFromScheduleArray(sch Schedule) {
	fmt.Println(scheduleState)
	if len(scheduleState) == 1 {
		scheduleState = scheduleState[:0]
		return
	}
	index := -1
	for i, v := range scheduleState {
		if v.Channel == sch.Channel && v.Message == sch.Message && v.Pattern == sch.Pattern {
			fmt.Println("Found match")
			fmt.Println(i)
			index = i
		}
	}
	if index != -1 {
		scheduleState = append(scheduleState[:index], scheduleState[index+1:]...)
	}
	fmt.Println(scheduleState)
}

func removeSchedule(id cron.EntryID) error {
	removeFromScheduleArray(allSchedules[id])
	schedule.Remove(id)
	delete(allSchedules, id)
	return Edi.Store.Set(scheduleStoreKey, scheduleState)
}

func ScheduleNew(msg joe.Message) error {
	// Unpack the matches
	chn := msg.Matches[0]
	ptn := msg.Matches[1]
	txt := msg.Matches[2]
	id, sch, err := addSchedule(ptn, txt, chn)
	if err != nil {
		return err
	}
	scheduleState = append(scheduleState, sch)
	err = Edi.Store.Set(scheduleStoreKey, scheduleState)
	msg.Respond(fmt.Sprintf("Added schedule: %d", id))
	return err
}

func ScheduleRemove(msg joe.Message) error {
	msgId, err := strconv.Atoi(msg.Matches[1])
	if err != nil {
		return err
	}
	id := cron.EntryID(msgId)

	_, ok := allSchedules[id]
	if !ok {
		msg.Respond(fmt.Sprintf("No schedule with id %d", id))
	} else {
		err = removeSchedule(id)
		msg.Respond(fmt.Sprintf("Removed schedule %d", id))
	}
	return err
}

func ScheduleList(msg joe.Message) error {
	if len(allSchedules) == 0 {
		msg.Respond("No schedules set yet")
	}
	var output []string
	for id, sch := range allSchedules {
		line := fmt.Sprintf("%d: [%s] #%s %s", id, sch.Pattern, sch.Channel, sch.Message)
		output = append(output, line)
	}
	msg.Respond(strings.Join(output, "\n"))
	return nil
}
