package fedwiki

func Slugify(title string) string {
	slug := []byte{}
	for _, r := range title {
		if r == ' ' {
			slug = append(slug, '-')
		} else if 'A' <= r && r <= 'Z' {
			slug = append(slug, byte(r-'A'+'a'))
		} else if 'a' <= r && r <= 'z' {
			slug = append(slug, byte(r))
		} else if '0' <= r && r <= '9' {
			slug = append(slug, byte(r))
		}
	}
	return string(slug)
}
