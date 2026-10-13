package scraper

import "testing"

func TestParseDDG(t *testing.T) {
	html := `
	<h2 class="result__title">
		<a rel="nofollow" class="result__a"
		   href="https://en.wikipedia.org/wiki/Coffee">
			Coffee - Wikipedia
		</a>
	</h2>`

	res := ParseDDG(html)
	t.Log(res)
}
