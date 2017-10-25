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

type Context struct {
	Stage  *Stage
	Status Status

	Error error
	Title string
	URL   string

	Page *Page
}

func NewContext(title, url string) *Context {
	context := &Context{}
	context.Status = Loading
	context.Title = title
	context.URL = url
	return context
}

func (context *Context) Attach(stage *Stage) {
	context.Stage = stage

	context.Update()
	go context.fetch()
}

func (context *Context) Detach() {
	// TODO: cancel all pending requests
	context.Stage = nil
}

func (context *Context) Update() {
	context.Stage.SetTag("loading", context.Status == Loading)
	context.Stage.SetSlug(h.Text("[loov.io] " + context.URL))
	context.Stage.SetButtons(h.Div("button", h.Text("Edit")))

	page := h.Div("page")
	switch context.Status {
	case Loading:
	case Errored:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text("Error")),
			h.P(context.Error.Error()),
		))
	case Denied:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text("Access Denied")),
			h.P(context.Error.Error()),
		))
	case Missing:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text("Page missing")),
			h.P(context.Error.Error()),
		))
	case Loaded:
		page.AppendChild(h.Fragment(
			h.Div("title", h.Text(context.Page.Title)),
			h.Div("story", context.RenderAll(context.Page.Story...)...),
		))
	}

	context.Stage.SetContent(page)
}

func (context *Context) fetch() {
	defer context.Update()

	// TODO: proper threading
	r, err := http.Get(context.URL)
	if err != nil {
		context.Status = Errored
		context.Error = err
		return
	}

	if r.StatusCode == 404 {
		context.Status = Missing
		context.Error = errors.New("Page missing")
		return
	}

	page := &Page{}
	err = json.NewDecoder(r.Body).Decode(page)
	if err != nil {
		context.Status = Errored
		context.Error = fmt.Errorf("Invalid page: %v", err)
		return
	}

	context.Status = Loaded
	context.Page = page
}

func (context *Context) RenderAll(items ...Item) []dom.Node {
	rs := []dom.Node{}
	for _, item := range items {
		rs = append(rs, context.Render(item))
	}
	return rs
}

func (context *Context) Render(item Item) dom.Element {
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
				link.AddEventListener("click", false, context.LinkClicked)
				p.AppendChild(link)
			},
		}).Run(item.String("text"))
	default:
		el.Class().Add("missing")
		el.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return el
}

func (context *Context) LinkClicked(ev dom.Event) {
	target := ev.Target()
	ev.StopPropagation()
	ev.PreventDefault()

	child := NewContext(target.TextContent(), target.GetAttribute("href"))
	context.Stage.Open(child)
}
