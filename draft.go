package main

import (
	"fmt"
	"github.com/lab-d8/lol-at-pitt/ols"
	"strconv"
	"time"
)

type DraftHandler func(msg Message, room *DraftRoom)

type DraftPlayer struct {
	ols.Player
	CurrentBid    int
	HighestBidder string
}

type Captain struct {
	ID     string // player relative marker
	Points int
	Team   string
}

type DraftClient struct {
	// TODO: Put maria's draft client in here.
}

/////////////////////////
const (
	startingCountdownTime = 20
	countdownEventTime    = 5
)

var (
	currentCountdown                         = startingCountdownTime
	allowTicks       bool                    = true // If this is false, dont continue to count down
	mainHandler      map[string]DraftHandler = map[string]DraftHandler{}
)

/////////////////////////

func RegisterDraftHandler(msg_type string, handle DraftHandler) {
	mainHandler[msg_type] = handle
}

func Handle(msg Message) {
	if mainHandler[msg.Type] != nil {
		mainHandler[msg.Type](msg, room)
	}
}

func Init() {
	timer_handler()
	RegisterDraftHandler("login", handle_login)
	RegisterDraftHandler("bid", handle_bid)
	RegisterDraftHandler("event", handle_event)
	RegisterDraftHandler("captains", handle_captains)
	RegisterDraftHandler("timer_reset", handle_timer_reset)
	RegisterDraftHandler("captains", handle_captains)
	RegisterDraftHandler("upcoming", handle_upcoming)
	//go timer()
}

func handle_bid(msg Message, room *DraftRoom) {
	// TODO: Update with maria code
	amt, err := strconv.Atoi(msg.Text)

	if err == nil {
		formattedStr := fmt.Sprintf("<h5>Amount: <span  class='text-success'>%d</span></h5>", amt)
		go Handle(Message{Type: "event", Text: formattedStr})
	}
}

func handle_event(msg Message, room *DraftRoom) {
	room.broadcast(&msg)
}

func handle_captains(msg Message, room *DraftRoom) {
	// TODO: do formatting of text here. Make it a json blob
	room.broadcast(&Message{Type: "captains", Text: "captains"})
}

func handle_upcoming(msg Message, room *DraftRoom) {
	//TODO: do formatting here!
	room.broadcast(&Message{Type: "upcoming", Text: "upcoming"})
}

// handle_login will give the player their stats, captains, current player, and upcoming players.
func handle_login(msg Message, room *DraftRoom) {
	Handle(Message{Type: "captains"})
	Handle(Message{Type: "upcoming"})
	Handle(Message{Type: "current-player"})
}

func handle_timer_reset(msg Message, room *DraftRoom) {
	currentCountdown = startingCountdownTime
}

func timer_handler() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			_ = now

			if !allowTicks {
				continue
			}

			currentCountdown--

			if currentCountdown < countdownEventTime {
				Handle(Message{Type: "event", Text: "counting down..."})
			}

			if currentCountdown == 0 {
				allowTicks = false
			}
		}
	}()
}
