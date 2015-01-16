package db

import (
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
	"sort"
)

var MatchesCollectionName string = "olsmatches"

type MatchesDAO struct {
	DAO
}

func NewMatchesContext(db *mgo.Database) *MatchesDAO {
	dao := MatchesDAO{DAO{db, db.C(MatchesCollectionName)}}
	return &dao
}

func (m *MatchesDAO) Load(id int64) ols.Match {
	var match ols.Match
	m.Collection.Find(map[string]int64{"id": id}).One(&match)
	return match
}

func (m *MatchesDAO) Save(match ols.Match) {
	m.DAO.Save(map[string]interface{}{"id": match.Id, "week": match.Week, "blueteam": match.BlueTeam, "redteam": match.RedTeam}, match)
}

func (m *MatchesDAO) LoadTeamMatches(team string) []*ols.Match {
	var matches []*ols.Match
	var matchesRed []*ols.Match
	m.Collection.Find(map[string]string{"blueteam": team}).All(&matches)
	m.Collection.Find(map[string]string{"redteam": team}).All(&matchesRed)

	allMatches := append(matches, matchesRed...)
	sort.Sort(ols.Matches(allMatches))
	return allMatches
}

func (m *MatchesDAO) LoadWinningMatches(team string) []*ols.Match {
	var matches []*ols.Match
	m.Collection.Find(map[string]string{"winner": team}).All(&matches)

	sort.Sort(ols.Matches(matches))
	return matches
}

func (m *MatchesDAO) Delete(match ols.Match) {
	m.Collection.Remove(match)
}