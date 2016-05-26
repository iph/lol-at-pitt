package draft

import (
	"github.com/lab-D8/lol-at-pitt/ols"
	"testing"
)

func TestInitNothing(t *testing.T) {
	Init("ols-summer-2016")

	if Snapshot.EventID != -1 {
		t.Error("Snapshot Event ID should be -1 when initialized empty.")
	}
}

func TestInitSavedData(t *testing.T) {
	currentData := LeagueDraftSnapshot{
		EventID: 100,
		ID:      "ols-summer-test",
	}

	db := ols.InitDB()
	collection := db.C("auctions")
	err := collection.Insert(currentData)

	if err != nil {
		panic(err)
	}

	Init("ols-summer-test")
	//collection.Remove(data)
	if Snapshot.EventID != 100 {
		t.Error("Snapshot Event ID should be 100, but was ", Snapshot.EventID)
	}
}

func TestProcessBidEventWrongId(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:        "woah",
			HighestBid: 10,
			Team:       "test",
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     101,
				},
			},
		},
	}

	event := BidEvent{
		Amount:            100,
		PlayerLeagueID:    1,
		DrafterFacebookID: "wrong_id",
	}

	result := ProcessBidEvent(event)

	if result {
		t.Error("Bid Event with a wrong id is fucked")
	}
}

func TestProcessBidEventHappy(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:        "woah",
			HighestBid: 10,
			Team:       "other_team",
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     101,
				},
			},
		},
	}

	event := BidEvent{
		Amount:            100,
		PlayerLeagueID:    1,
		DrafterFacebookID: "test_id",
	}

	result := ProcessBidEvent(event)

	if !result {
		t.Error("Bid Event with a wrong id is fucked")
	}

	if Snapshot.CurrentPlayer.Team != Snapshot.Teams[0].Name {
		t.Error("Bid Event did not register the team correctly")
	}
}

func TestProcessBidEventSameTeam(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:        "woah",
			HighestBid: 10,
			Team:       "test",
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     101,
				},
			},
		},
	}

	event := BidEvent{
		Amount:            100,
		PlayerLeagueID:    1,
		DrafterFacebookID: "test_id",
	}

	ProcessBidEvent(event)

	if Snapshot.CurrentPlayer.HighestBid != 10 {
		t.Error("Bid Event did not register the team correctly")
	}
}

func TestProcessBidEventNotEnoughPoints(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:        "woah",
			HighestBid: 10,
			Team:       "test",
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     99,
				},
			},
		},
	}

	event := BidEvent{
		Amount:            100,
		PlayerLeagueID:    1,
		DrafterFacebookID: "test_id",
	}

	ProcessBidEvent(event)

	if Snapshot.CurrentPlayer.HighestBid != 10 {
		t.Error("Bid was actually placed when it should not have been")
	}
}

func TestWinEventHappy(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:        "woah",
			HighestBid: 10,
			Team:       "test",
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     99,
				},
				Players: []*DraftPlayer{},
			},
		},
		FuturePlayers: []*DraftPlayer{
			&DraftPlayer{
				Ign: "next_player",
			},
		},
	}

	event := WinEvent{
		Amount:            10,
		PlayerLeagueID:    1,
		CaptainFacebookID: "test_id",
		TeamName:          "test",
	}

	ProcessWinEvent(event)

	if len(Snapshot.FuturePlayers) != 0 {
		t.Error("Future Players wasn't decremented off a win event")
	}

	if Snapshot.CurrentPlayer.Ign != "next_player" {
		t.Error("Current Player wasn't updated to new player")
	}

	if Snapshot.Teams[0].Captain.Points != 89 {
		t.Error("Bid wasnt decremented properly for win")
	}

	if len(Snapshot.Teams[0].Players) != 1 {
		t.Error("Player wasn't updated to new team.")
	}
}

