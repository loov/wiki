package main

import "honnef.co/go/js/dom"

func h(tag string, className string, children ...dom.Node) dom.Element {
	el := dom.GetWindow().Document().CreateElement(tag)
	if className != "" {
		el.Class().Add(className)
	}
	for _, child := range children {
		el.AppendChild(child)
	}
	return el
}

func text(text string) dom.Node {
	return dom.GetWindow().Document().CreateTextNode(text)
}

func frag(children ...dom.Node) dom.DocumentFragment {
	el := dom.GetWindow().Document().CreateDocumentFragment()
	for _, child := range children {
		el.AppendChild(child)
	}
	return el
}
