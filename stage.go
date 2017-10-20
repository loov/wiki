package main

import "honnef.co/go/js/dom"

type Stage struct {
	Stages *Stages
	Node   dom.Element

	URL     string
	Loading bool
	Editing bool

	Page *Page
}
