package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/lab-D8/lol-at-pitt/draft"
	"github.com/lab-D8/lol-at-pitt/utils"
)

type DraftHandler func(msg Message, room *DraftRoom)

/////////////////////////
const (
	startingCountdownTime = 10
	countUpEventTime      = 10
	countdownEventTime    = 5
)

var (
	currentCountdown                         = startingCountdownTime
	allowTicks       bool                    = false // If this is false, dont continue to count down
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
	draft.Init("ols-summer-2016")
	timer_handler()

	RegisterDraftHandler("login", handle_update)
	RegisterDraftHandler("bid", handle_bid)
	RegisterDraftHandler("bid-more", handle_more_bid)
	RegisterDraftHandler("event", handle_event)
	RegisterDraftHandler("refresh", handle_refresh)
}

func handle_refresh(msg Message, room *DraftRoom) {
	//draft.Init()
	Handle(Message{Type: "update"})
}

func handle_more_bid(msg Message, room *DraftRoom) {
	amt, err := strconv.Atoi(msg.Text)
	log.Println(msg, err)
	if err == nil {
		amount := draft.Snapshot.CurrentPlayer.HighestBid + amt
		Handle(Message{Type: "bid", From: msg.From, Text: strconv.Itoa(amount)})
	}
}

func handle_bid(msg Message, room *DraftRoom) {
	amt, err := strconv.Atoi(msg.Text)
	log.Println(msg)
	if err == nil {
		event := draft.BidEvent{
			Amount:            amt,
			PlayerLeagueID:    draft.Snapshot.CurrentPlayer.LeagueId,
			DrafterFacebookID: msg.From,
		}
		bidSuccess := draft.ProcessBidEvent(event)
		if bidSuccess {
			go Handle(Message{Type: "event", Text: util.JsonStringifyIgnoreError(event)})
			currentCountdown = startingCountdownTime
			allowTicks = true
		}
	}
}

func handle_event(msg Message, room *DraftRoom) {
	room.broadcast(&msg)
}

func handle_captains(msg Message, room *DraftRoom) {
	room.broadcast(&Message{Type: "captains", Text: util.JsonStringifyIgnoreError(draft.Snapshot)})
}

func handle_upcoming(msg Message, room *DraftRoom) {
	room.broadcast(&Message{Type: "upcoming", Text: util.JsonStringifyIgnoreError(draft.Snapshot)})
}

func handle_bidder(msg Message, room *DraftRoom) {
	team := draft.Snapshot.GetTeamByFacebookId(msg.From)
	if team != nil {
		str := fmt.Sprintf("%d", team.Captain.Points)
		room.messageWithID(msg.From, &Message{Type: "points", Text: str})
		room.messageWithID(msg.From, &Message{Type: "team", Text: strconv.Itoa(team.Captain.Points)})
	}

}

// handle_login will give the player their stats, captains, current player, and upcoming players.
func handle_update(msg Message, room *DraftRoom) {
	Handle(Message{Type: "captains"})
	Handle(Message{Type: "upcoming"})
	Handle(Message{Type: "current-header"})
	//Handle(Message{Type: "event", Text: "Currently waiting to bid on.." + draft.Snapshot.CurrentPlayer.Ign})
	for _, client := range room.clients {
		Handle(Message{Type: "bidder", From: client.ID})
	}

}

func handle_winner(msg Message, room *DraftRoom) {
	Handle(Message{Type: "update"})
	draft.Paused = true
}

func handle_timer_end(msg Message, room *DraftRoom) {
	current := draft.Snapshot.CurrentPlayer
	if current.HighestBid > 0 {
		draft.Paused = true
		handle_winner(msg, room)
	}
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
				res := fmt.Sprintf("%d seconds remaining...", currentCountdown)
				Handle(Message{Type: "event", Text: res})
			}

			if currentCountdown == 0 {
				allowTicks = false
				Handle(Message{Type: "timer-end"})
			}
		}
	}()
}
