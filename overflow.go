package main

import (
	"time"

	"honnef.co/go/js/dom"
)

func AttachOverflowIndicator(elem dom.Element) {
	token := 0

	update := func(_ time.Duration) {
		token = 0
		obj := elem.Underlying()
		scrollTop := obj.Get("scrollTop").Int()
		clientHeight := obj.Get("clientHeight").Int()
		scrollHeight := obj.Get("scrollHeight").Int()

		if scrollTop > 0 {
			elem.Class().Add("overflow-top")
		} else {
			elem.Class().Remove("overflow-top")
		}

		if scrollTop+clientHeight < scrollHeight {
			elem.Class().Add("overflow-bottom")
		} else {
			elem.Class().Remove("overflow-bottom")
		}
	}

	request := func(_ dom.Event) {
		win := dom.GetWindow()
		win.CancelAnimationFrame(token)
		token = dom.GetWindow().RequestAnimationFrame(update)
	}

	elem.AddEventListener("scroll", false, request)

	request(nil)
}
