package main

import "honnef.co/go/js/dom"

func h(tag string, className string) dom.Element {
	el := dom.GetWindow().Document().CreateElement(tag)
	if className != "" {
		el.Class().Add(className)
	}
	return el
}

func text(text string) dom.Node {
	return dom.GetWindow().Document().CreateTextNode(text)
}
