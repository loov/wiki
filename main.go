package main

import (
	"github.com/loov/wiki/client"
	"github.com/loov/wiki/client/fedwiki"
	"honnef.co/go/js/dom"
)

func main() {
	cl := client.New()

	cl.Lineup.Providers["fedwiki"] = &fedwiki.Provider{}
	cl.Lineup.Providers[""] = cl.Lineup.Providers["fedwiki"]

	cl.Lineup.Open("fedwiki", "Welcome", "/data/welcome.json")

	dom.GetWindow().
		Document().
		DocumentElement().
		AppendChild(cl.Node)
}
