package main

import (
	"honnef.co/go/js/dom"
)

type Site struct {
	Node   dom.Element
	Search *Search
	Stages *Stages

	Wiki *Wiki
}

func NewSite() *Site {
	site := &Site{}
	site.Search = NewSearch()
	site.Stages = NewStages()
	site.Node = h("div", "app",
		h("div", "header",
			site.Search.Node,
		),
		site.Stages.Node,
	)
	return site
}

type Search struct {
	Node dom.Element
}

func NewSearch() *Search {
	search := &Search{}
	search.Node = h("form", "search",
		h("input", "search-input"),
		h("button", "", text("Search")))
	return search
}
