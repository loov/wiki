package main

import (
	"honnef.co/go/js/dom"

	"github.com/loov/wiki/fed"
	"github.com/loov/wiki/h"
)

type Lineup struct {
	Node dom.Element
	List []*fed.Stage
}

func NewLineup() *Lineup {
	lineup := &Lineup{}
	lineup.Node = h.Div("lineup")
	return lineup
}

func (lineup *Lineup) Open(title, url string) {
	stage := fed.NewStage(title, url)
	lineup.Node.AppendChild(stage.Node)
}
