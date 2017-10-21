package client

import (
	"fmt"

	"honnef.co/go/js/dom"

	"github.com/loov/wiki/h"
)

type Stage struct {
	Lineup *Lineup
	Node   dom.Element

	Title   string
	Context string
	URL     string
	Loading bool
	Editing bool

	PageNode dom.Element
	Page     *Page
}

func NewStage(lineup *Lineup, title, url string) *Stage {
	stage := &Stage{}
	stage.Lineup = lineup
	stage.Title = title
	stage.URL = url
	stage.Loading = true

	stage.PageNode = h.Div("page")
	h.AttachOverflowIndicator(stage.PageNode)

	stage.Node = h.Div("stage",
		h.Div("indicator"),
		h.Div("status",
			h.Div("url", h.Text("[loov.io] "+stage.URL)),
			h.Div("icon", h.Text("Edit")),
		),
		stage.PageNode,
	)
	stage.Update()

	// simulate fetch
	dom.GetWindow().SetTimeout(func() {
		stage.Loading = false
		stage.Page = &Welcome
		stage.Update()
	}, 3000)

	return stage
}

func (stage *Stage) Close() {

}

func (stage *Stage) Update() {
	stage.PageNode.SetInnerHTML("")
	if stage.Page == nil || stage.Loading {
		stage.Node.Class().Add("loading")
		stage.PageNode.AppendChild(h.Fragment(
			h.Div("title", h.Text(stage.Title)),
		))
	} else {
		stage.Node.Class().Remove("loading")
		stage.PageNode.AppendChild(h.Fragment(
			h.Div("title", h.Text(stage.Page.Title)),
			h.Div("story", stage.RenderAll(stage.Page.Story...)...),
		))
	}
}

func (stage *Stage) RenderAll(items ...Item) []dom.Node {
	rs := []dom.Node{}
	for _, item := range items {
		rs = append(rs, stage.Render(item))
	}
	return rs
}

func (stage *Stage) Render(item Item) dom.Element {
	el := h.Div("item")

	switch item := item.(type) {
	case *Paragraph:
		el.Class().Add("paragraph")
		p := h.P()
		(&Parser{
			Begin: func() { p = h.P() },
			End:   func() { el.AppendChild(p); p = nil },
			Text: func(s string) {
				p.AppendChild(h.Text(s))
			},
			Link: func(spec string) {
				link := h.A("", spec, h.Text(spec))
				link.AddEventListener("click", false, stage.LinkClicked)
				p.AppendChild(link)
			},
		}).Run(item.Text)
	default:
		el.Class().Add("missing")
		el.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return el
}

func (stage *Stage) LinkClicked(ev dom.Event) {
	stage.Lineup.CloseTrailing(stage)
	target := ev.Target()
	stage.Lineup.Add(NewStage(stage.Lineup, target.TextContent(), target.GetAttribute("href")))

	ev.StopPropagation()
	ev.PreventDefault()
}
