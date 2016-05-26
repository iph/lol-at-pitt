package draft

import (
	"sync"

	"encoding/csv"
	"github.com/TrevorSStone/goriot"
	"github.com/lab-D8/lol-at-pitt/ols"
	"os"
	"strconv"
)

var (
	Snapshot LeagueDraftSnapshot
	Events   DraftStore
	OlsId    string
	Paused   bool       = true
	lock     sync.Mutex = sync.Mutex{}
)

func Init(olsId string) {
	Events = DraftStore{
		DraftEvents: []*DraftEvent{},
		ID:          olsId,
	}
	OlsId = olsId
	db := ols.InitDB()
	collection := db.C("auctions")
	mapa := map[string]string{
		"id": olsId,
	}
	data := collection.Find(mapa)
	amt, err := data.Count()

	if err != nil {
		panic("Things blew up in the database")
	}

	if amt == 0 {
		Snapshot = LeagueDraftSnapshot{
			CurrentPlayer:  nil,
			FuturePlayers:  []*DraftPlayer{},
			Teams:          InitInitialCaptains("resources/captains.csv"),
			SkippedPlayers: []*DraftPlayer{},
			EventID:        -1,
			ID:             olsId,
		}
		collection.Insert(Snapshot)
	} else {
		data.One(&Snapshot)
	}
}

func InitInitialCaptains(filename string) []*DraftTeam {
	r, err := os.Open(filename)

	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(r)
	allData, err := csvReader.ReadAll()

	if err != nil {
		panic(err)
	}
	teams := []*DraftTeam{}
	for _, record := range allData[1:] {

		normalizedSummonerName := goriot.NormalizeSummonerName(record[1])[0]
		summ, _ := goriot.SummonerByName(goriot.NA, normalizedSummonerName)
		points, _ := strconv.Atoi(record[3])
		team := &DraftTeam{
			Players: []*DraftPlayer{},
			Name:    record[0] + "'s Team",
			Captain: DraftCaptain{
				Name:          record[0],
				LeagueID:      summ[normalizedSummonerName].ID,
				Points:        points,
				NormalizedIgn: normalizedSummonerName,
			},
		}

		teams = append(teams, team)
	}

	return teams
}

func Save() {
	db := ols.InitDB()
	collection := db.C("auctions")
	mapa := map[string]string{
		"id": OlsId,
	}
	collection.Update(mapa, Snapshot)
	collection.Update(mapa, Events)
}

func ProcessBidEvent(event BidEvent) bool {
	team := Snapshot.GetTeamByFacebookId(event.DrafterFacebookID)

	if team == nil {
		return false
	}

	lock.Lock()
	currentHighestBid := Snapshot.CurrentPlayer.HighestBid
	currentHighestTeam := Snapshot.CurrentPlayer.Team
	bidSuccessful := false
	// Don't bid if you are the current highest bidder, if you havent bid MORE than the current one, or you don't have the points.
	if team.Name == currentHighestTeam || currentHighestBid >= event.Amount || event.Amount > team.Captain.Points {
		bidSuccessful = false
	} else {
		Snapshot.CurrentPlayer.HighestBid = event.Amount
		Snapshot.CurrentPlayer.Team = team.Name
		bidSuccessful = true
	}
	lock.Unlock()
	return bidSuccessful

}

func ProcessWinEvent(event WinEvent) {
	lock.Lock()
	team := Snapshot.GetTeamByFacebookId(event.CaptainFacebookID)
	if team != nil {
		team.Captain.Points -= event.Amount
		team.Players = append(team.Players, Snapshot.CurrentPlayer)
	}

	Paused = true

	if len(Snapshot.FuturePlayers) != 0 {
		Snapshot.CurrentPlayer = Snapshot.FuturePlayers[0]
		Snapshot.FuturePlayers = Snapshot.FuturePlayers[1:]
	} else {
		Snapshot.CurrentPlayer = nil
	}
	lock.Unlock()
}

func UnprocessWinEvent(event WinEvent) {
	lock.Lock()

	if Snapshot.CurrentPlayer != nil {
		Snapshot.FuturePlayers = append([]*DraftPlayer{Snapshot.CurrentPlayer}, Snapshot.FuturePlayers...)
	}

	team := Snapshot.GetTeamByFacebookId(event.CaptainFacebookID)
	if team != nil {
		team.Captain.Points += event.Amount
		for i, player := range team.Players {
			if player.LeagueId == event.PlayerLeagueID {
				currentPlayer := team.Players[i]
				Snapshot.CurrentPlayer = currentPlayer
				team.Players = append(team.Players[0:i], team.Players[i+1:]...)
			}
		}
	}

	Paused = true
	lock.Unlock()
}

func TogglePause() {
	Paused = !Paused
}

func ProcessSkipEvent(event SkipEvent) {
	lock.Lock()
	if Snapshot.CurrentPlayer != nil {
		Snapshot.SkippedPlayers = append(Snapshot.SkippedPlayers, Snapshot.CurrentPlayer)
		Snapshot.CurrentPlayer = nil
	}

	if len(Snapshot.FuturePlayers) > 0 {
		Snapshot.CurrentPlayer = Snapshot.FuturePlayers[0]
		Snapshot.FuturePlayers = Snapshot.FuturePlayers[1:]
	}

	Paused = true
	lock.Unlock()
}

func UnprocessSkipEvent(event SkipEvent) {
	lock.Lock()
	if Snapshot.CurrentPlayer != nil {
		Snapshot.FuturePlayers = append([]*DraftPlayer{Snapshot.CurrentPlayer}, Snapshot.FuturePlayers...)
		Snapshot.CurrentPlayer = nil
	}
	for i, player := range Snapshot.SkippedPlayers {
		if player.LeagueId == event.PlayerLeagueID {
			Snapshot.CurrentPlayer = Snapshot.SkippedPlayers[i]
			Snapshot.SkippedPlayers = append(Snapshot.SkippedPlayers[0:i], Snapshot.SkippedPlayers[i+1:]...)
		}
	}

	Paused = true
	lock.Unlock()
}

func ProcessAddTeamEvent(event AddTeamEvent) {
	lock.Lock()

	lock.Unlock()
}
