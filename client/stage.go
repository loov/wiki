package client

import (
	"honnef.co/go/js/dom"

	"github.com/loov/wiki/h"
)

type Stage struct {
	Lineup *Lineup
	View   View

	Node    dom.Element
	Slug    dom.Element
	Buttons dom.Element
	Content dom.Element
}

func NewStage(lineup *Lineup, view View) *Stage {
	stage := &Stage{}
	stage.Lineup = lineup
	stage.View = view

	stage.Content = h.Div("content")
	stage.Slug = h.Div("slug")
	stage.Buttons = h.Div("buttons")
	stage.Node = h.Div("stage",
		h.Div("indicator"),
		h.Div("status",
			stage.Slug,
			stage.Buttons,
		),
		stage.Content,
	)

	h.AttachOverflowIndicator(stage.Content)
	stage.View.Attach(stage)

	return stage
}

func (stage *Stage) SetTag(class string, state bool) {
	if state {
		stage.Node.Class().Add(class)
	} else {
		stage.Node.Class().Remove(class)
	}
}

func (stage *Stage) SetSlug(node dom.Node) {
	stage.Slug.SetInnerHTML("")
	stage.Slug.AppendChild(node)
}

func (stage *Stage) SetButtons(node dom.Node) {
	stage.Buttons.SetInnerHTML("")
	stage.Buttons.AppendChild(node)
}

func (stage *Stage) SetContent(node dom.Node) {
	stage.Content.SetInnerHTML("")
	stage.Content.AppendChild(node)
}

func (stage *Stage) OpenNext(view View) {
	// TODO: is this the best place for it?
	next := NewStage(stage.Lineup, view)

	stage.Lineup.CloseTrailing(stage)
	stage.Lineup.Add(next)

	h.ScrollIntoView(next.Node)
}

func (stage *Stage) OpenLast(view View) {
	// TODO: is this the best place for it?
	next := NewStage(stage.Lineup, view)

	stage.Lineup.Add(next)

	h.ScrollIntoView(next.Node)
}

func (stage *Stage) Close() {
	stage.View.Detach()
}
