package main

import (
	"github.com/TrevorSStone/goriot"
	"github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
	"log"
)

const (
	BLUE_TEAM = 100
	RED_TEAM  = 200
)

type MatchRule func(player ols.Player, match goriot.Game) bool

func CheckGames() {
	players := db.GetPlayersDAO().All()
	matchRules := []MatchRule{CorrectPlayerTeam, CorrectOtherTeam, AlreadyChecked, CorrectGameType, CorrectGameMode}
	var usedSet map[int64]bool
	for _, player := range players {
		summonerId := player.Id
		games, err := goriot.RecentGameBySummoner("na", summonerId)
		if err != nil {
			log.Printf("Error: ", err.Error())
		}

		for _, game := range games {
			_, ok := usedSet[game.GameID]
			if !ok {
				continue
			}

			allowedGame := ApplyRules(*player, game, matchRules)
			if allowedGame {
				usedSet[game.GameID] = true
				createMatch(*player, game)
			}

		}
	}

}
func ApplyRules(player ols.Player, match goriot.Game, rules []MatchRule) bool {
	allowedGame := true
	for _, rule := range rules {
		allowedGame = allowedGame && rule(player, match)
		if !allowedGame {
			break
		}
	}

	return allowedGame
}

func CorrectPlayerTeam(player ols.Player, match goriot.Game) bool {
	myTeamMatch := 0

	playingTeam := db.GetTeamsDAO().LoadPlayer(player.Id)
	myTeam := match.TeamID

	for _, fellowPlayer := range match.FellowPlayers {
		if fellowPlayer.TeamID == myTeam && playingTeam.IsPlayerOnTeam(fellowPlayer.SummonerID) {
			myTeamMatch++
		}
	}

	return myTeamMatch >= 3
}

func CorrectOtherTeam(player ols.Player, match goriot.Game) bool {
	// Annoying.
	otherTeamId := BLUE_TEAM
	matchAmount := 0
	myTeam := match.TeamID

	if myTeam == otherTeamId {
		otherTeamId = RED_TEAM
	}

	otherTeam := getOtherTeam(player, match)
	for _, fellowPlayer := range match.FellowPlayers {
		if fellowPlayer.TeamID == otherTeamId && otherTeam.IsPlayerOnTeam(fellowPlayer.SummonerID) {
			matchAmount++
		}
	}

	return matchAmount >= 3

}

func getOtherTeam(player ols.Player, match goriot.Game) ols.Team {
	var otherTeam ols.Team
	otherTeamId := BLUE_TEAM
	myTeam := match.TeamID
	if myTeam == otherTeamId {
		otherTeamId = RED_TEAM
	}

	// Get other player on other team..
	for _, fellowPlayer := range match.FellowPlayers {
		if fellowPlayer.TeamID != myTeam {
			otherTeam = db.GetTeamsDAO().LoadPlayer(fellowPlayer.SummonerID)
			if otherTeam.Name != "" {
				break
			}
		}
	}
	return otherTeam
}

// Returns false if it is saved.
func AlreadyChecked(player ols.Player, match goriot.Game) bool {
	return !db.GetMatchesDAO().IsSaved(match.GameID)

}

func CorrectGameType(player ols.Player, match goriot.Game) bool {
	return match.GameType == "CUSTOM_GAME" && match.SubType == "NONE"
}

func CorrectGameMode(player ols.Player, match goriot.Game) bool {
	return match.GameMode == "CLASSIC"
}

// Hell function that turns bullshit into magic
func createMatch(player ols.Player, game goriot.Game) {
	match, err := goriot.MatchByMatchID("na", true, game.GameID)
	var participants map[int]*ols.Participant // championid -> participant
	participants[game.ChampionID] = &ols.Participant{Id: player.Id}
	if err != nil {
		log.Println("Match with id: ", game.GameID, " had an error:", err.Error())
		return
	}

	// Connect participants of an anonymous game to one of a recent game...
	for _, fellowPlayer := range game.FellowPlayers {
		participants[fellowPlayer.ChampionID] = &ols.Participant{Id: fellowPlayer.SummonerID}
	}

	// All info is connected now!
	for _, matchPlayer := range match.Participants {
		participant := participants[matchPlayer.ChampionID]
		participant.ParticipantId = matchPlayer.ParticipantID
	}

	var matchParticipants []ols.Participant
	for _, participant := range participants {
		matchParticipants = append(matchParticipants, *participant)
	}
	///////////////////////
	blueTeam := getTeamName(game, BLUE_TEAM)
	redTeam := getTeamName(game, RED_TEAM)
	winnerTeam := blueTeam

	if game.Statistics.Win && game.TeamID == RED_TEAM {
		winnerTeam = redTeam
	}
	week := db.GetMatchesDAO().LoadWeekForMatch(blueTeam, redTeam)
	olsMatch := ols.Match{
		Participants: matchParticipants,
		BlueTeam:     blueTeam,
		RedTeam:      redTeam,
		Played:       true,
		Week:         week,
		Winner:       winnerTeam,
		Id:           game.GameID,
	}

	db.GetMatchesDAO().Save(olsMatch)
}

func getTeamName(game goriot.Game, teamCode int) string {
	for _, fellowPlayer := range game.FellowPlayers {
		if teamCode == fellowPlayer.TeamID {
			team := db.GetTeamsDAO().LoadPlayer(fellowPlayer.SummonerID)
			if team.Name != "" {
				return team.Name
			}
		}
	}

	return ""
}
