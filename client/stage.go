package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	go stage.fetch()

	return stage
}

func (stage *Stage) Close() {

}

func (stage *Stage) Update() {
	stage.PageNode.SetInnerHTML("")
	if stage.Page == nil || stage.Loading {
		stage.Node.Class().Add("loading")
	} else {
		stage.Node.Class().Remove("loading")
		stage.PageNode.AppendChild(h.Fragment(
			h.Div("title", h.Text(stage.Page.Title)),
			h.Div("story", stage.RenderAll(stage.Page.Story...)...),
		))
	}
}

func (stage *Stage) fetch() {
	// TODO: proper threading
	r, err := http.Get("/data/welcome.json")
	if err != nil {
		log.Println(err)
		return
	}

	page := &Page{}
	err = json.NewDecoder(r.Body).Decode(page)
	if err != nil {
		log.Println(err)
		return
	}

	stage.Loading = false
	stage.Page = page
	stage.Update()
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

	switch item.Type() {
	case "paragraph":
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
		}).Run(item.String("text"))
	default:
		el.Class().Add("missing")
		el.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return el
}

func (stage *Stage) LinkClicked(ev dom.Event) {
	target := ev.Target()
	ev.StopPropagation()
	ev.PreventDefault()

	stage.Lineup.CloseTrailing(stage)

	next := NewStage(stage.Lineup, target.TextContent(), target.GetAttribute("href"))
	stage.Lineup.Add(next)
	h.ScrollIntoView(next.Node)
}
