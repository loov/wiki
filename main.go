package main

import (
	"honnef.co/go/js/dom"

	"github.com/loov/wiki/client"
	"github.com/loov/wiki/client/fedwiki"
	"github.com/loov/wiki/client/mark"
)

func main() {
	_ = fedwiki.Context{}
	_ = mark.Context{}

	cl := client.New()
	cl.Lineup.Contexts[""] = mark.NewContext("/data/")
	cl.Lineup.Open("", "Welcome", "welcome.md")

	dom.GetWindow().
		Document().
		DocumentElement().
		AppendChild(cl.Node)
}
