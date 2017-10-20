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
	site.Node = h.Tag("div", "app",
		h.Tag("div", "header",
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
	search.Node = h.Tag("form", "search",
		h.Tag("input", "search-input"),
		h.Tag("button", "", h.Text("Search")))
	return search
}
