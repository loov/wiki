package main

import (
	"honnef.co/go/js/dom"

	"github.com/loov/wiki/client"
	"github.com/loov/wiki/client/fedwiki"
	"github.com/loov/wiki/client/mark"
)

func main() {
	_ = fedwiki.Server{}
	_ = mark.Server{}

	cl := client.New()
	cl.Lineup.Servers[""] = mark.NewServer("/data/")
	cl.Lineup.Open("", "Welcome", "welcome")

	dom.GetWindow().
		Document().
		DocumentElement().
		AppendChild(cl.Node)
}
