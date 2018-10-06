package jira

import (
	"regexp"
)

var (
	rxCaption     = regexp.MustCompile(`(?m)^\s*(h[1-6])\.\s*(.*)\s*$`)
	rxBold        = regexp.MustCompile(`([^\s])\*([^\n\*]+)\*`)
	rxEmphasis    = regexp.MustCompile(`([^\s])\_([^\n\_]+)\_`)
	rxCite        = regexp.MustCompile(`([^\s])\?\?([^\n\?]+)\?\?`)
	rxDeleted     = regexp.MustCompile(`([^\s])\-([^\n-]+)\-`)
	rxInserted    = regexp.MustCompile(`([^\s])\+([^\n\+]+)\+`)
	rxSuperscript = regexp.MustCompile(`([^\s])\^([^\n\^]+)\^`)
	rxSubscript   = regexp.MustCompile(`([^\s])\~([^\n\~]+)\~`)
	rxMonospace   = regexp.MustCompile(`([^\s])\{\{([^\n\}]+)\}\}`)
	rxBlockquote  = regexp.MustCompile(`(?m)^bq\.\s*(.*)\s*$`)
	rxQuote       = regexp.MustCompile(`(?m)\{quote\}\s*([.\n]*)\s*\{quote\}`)
	rxColor       = regexp.MustCompile(`(?m)\{color:([^}]*)\}\s*(.*)\s*\{color\}`)
)

func ToHTML(s string) string {
	s = rxCaption.ReplaceAllString(s, `<$1>$2</$1>`)
	s = rxBold.ReplaceAllString(s, `$1<b>$2</b>`)
	s = rxEmphasis.ReplaceAllString(s, `$1<em>$2</em>`)
	s = rxCite.ReplaceAllString(s, `$1<cite>$2</cite>`)
	s = rxDeleted.ReplaceAllString(s, `$1<del>$2<del/>`)
	s = rxInserted.ReplaceAllString(s, `$1<ins>$2<ins/>`)
	s = rxSuperscript.ReplaceAllString(s, `$1<sup>$2<sup/>`)
	s = rxSubscript.ReplaceAllString(s, `$1<sub>$2<sub/>`)
	s = rxMonospace.ReplaceAllString(s, `$1<tt>$2</tt>`)
	s = rxBlockquote.ReplaceAllString(s, `<blockquote>$1</blockquote>`)
	s = rxQuote.ReplaceAllString(s, `<blockquote>$1</blockquote>`)
	s = rxColor.ReplaceAllString(s, `<span style="color:$1">$2</span>`)
	return s
}
