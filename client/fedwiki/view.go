package fedwiki

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"honnef.co/go/js/dom"

	"github.com/loov/wiki/client"
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

type View struct {
	Stage  *client.Stage
	Status Status

	Error error
	Title string
	URL   string

	Page *Page
}

func NewView(title, url string) *View {
	view := &View{}
	view.Status = Loading
	view.Title = title
	view.URL = url
	return view
}

func (view *View) Attach(stage *client.Stage) {
	view.Stage = stage

	view.Update()
	go view.fetch()
}

func (view *View) Detach() {
	// TODO: cancel all pending requests
	view.Stage = nil
}

func (view *View) Update() {
	view.Stage.SetTag("loading", view.Status == Loading)
	view.Stage.SetSlug(h.Text("[loov.io] " + view.URL))
	view.Stage.SetButtons(h.Div("button", h.Text("Edit")))

	page := h.Div("page")
	switch view.Status {
	case Loading:
	case Errored:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text("Error")),
			h.P(view.Error.Error()),
		))
	case Denied:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text("Access Denied")),
			h.P(view.Error.Error()),
		))
	case Missing:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text("Page missing")),
			h.P(view.Error.Error()),
		))
	case Loaded:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text(view.Page.Title)),
			h.Div("story", view.RenderAll(view.Page.Story...)...),
		))
	}

	view.Stage.SetContent(page)
}

func (view *View) fetch() {
	defer view.Update()

	// TODO: proper threading
	r, err := http.Get(view.URL)
	if err != nil {
		view.Status = Errored
		view.Error = err
		return
	}

	if r.StatusCode == 404 {
		view.Status = Missing
		view.Error = errors.New("Page missing")
		return
	}

	page := &Page{}
	err = json.NewDecoder(r.Body).Decode(page)
	if err != nil {
		view.Status = Errored
		view.Error = fmt.Errorf("Invalid page: %v", err)
		return
	}

	view.Status = Loaded
	view.Page = page
}

func (view *View) RenderAll(items ...Item) []dom.Node {
	rs := []dom.Node{}
	for _, item := range items {
		rs = append(rs, view.Render(item))
	}
	return rs
}

func (view *View) Render(item Item) dom.Element {
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
				link.AddEventListener("click", false, view.LinkClicked)
				p.AppendChild(link)
			},
		}).Run(item.String("text"))
	default:
		el.Class().Add("missing")
		el.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return el
}

func (view *View) LinkClicked(ev dom.Event) {
	target := ev.Target()
	ev.StopPropagation()
	ev.PreventDefault()

	child := NewView(target.TextContent(), target.GetAttribute("href"))
	view.Stage.OpenNext(child)
}
