package h

import (
	"honnef.co/go/js/dom"
)

func Tag(tag string, className string, children ...dom.Node) dom.Element {
	el := dom.GetWindow().Document().CreateElement(tag)
	if className != "" {
		el.Class().Add(className)
	}
	for _, child := range children {
		el.AppendChild(child)
	}
	return el
}

func Div(className string, children ...dom.Node) dom.Element {
	return Tag("div", className, children...)
}

func H1(className string, children ...dom.Node) dom.Element {
	return Tag("h1", className, children...)
}
func H2(className string, children ...dom.Node) dom.Element {
	return Tag("h2", className, children...)
}
func H3(className string, children ...dom.Node) dom.Element {
	return Tag("h3", className, children...)
}

func P(children ...string) dom.Element {
	el := dom.GetWindow().Document().CreateElement("p")
	for _, child := range children {
		el.AppendChild(Text(child))
	}
	return el
}

func A(className string, href string, children ...dom.Node) dom.Element {
	el := Tag("a", className, children...)
	el.SetAttribute("href", href)
	return el
}

func Form(className string, children ...dom.Node) dom.Element {
	return Tag("form", className, children...)
}

func Input(className string, children ...dom.Node) dom.Element {
	return Tag("input", className, children...)
}

func Button(className string, children ...dom.Node) dom.Element {
	return Tag("button", className, children...)
}

func Text(text string) dom.Node {
	return dom.GetWindow().Document().CreateTextNode(text)
}

func Fragment(children ...dom.Node) dom.DocumentFragment {
	el := dom.GetWindow().Document().CreateDocumentFragment()
	for _, child := range children {
		el.AppendChild(child)
	}
	return el
}
