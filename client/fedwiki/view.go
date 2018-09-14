package fedwiki

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"honnef.co/go/js/dom"
	"honnef.co/go/js/xhr"

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
	Server client.Server
	Stage  *client.Stage
	Status Status

	Error error
	Title string
	URL   string

	Page *Page
}

func NewView(server client.Server, title, url string) *View {
	view := &View{}
	view.Server = server
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
	view.Stage.SetSlug(h.Text(view.URL))
	view.Stage.SetButtons(h.Div("button", h.Text("Edit")))

	page := h.Div("page")
	switch view.Status {
	case Loading:
	case Errored:
		page.AppendChild(h.Fragment(
			h.H1("", h.Text("Error")),
			h.P(view.Error.Error()),
		))
	case Denied:
		page.AppendChild(h.Fragment(
			h.H1("", h.Text("Access Denied")),
			h.P(view.Error.Error()),
		))
	case Missing:
		page.AppendChild(h.Fragment(
			h.H1("", h.Text("Page missing")),
			h.P(view.Error.Error()),
		))
	case Loaded:
		page.AppendChild(h.Fragment(
			h.H1("", h.Text(view.Page.Title)),
			h.Div("story", view.RenderAll(view.Page.Story...)...),
		))
	}

	view.Stage.SetContent(page)
}

func (view *View) fetch() {
	defer view.Update()

	r := xhr.NewRequest("GET", view.URL)
	// r.Timeout = 5000
	r.ResponseType = xhr.JSON

	err := r.Send(nil)
	if err != nil {
		view.Status = Errored
		view.Error = err
		return
	}
	if r.Status == 404 {
		view.Status = Missing
		view.Error = errors.New("Page missing")
		return
	}

	page := &Page{}
	err = json.NewDecoder(strings.NewReader(r.ResponseText)).Decode(page)
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
				slug := Slugify(spec)
				url := view.Server.CreateURL(slug)
				link := h.A("", url, h.Text(spec))
				link.SetAttribute("data-slug", slug)
				link.AddEventListener("click", false, view.LinkClicked)
				p.AppendChild(link)
			},
		}).Run(item.String("text"))
	case "factory":
		el.Class().Add("factory")
		el.SetTextContent(item.String("prompt"))
	case "html":
		el.Class().Add("html")
		// TODO: add sanitization
		el.SetInnerHTML(item.String("text"))
	case "reference":
		el.Class().Add("reference")
		site := item.String("site")
		link := h.A("", "http://"+site, h.Text(item.String("title")))
		link.SetAttribute("data-site", site)
		link.AddEventListener("click", false, view.ReferenceLinkClicked)

		el.AppendChild(h.Tag("p", "",
			link,
			h.Text(" - "),
			h.Text(item.String("text")),
		))
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

	slug := target.GetAttribute("data-slug")
	url := target.GetAttribute("href")
	if slug == "" {
		slug = url
	}
	child := view.Server.Open(target.TextContent(), slug)
	view.Stage.OpenNext(child)
}

func (view *View) ReferenceLinkClicked(ev dom.Event) {
	target := ev.Target()
	ev.StopPropagation()
	ev.PreventDefault()

	site := target.GetAttribute("data-site")
	server := &Server{"http://" + site + "/"}
	child := server.Open(target.TextContent(), "welcome-visitors")
	view.Stage.OpenNext(child)
}
