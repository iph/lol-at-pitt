package ols

import (
	"github.com/lab-D8/lol-at-pitt/site"
	"gopkg.in/mgo.v2"
)

type UsersDAO struct {
	DAO
}

var UsersCollectionName string = "users"

func (u *UsersDAO) GetUserFB(facebookId string) site.User {
	var user site.User
	u.Collection.Find(map[string]string{"facebookid": facebookId}).One(&user)
	return user
}

func (u *UsersDAO) GetUserLeague(leagueId int64) site.User {
	var user site.User
	u.Collection.Find(map[string]int64{"leagueid": leagueId}).One(&user)
	return user
}

func (u *UsersDAO) GetUserByIgn(ign string) site.User {
	var user site.User
	player := GetPlayersDAO().LoadNormalizedIGN(ign)
	u.Collection.Find(map[string]int64{"leagueid": player.Id}).One(&user)
	return user
}

func (u *UsersDAO) Save(user site.User) {
	u.DAO.Save(map[string]string{"facebookid": user.FacebookId}, user)
}

func NewUserDAO(db *mgo.Database) *UsersDAO {
	d := DAO{db: db, Collection: db.C(UsersCollectionName)}
	return &UsersDAO{d}
}
