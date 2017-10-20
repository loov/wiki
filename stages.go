package main

import "honnef.co/go/js/dom"

type Stages struct {
	Node dom.Element
	List []*Stage
}

func NewStages(parent dom.Element) *Stages {
	stages := &Stages{}
	stages.Node = h("div", "stages")
	return stages
}

func (stages *Stages) Open(title, url string) {

}
