package h

import "honnef.co/go/js/dom"

func IsMiddleClick(ev dom.Event) bool {
	if click, ok := ev.(*dom.MouseEvent); ok {
		return click.Button == 1
	}
	return false
}
