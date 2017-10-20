package main

import (
	"honnef.co/go/js/dom"
)

func main() {
	wiki := &Wiki{}
	wiki.Pages = append(
		wiki.Pages,
		&Welcome,
		&Second,
	)

	site := &Site{}
	site.Doc = dom.GetWindow().Document()
	site.Node = site.Doc.GetElementByID("app")
	site.Init()

	site.Wiki = wiki
	site.UpdateStages()
}
