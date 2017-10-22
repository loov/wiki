package h

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

func ScrollIntoView(el dom.Element) {
	el.Underlying().Call("scrollIntoView", js.M{
		"behavior": "smooth",
		"block":    "center",
		"inline":   "center",
	})
}
