package client

import (
	"honnef.co/go/js/dom"

	"github.com/loov/wiki/h"
)

type Lineup struct {
	Node dom.Element
	List []*Stage
}

func NewLineup() *Lineup {
	lineup := &Lineup{}
	lineup.Node = h.Div("lineup")
	return lineup
}

func (lineup *Lineup) Open(title, url string) {
	lineup.Add(NewStage(lineup, title, url))
}

func (lineup *Lineup) indexOf(stage *Stage) int {
	for i, x := range lineup.List {
		if x == stage {
			return i
		}
	}
	return -1
}

func (lineup *Lineup) CloseTrailing(target *Stage) {
	i := lineup.indexOf(target)
	if i < 0 {
		return
	}

	for k := len(lineup.List) - 1; k > i; k-- {
		lineup.closeIndex(k)
	}
}

func (lineup *Lineup) closeIndex(i int) {
	// TODO: don't close stages that are being edited
	stage := lineup.List[i]
	lineup.List = append(lineup.List[:i], lineup.List[i+1:]...)
	lineup.Node.RemoveChild(stage.Node)
	stage.Close()
}

func (lineup *Lineup) Add(stage *Stage) {
	lineup.List = append(lineup.List, stage)
	lineup.Node.AppendChild(stage.Node)
}
