package draft

import (
	"log"
)

type LeagueDraftSnapshot struct {
	CurrentPlayer  *DraftPlayer
	FuturePlayers  []*DraftPlayer
	SkippedPlayers []*DraftPlayer
	Teams          []*DraftTeam
	EventID        int
	ID             string
}

func (ld *LeagueDraftSnapshot) GetTeamByFacebookId(id string) *DraftTeam {
	for _, val := range ld.Teams {
		if val.Captain.FacebookID == id {
			return val
		}
	}
	return nil
}

type DraftTeam struct {
	Captain DraftCaptain
	Players []*DraftPlayer
	Name    string
}

type DraftPlayer struct {
	Ign        string
	HighestBid int
	LeagueId   int64
	RealName   string
	Team       string
	Tier       string
	Roles      []LeagueRole
	Score      int
}

type DraftCaptain struct {
	Name          string
	FacebookID    string
	LeagueID      int64
	Points        int
	NormalizedIgn string // used for easy lookup
}

type LeagueRole struct {
	Position string
	Comments string
	Score    int
}

type DraftPlayers []*DraftPlayer
type DraftCaptains []*DraftCaptain

func (p *DraftCaptains) Print() {
	for _, player := range *p {
		log.Println(player)
	}
}

// Sorting fun-ness
func (p DraftCaptains) Len() int {
	return len(p)
}

func (p DraftCaptains) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p DraftCaptains) Less(i, j int) bool {
	return p[i].Points > p[j].Points
}
