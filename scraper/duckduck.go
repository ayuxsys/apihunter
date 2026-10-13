package scraper

import (
	"apihunter/log"
	"apihunter/must"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

func SearchDDG(query string, timeout time.Duration) (string, error) {
	data := url.Values{}
	data.Set("q", query)
	log.Debug.Printf("searching query %q", query)

	req := must.Eval(
		http.NewRequest(
			http.MethodPost,
			"https://html.duckduckgo.com/html/",
			strings.NewReader(data.Encode()),
		),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:157.0) Gecko/20100101 Firefox/157.0")

	client := http.DefaultClient
	client.Timeout = timeout

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request to duckduckgo API: %v", err)
	}
	defer resp.Body.Close()

	respBody := string(must.Eval(io.ReadAll(resp.Body)))

	if strings.Contains(respBody, DDGNoResultsPrefix+query) {
		return respBody, errors.New(DDGNoResultsPrefix + query)
	}

	return respBody, validateResp(resp)
}

func ParseDDG(resp string) []Result {
	doc := must.Eval(
		goquery.NewDocumentFromReader(strings.NewReader(resp)),
	)

	var results []Result
	doc.Find("h2.result__title a.result__a").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		link, _ := s.Attr("href")
		results = append(results, Result{Title: title, Href: link})
	})

	return results
}
