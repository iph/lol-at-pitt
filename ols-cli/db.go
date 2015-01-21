package main

import (
	"encoding/json"
	"github.com/lab-d8/lol-at-pitt/ols"
	"io/ioutil"
	"labix.org/v2/mgo"
	"os"
)

func dumpDb(filename string) {
	session, _ := mgo.Dial(MongoLocation)
	db := session.DB(DatabaseName)
	json_blob := map[string]interface{}{}
	var teams ols.Teams
	db.C("teams").Find(map[string]string{}).All(&teams)
	json_blob["Teams"] = teams

	var players ols.Players
	db.C("players").Find(map[string]string{}).All(&players)
	db.CollectionNames()
	json_blob["Players"] = players

	data, _ := json.MarshalIndent(json_blob, "", "  ")
	ioutil.WriteFile(filename, data, 0644)
}

func initDbPlayers(players ols.Players) {
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.C("players").DropCollection()
	for _, player := range players {
		db.C("players").Insert(player)
	}
	session.Close()

}

func initDbTeams(teams ols.Teams) {
	session, err := mgo.Dial(MongoLocation)
	defer session.Close()

	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.C("teams").DropCollection()

	for _, team := range teams {
		db.C("teams").Insert(team)
	}
}

func upload(json_file string) {
	var db_blob DB
	file, _ := os.Open(json_file)
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	json.Unmarshal(data, &db_blob)
	players := db_blob.Players
	initDbPlayers(players)

	teams := db_blob.Teams
	initDbTeams(teams)
}

func deleteDb() {
	session, err := mgo.Dial(MongoLocation)
	defer session.Close()

	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.DropDatabase()
}
