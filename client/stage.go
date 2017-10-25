package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"honnef.co/go/js/dom"

	"github.com/loov/wiki/h"
)

type Status string

const (
	Loading Status = "loading"
	Errored        = "error"
	Denied         = "denied"
	Missing        = "missing"
	Loaded         = "loaded"
)

type Stage struct {
	Lineup *Lineup
	Node   dom.Element

	Title   string
	Context string
	URL     string
	Status  Status
	Error   error

	PageNode dom.Element
	Page     *Page
}

func NewStage(lineup *Lineup, title, url string) *Stage {
	stage := &Stage{}
	stage.Lineup = lineup
	stage.Title = title
	stage.URL = url
	stage.Status = Loading

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
	if stage.Status == Loading {
		stage.Node.Class().Add("loading")
	} else {
		stage.Node.Class().Remove("loading")
	}

	stage.PageNode.SetInnerHTML("")
	switch stage.Status {
	case Loading:
	case Errored:
		stage.PageNode.AppendChild(h.Fragment(
			h.Div("title", h.Text("Error")),
			h.P(stage.Error.Error()),
		))
	case Denied:
		stage.PageNode.AppendChild(h.Fragment(
			h.Div("title", h.Text("Access Denied")),
			h.P(stage.Error.Error()),
		))
	case Missing:
		stage.PageNode.AppendChild(h.Fragment(
			h.Div("title", h.Text("Page missing")),
			h.P(stage.Error.Error()),
		))
	case Loaded:
		stage.PageNode.AppendChild(h.Fragment(
			h.Div("title", h.Text(stage.Page.Title)),
			h.Div("story", stage.RenderAll(stage.Page.Story...)...),
		))
	}
}

func (stage *Stage) fetch() {
	defer stage.Update()

	// TODO: proper threading
	r, err := http.Get(stage.URL)
	if err != nil {
		stage.Status = Errored
		stage.Error = err
		return
	}

	if r.StatusCode == 404 {
		stage.Status = Missing
		stage.Error = errors.New("Page missing")
		return
	}

	page := &Page{}
	err = json.NewDecoder(r.Body).Decode(page)
	if err != nil {
		stage.Status = Errored
		stage.Error = fmt.Errorf("Invalid page: %v", err)
		return
	}

	stage.Status = Loaded
	stage.Page = page
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
				slug := strings.ToLower(spec)
				link := h.A("", "/data/"+slug+".json", h.Text(spec))
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
