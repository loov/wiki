package main

import (
	"github.com/loov/wiki/client"
	"honnef.co/go/js/dom"
)

func main() {
	cl := client.New()
	cl.Lineup.Open("fedwiki", "Welcome", "/data/welcome.json")

	dom.GetWindow().
		Document().
		DocumentElement().
		AppendChild(cl.Node)
}
