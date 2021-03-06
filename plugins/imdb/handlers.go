// Package imdb implements a plugin to search IMDb titles.
package imdb

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/transport"
	imdb "github.com/StalkR/imdb/appengine"
)

func search(e *bot.Event) {
	q := strings.TrimSpace(e.Args)
	if len(q) == 0 {
		return
	}
	c, err := transport.Client(imdb.AppURL)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	titles, err := imdb.SearchTitle(c, q)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(titles) == 0 {
		e.Bot.Privmsg(e.Target, "No results found.")
		return
	}
	title, err := imdb.NewTitle(c, titles[0].ID)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, title.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("imdb", bot.Command{
		Help:    "imdb <title> - search a Title on IMDb",
		Handler: search,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
