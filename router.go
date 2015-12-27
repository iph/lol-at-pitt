package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/TrevorSStone/goriot"
	"github.com/go-martini/martini"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/lab-d8/lol-at-pitt/site"
	"github.com/lab-d8/oauth2"
	"github.com/martini-contrib/render"
)

// Register is used for derp
type Register struct {
	Id   string
	Name string
	Next string
}

func main() {
	m := martini.Classic()
	goriot.SetAPIKey(LeagueApiKey)
	goriot.SetLongRateLimit(LongLeagueLimit, 10*time.Minute)
	goriot.SetSmallRateLimit(ShortLeagueLimit, 10*time.Second)
	// Setup middleware to be attached to the controllers on every call.
	if Debug {
		InitDebugMiddleware(m)
	} else {
		InitMiddleware(m)
	}

	teamHandler := func(renderer render.Render) {
		teams := ols.GetTeamsDAO().All()
		renderer.JSON(200, teams)
	}

	individualTeamHandler := func(params martini.Params, renderer render.Render) {
		team := ols.GetTeamsDAO().Load(params["name"])
		renderer.JSON(200, team)
	}

	m.Get("/error", func(urls url.Values, renderer render.Render) {
		renderer.HTML(200, "error", urls.Get("status"))
	})
	m.Get("/players/:name", func(params martini.Params, renderer render.Render) {
		normalizedName := params["name"]
		fmt.Println(normalizedName)
		player := ols.GetPlayersDAO().LoadNormalizedIGN(normalizedName)
		fmt.Println(player)
		displayTeam := ols.GetTeamsDAO().LoadPlayerDisplay(player.Id)
		fmt.Println(displayTeam)
		renderer.HTML(200, "team", displayTeam)
	})

	m.Get("/results", func(renderer render.Render) {
		players := ols.GetPlayersDAO().All()
		allPlayerDisplay := []ols.PlayerDraftResult{}
		for _, player := range players {
			team := ols.GetTeamsDAO().LoadPlayer(player.Id)
			playerDisplay := ols.PlayerDraftResult{Ign: player.Ign, Team: team.Name, NormalizedIgn: player.NormalizedIgn}
			allPlayerDisplay = append(allPlayerDisplay, playerDisplay)
		}

		renderer.HTML(200, "drafted", allPlayerDisplay)
	})
	m.Get("/teams", teamHandler)
	m.Get("/team/:name", individualTeamHandler)
	m.Get("/", func(renderer render.Render) {
		renderer.HTML(200, "main", 1)
	})
	m.Get("/register", LoginRequired, func(urls url.Values, renderer render.Render) {
		renderer.HTML(200, "register", Register{Next: urls.Get("next")})
	})

	m.Get("/oauth2error", func(token oauth2.Tokens, renderer render.Render) {
		renderer.JSON(200, token)
	})

	m.Get("/register/complete", LoginRequired, func(urls url.Values, renderer render.Render, token oauth2.Tokens, w http.ResponseWriter, r *http.Request) {
		summonerName := urls.Get("summoner")
		teamName := urls.Get("team")

		if token.Expired() {
			http.Redirect(w, r, "/error?status=InvalidFacebook", 302)
			return
		}

		id, err := GetId(token.Access())
		if err != nil {
			renderer.Status(404)
			return
		}

		normalizedSummonerName := goriot.NormalizeSummonerName(summonerName)[0]
		player := ols.GetPlayersDAO().LoadNormalizedIGN(normalizedSummonerName)
		if player.Id == 0 {
			http.Redirect(w, r, "/error?status=NoPlayerFound", 302)
		}

		user := ols.GetUserDAO().GetUserFB(id)

		// User is registered registered
		if user.LeagueId != 0 {
			http.Redirect(w, r, "/error?status=AlreadyRegistered", 302)
			return
		}

		user = site.User{LeagueId: player.Id, FacebookId: id}
		log.Println("User registered:", user)
		if player.Id == 0 {
			// new player not in our db
			ols.GetPlayersDAO().Save(player)
		}
		team := ols.GetTeamsDAO().LoadPlayerByCaptain(player.Id)
		newTeam := team
		if team.Name != "" && teamName != "" {
			newTeam.Name = teamName
			ols.GetTeamsDAO().Update(team, newTeam)

		}
		ols.GetUserDAO().Save(user)
		//next := urls.Get("next")
		renderer.HTML(200, "register_complete", 1)
	})

	err := http.ListenAndServe(":6060", m) // Nginx needs to redirect here, so we don't need sudo priv to test.
	if err != nil {
		log.Println(err)
	}

}
