package html

import (
	"errors"
	"net/url"

	"honnef.co/go/js/dom"
	"honnef.co/go/js/xhr"

	"github.com/microcosm-cc/bluemonday"

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
	Context client.Context
	Stage   *client.Stage
	Status  Status

	Error error
	URL   string

	Content string
}

func NewView(context client.Context, title, url string) *View {
	view := &View{}
	view.Context = context
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

	u := url.URL{}
	u.Path = "/proxy"
	u.RawQuery = url.Values{"url": []string{view.URL}}.Encode()

	r := xhr.NewRequest("GET", u.RequestURI())
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
	// Sanitize the HTML.
	safehtml := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(view.Content)))

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
	child := view.Context.Open(target.TextContent(), slug)

	view.Stage.Open(child, h.IsMiddleClick(ev))
}
