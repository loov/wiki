package main

import (
	"honnef.co/go/js/dom"

	"github.com/loov/wiki/client"
	"github.com/loov/wiki/client/fedwiki"
)

func main() {
	cl := client.New()
	cl.Lineup.Servers[""] = fedwiki.NewServer("")
	cl.Lineup.Open("", "Welcome", "welcome")

	dom.GetWindow().
		Document().
		DocumentElement().
		AppendChild(cl.Node)
}
