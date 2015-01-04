package main

import (
	"github.com/go-martini/martini"
	dao "github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/draft"
	"github.com/lab-d8/lol-at-pitt/site"
	"github.com/martini-contrib/render"
	"log"
	"net/url"
	"strconv"
)

func initDraftRouter(m *martini.ClassicMartini) {
	m.Get("/draft/pause", func(renderer render.Render, draft *draft.Draft) {
		draft.Pause()
		renderer.Data(200, []byte("success"))
	})

	m.Get("/draft/resume", func(renderer render.Render, draft *draft.Draft) {
		draft.Resume()
		renderer.Data(200, []byte("success"))
	})

	m.Get("/draft/start", func(renderer render.Render, draft *draft.Draft) {
		draft.Start()
		renderer.Data(200, []byte("success"))
	})

	m.Get("draft/next", func(renderer render.Render, draft *draft.Draft) {
		draft.Next()
		renderer.Data(200, []byte("success"))
	})

	m.Get("/draft/bid", CaptainRequired, func(urls url.Values, renderer render.Render, d *draft.Draft, user site.User) {
		player := dao.GetPlayersDAO().Load(user.LeagueId)
		bidAmount, err := strconv.Atoi(urls.Get("amount"))
		if err != nil {
			log.Println("error bidding: ", err.Error())
		} else {
			d.Bid(bidAmount, player.Team)
		}

	})
	m.Get("/draft/history", func(renderer render.Render, d *draft.Draft) {
		renderer.JSON(200, d.History.Values)
	})

	m.Get("/draft/status", CaptainRequired, func(renderer render.Render, d *draft.Draft) {
		renderer.JSON(200, d)
	})

}
