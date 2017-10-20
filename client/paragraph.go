package client

type Paragraph struct {
	Text string
}

type Parser struct {
	Begin func()
	Text  func(s string)
	Link  func(spec string)
	End   func()
}

func (parser *Parser) Run(text string) {
	parser.Begin()

	s := 0
	inLink := false
	for p, r := range text {
		switch r {
		case '[':
			if p+1 < len(text) && text[p+1] == '[' {
				if s < p {
					parser.Text(text[s:p])
				}
				inLink = true
				s = p + 2
			}
		case ']':
			if p+1 < len(text) && text[p+1] == ']' {
				if inLink {
					if s < p {
						parser.Link(text[s:p])
					}
				} else {
					if s < p {
						parser.Text(text[s:p])
					}
				}
				inLink = false
				s = p + 2
			}
		case '\n':
			if p+1 < len(text) && text[p+1] == '\n' {
				if s < p {
					parser.Text(text[s:p])
				}
				parser.End()
				parser.Begin()
				s = p + 2
			}
		}
	}
	if s < len(text) {
		parser.Text(text[s:])
	}

	parser.End()
}
