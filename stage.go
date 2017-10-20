package main

import (
	"fmt"

	"honnef.co/go/js/dom"
)

type Stage struct {
	Stages *Stages

	Node dom.Element

	URL     string
	Loading bool
	Editing bool

	PageNode dom.Element
	Page     *Page
}

func NewStage(url string) *Stage {
	stage := &Stage{}
	stage.URL = url

	stage.PageNode = h("div", "page")
	AttachOverflowIndicator(stage.PageNode)

	stage.Node = h("div", "stage",
		h("div", "status",
			h("div", "icon", text("Edit")),
			h("div", "url", text(stage.URL)),
		),
		stage.PageNode,
	)

	stage.Page = &Welcome
	stage.UpdatePage()

	return stage
}

func (stage *Stage) UpdatePage() {
	stage.PageNode.SetInnerHTML("")
	stage.PageNode.AppendChild(frag(
		h("div", "title", text(stage.Page.Title)),
		h("div", "story", stage.RenderAll(stage.Page.Story...)...),
	))
}

func (stage *Stage) RenderAll(items ...Item) []dom.Node {
	rs := []dom.Node{}
	for _, item := range items {
		rs = append(rs, stage.Render(item))
	}
	return rs
}

func (stage *Stage) Render(item Item) dom.Element {
	el := h("div", "item")

	switch item := item.(type) {
	case *Paragraph:
		el.Class().Add("paragraph")
		p := h("p", "")
		(&Parser{
			Begin: func() { p = h("p", "") },
			End:   func() { el.AppendChild(p); p = nil },
			Text: func(s string) {
				p.AppendChild(text(s))
			},
			Link: func(spec string) {
				link := h("a", "")
				link.SetAttribute("href", spec)
				link.SetTextContent(spec)
				p.AppendChild(link)
			},
		}).Run(item.Text)
	default:
		el.Class().Add("missing")
		el.SetInnerHTML(fmt.Sprintf("%+v", item))
	}

	return el
}
