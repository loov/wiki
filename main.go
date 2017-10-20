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
	site.Node = dom.GetWindow().Document().GetElementByID("app")
	site.Init()

	site.Wiki = wiki
	site.UpdateStages()
}
