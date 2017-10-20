package main

import (
	"honnef.co/go/js/dom"
)

func main() {
	site := NewSite()
	site.Lineup.Open("Welcome", "welcome")
	dom.GetWindow().Document().DocumentElement().AppendChild(site.Node)
}
