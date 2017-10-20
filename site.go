package main

import (
	"github.com/loov/wiki/h"
	"honnef.co/go/js/dom"
)

type Site struct {
	Node   dom.Element
	Search *Search
	Lineup *Lineup
}

func NewSite() *Site {
	site := &Site{}
	site.Search = NewSearch()
	site.Lineup = NewLineup()
	site.Node = h.Div("app",
		h.Div("header",
			site.Search.Node,
		),
		site.Lineup.Node,
	)
	return site
}

type Search struct {
	Node dom.Element
}

func NewSearch() *Search {
	search := &Search{}
	search.Node = h.Form("search",
		h.Input("search-input"),
		h.Button("", h.Text("Search")))
	return search
}
