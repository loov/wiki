package mark

import (
	"errors"

	"honnef.co/go/js/dom"
	"honnef.co/go/js/xhr"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"

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
	URL   string

	Content string
}

func NewView(server client.Server, title, url string) *View {
	view := &View{}
	view.Server = server
	view.Status = Loading
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
		page.AppendChild(view.Render())
	}

	view.Stage.SetContent(page)
}

func (view *View) fetch() {
	defer view.Update()

	r := xhr.NewRequest("GET", view.URL)
	// r.Timeout = 5000
	r.ResponseType = xhr.Text

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

	view.Status = Loaded
	view.Content = r.ResponseText
}

func (view *View) Render() dom.Node {
	// Render the markdown input into HTML using Blackfriday.
	unsafehtml := blackfriday.Run([]byte(view.Content))

	// Sanitize the HTML.
	safehtml := string(bluemonday.UGCPolicy().SanitizeBytes(unsafehtml))

	mark := h.Div("markdown")
	mark.SetInnerHTML(safehtml)

	for _, link := range mark.GetElementsByTagName("a") {
		link.AddEventListener("click", false, view.LinkClicked)
		link.AddEventListener("auxclick", false, view.LinkClicked)
	}

	return mark
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

	if click, ok := ev.(*dom.MouseEvent); ok {
		if click.Button == 1 {
			view.Stage.OpenLast(child)
			return
		}

		view.Stage.OpenNext(child)
	}
}
