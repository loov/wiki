package fedwiki

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	Context client.Context
	Stage   *client.Stage
	Status  Status

	Error error
	Title string
	URL   string

	Page *Page
}

func NewView(context client.Context, title, url string) *View {
	view := &View{}
	view.Context = context
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
	if view.Stage == nil {
		return
	}

	view.Stage.SetTag("loading", view.Status == Loading)
	view.Stage.SetSlug(h.Text(view.URL))
	// view.Stage.SetButtons(h.Div("button", h.Text("Edit")))

	page := h.Div("page")
	switch view.Status {
	case Loading:
		page.AppendChild(h.Fragment(
			h.H1("story-header", h.Text(view.Title)),
		))
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
			h.H1("story-header", h.Text(view.Page.Title)),
			h.Div("story", view.RenderAll(view.Page.Story...)...),
		))
	}

	view.Stage.SetContent(page)
}

func (view *View) fetch() {
	defer view.Update()

	r := xhr.NewRequest("GET", view.URL)
	// r.Timeout = 5000
	r.ResponseType = xhr.Text // TODO: use json

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
	if page.Title != "" {
		view.Title = page.Title
	}
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
				url := view.Context.CreateURL(slug)
				link := h.A("", url, h.Text(spec))
				link.SetAttribute("data-slug", slug)
				link.AddEventListener("click", false, view.LinkClicked)
				link.AddEventListener("auxclick", false, view.LinkClicked)
				p.AppendChild(link)
			},
		}).Run(item.String("text"))

	case "markdown":
		el.Class().Add("markdown")

		// Render the markdown input into HTML using Blackfriday.
		unsafehtml := blackfriday.Run([]byte(item.String("text")))
		// Sanitize the HTML.
		safehtml := string(bluemonday.UGCPolicy().SanitizeBytes(unsafehtml))

		el.SetInnerHTML(safehtml)

		for _, link := range el.GetElementsByTagName("a") {
			link.AddEventListener("click", false, view.LinkClicked)
			link.AddEventListener("auxclick", false, view.LinkClicked)
		}

	case "factory":
		el.Class().Add("factory")
		el.SetTextContent(item.String("prompt"))

	case "html":
		el.Class().Add("html")
		safehtml := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(item.String("text"))))
		el.SetInnerHTML(safehtml)

	case "reference":
		el.Class().Add("reference")
		site := item.String("site")
		//TODO: fix reference and moving between sites
		link := h.A("", site, h.Text(item.String("title")))
		link.SetAttribute("data-site", site)
		link.AddEventListener("click", false, view.ReferenceLinkClicked)
		link.AddEventListener("auxclick", false, view.ReferenceLinkClicked)

		el.AppendChild(h.Tag("p", "",
			link,
			h.Text(" - "),
			h.Text(item.String("text")),
		))

	case "image":
		el.Class().Add("image")
		el.AppendChild(h.Img("thumbnail", item.String("url")))
		el.AppendChild(h.P(item.String("text")))

	default:
		el.Class().Add("missing")
		el.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return el
}

func (view *View) LinkClicked(ev dom.Event) {
	if view.Stage == nil {
		return
	}

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

func (view *View) ReferenceLinkClicked(ev dom.Event) {
	target := ev.Target()
	ev.StopPropagation()
	ev.PreventDefault()

	site := target.GetAttribute("data-site")
	//TODO: fix reference and moving between sites
	// context := &Context{"http://" + site + "/"}
	child := view.Context.Open(target.TextContent(), site)

	view.Stage.Open(child, h.IsMiddleClick(ev))
}
