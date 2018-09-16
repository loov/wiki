package main

import (
	"honnef.co/go/js/dom"

	"github.com/loov/wiki/client"
	"github.com/loov/wiki/client/fedwiki"
	"github.com/loov/wiki/client/html"
	"github.com/loov/wiki/client/mark"
)

func main() {
	_ = fedwiki.Context{}
	_ = mark.Context{}
	_ = html.Context{}

	cl := client.New()
	cl.Lineup.Contexts["data"] = mark.NewContext("/data/")
	cl.Lineup.Open("data", "Welcome", "welcome.md")

	cl.Lineup.Contexts["neti.ee"] = html.NewContext("http://wikipedia.com")
	cl.Lineup.Open("neti.ee", "Home", "http://wikipedia.com")

	dom.GetWindow().
		Document().
		DocumentElement().
		AppendChild(cl.Node)
}