func TestWinEventEnd(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:        "woah",
			HighestBid: 10,
			Team:       "test",
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     99,
				},
				Players: []*DraftPlayer{},
			},
		},
		FuturePlayers: []*DraftPlayer{},
	}

	event := WinEvent{
		Amount:            10,
		PlayerLeagueID:    1,
		CaptainFacebookID: "test_id",
		TeamName:          "test",
	}

	ProcessWinEvent(event)

	if len(Snapshot.FuturePlayers) != 0 {
		t.Error("Future Players wasn't decremented off a win event")
	}

	if Snapshot.CurrentPlayer != nil {
		t.Error("Current Player wasn't updated to new player")
	}

	if Snapshot.Teams[0].Captain.Points != 89 {
		t.Error("Bid wasnt decremented properly for win")
	}

	if len(Snapshot.Teams[0].Players) != 1 {
		t.Error("Player wasn't updated to new team.")
	}
}

func TestUnWinEventHappy(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign: "next_player",
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     89,
				},
				Players: []*DraftPlayer{
					&DraftPlayer{
						Ign:        "woah",
						HighestBid: 10,
						Team:       "test",
						LeagueId:   1,
					},
				},
			},
		},
		FuturePlayers: []*DraftPlayer{},
	}

	event := WinEvent{
		Amount:            10,
		PlayerLeagueID:    1,
		CaptainFacebookID: "test_id",
		TeamName:          "test",
	}

	UnprocessWinEvent(event)

	if len(Snapshot.FuturePlayers) != 1 {
		t.Error("Future Players wasn't decremented off an unwin event ", len(Snapshot.FuturePlayers), Snapshot)
	}

	if Snapshot.Teams[0].Captain.Points != 99 {
		t.Error("Bid wasnt incremented properly for unwin")
	}

	if len(Snapshot.Teams[0].Players) != 0 {
		t.Error("Player wasn't updated out of old team.")
	}
}

func TestSkipHappy(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:      "next_player",
			LeagueId: 1,
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     89,
				},
			},
		},
		SkippedPlayers: []*DraftPlayer{},
		FuturePlayers: []*DraftPlayer{
			&DraftPlayer{
				Ign:        "woah",
				HighestBid: 10,
				Team:       "test",
				LeagueId:   10,
			},
		},
	}

	event := SkipEvent{
		PlayerLeagueID: 1,
	}

	ProcessSkipEvent(event)

	if len(Snapshot.FuturePlayers) != 0 {
		t.Error("Future players was wrong should be 0: ", len(Snapshot.FuturePlayers))
	}

	if Snapshot.CurrentPlayer.LeagueId != 10 {
		t.Error("Current player wasnt changed")
	}

	if len(Snapshot.SkippedPlayers) != 1 {
		t.Error("Skipped Players not updated.")
	}
}

func TestSkipNoFuture(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: &DraftPlayer{
			Ign:      "next_player",
			LeagueId: 1,
		},
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     89,
				},
			},
		},
		SkippedPlayers: []*DraftPlayer{},
		FuturePlayers: []*DraftPlayer{
			&DraftPlayer{
				Ign:        "woah",
				HighestBid: 10,
				Team:       "test",
				LeagueId:   10,
			},
		},
	}

	event := SkipEvent{
		PlayerLeagueID: 1,
	}

	ProcessSkipEvent(event)

	if Snapshot.CurrentPlayer.LeagueId != 10 {
		t.Error("Current player wasnt changed")
	}

	if len(Snapshot.SkippedPlayers) != 1 {
		t.Error("Skipped Players not updated.")
	}
}

func TestSkipNoCurrent(t *testing.T) {
	Snapshot = LeagueDraftSnapshot{
		CurrentPlayer: nil,
		Teams: []*DraftTeam{
			&DraftTeam{
				Name: "test",
				Captain: DraftCaptain{
					FacebookID: "test_id",
					Points:     89,
				},
			},
		},
		SkippedPlayers: []*DraftPlayer{},
		FuturePlayers:  []*DraftPlayer{},
	}

	event := SkipEvent{
		PlayerLeagueID: 1,
	}

	ProcessSkipEvent(event)

	if len(Snapshot.SkippedPlayers) != 0 {
		t.Error("Uhhh")
	}
}
