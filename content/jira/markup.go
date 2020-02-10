package jira

import (
	"regexp"
	"strings"
)

var (
	// See https://jira.atlassian.com/secure/WikiRendererHelpAction.jspa

	// headings
	rxCaption = regexp.MustCompile(`(?m)^\s*(h[1-6])\.\s*(.*)\s*$`)

	// text effects
	rxBold        = regexp.MustCompile(`([^\s])\*([^\n\*]+)\*`)
	rxEmphasis    = regexp.MustCompile(`([^\s])\_([^\n\_]+)\_`)
	rxCite        = regexp.MustCompile(`([^\s])\?\?([^\n\?]+)\?\?`)
	rxDeleted     = regexp.MustCompile(`([^\s])\-([^\n-]+)\-`)
	rxInserted    = regexp.MustCompile(`([^\s])\+([^\n\+]+)\+`)
	rxSuperscript = regexp.MustCompile(`([^\s])\^([^\n\^]+)\^`)
	rxSubscript   = regexp.MustCompile(`([^\s])\~([^\n\~]+)\~`)
	rxMonospace   = regexp.MustCompile(`([^\s])\{\{([^\n\}]+)\}\}`)
	rxBlockquote  = regexp.MustCompile(`(?m)^bq\.\s*(.*)\s*$`)
	rxQuote       = regexp.MustCompile(`(?ms)\{quote\}\s*(.*?)\s*\{quote\}`)
	rxColor       = regexp.MustCompile(`(?ms)\{color:([^}]*)\}\s*(.*?)\s*\{color\}`)

	// text breaks
	rxParagraph       = regexp.MustCompile(`\n\n+`)
	rxLinebreak       = regexp.MustCompile(`\\\\`)
	rxHorizontalRuler = regexp.MustCompile(`----`)
	rxEmDash          = regexp.MustCompile(`---`)
	rxEnDash          = regexp.MustCompile(`--`)

	// Links
	// Lists
	// Images
	// Attachments
	// Tables
	// Advanced Formatting
	// Misc
	// escaping

	rxEmotes = strings.NewReplacer(
		// Emotes
		`:)`, `🙂`,
		`:(`, `🙁`,
		`:P`, `😛`,
		`:D`, `😀`,
		`;)`, `😉`,
		`(y)`, `👍`,
		`(n)`, `👎`,
		`(i)`, `ℹ`,
		`(/)`, `✔`,
		`(x)`, `✘`,
		`(!)`, `⚠`,
		// Notation
		`(+)`, `⊕`,
		`(-)`, `⊝`,
		`(?)`, `❓`,
		`(on)`, `💡`,
		`(off)`, `<span style="color:gray">💡</span>`,
		`(*)`, `🌟`,
		`(*r)`, `<span class="outline" style="color:red;">★</span>`,
		`(*g)`, `<span class="outline" style="color:green;">★</span>`,
		`(*b)`, `<span class="outline" style="color:blue;">★</span>`,
		`(*y)`, `<span class="outline" style="color:yellow;">★</span>`,
		`(flag)`, `🏴`,
		`(flagoff)`, `🏳`,
	)
)

func ToHTML(s string) string {
	s = rxEmotes.Replace(s)

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
	s = `<p>` + rxParagraph.ReplaceAllString(s, `</p><p>`) + `</p>`
	s = rxLinebreak.ReplaceAllString(s, `<br>`)
	s = rxHorizontalRuler.ReplaceAllString(s, `<hr>`)
	s = rxEmDash.ReplaceAllString(s, `—`)
	s = rxEnDash.ReplaceAllString(s, `–`)

	return s
}
