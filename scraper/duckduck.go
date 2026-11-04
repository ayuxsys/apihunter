package scraper

import (
	"apihunter/log"
	"apihunter/must"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	DuckDuckGoURL = "https://html.duckduckgo.com/html/"
)

func (s *Scraper) DuckDuckGo() error {
	if err := s.Validate(); err != nil {
		return err
	}

	if s.Opts.Headless {
		opts := append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("start-maximized", true),
		)

		allocCtx, allocCancel := chromedp.NewExecAllocator(
			context.Background(),
			opts...,
		)
		defer allocCancel()

		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()
		if s.Query != "" {
			out, err := s.searchDDGHeadless(ctx)
			if err != nil {
				return err
			}
			s.logOutput(out)
		}
		if len(s.Queries) != 0 {
			for _, q := range s.Queries {
				if q == "" {
					continue
				}
				s.Query = q
				out, err := s.searchDDGHeadless(ctx)
				if err == CaptchaErr {
					log.Warn.Printf("error searching query %q. captcha error detected", q)
					continue
				}
				if err != nil {
					return err
				}
				s.logOutput(out)
				time.Sleep(s.Opts.DelayBtwnReqMS)
			}
		}
		return nil
	}

	if s.Query != "" {
		out, err := s.searchDDG()
		if err != nil {
			return err
		}
		s.logOutput(out)
	}

	if len(s.Queries) != 0 {
		for _, q := range s.Queries {
			if q == "" {
				continue
			}
			s.Query = q
			out, err := s.searchDDG()
			if err != nil {
				return err
			}
			s.logOutput(out)
			time.Sleep(s.Opts.DelayBtwnReqMS)
		}
	}

	return nil
}

func (s *Scraper) searchDDGHeadless(ctx context.Context) (Output, error) {
	out := Output{Query: s.Query}
	// ctx, cancel := chromedp.NewContext(context.Background())
	// ctx, cancel = context.WithTimeout(ctx, s.Timeout)
	// defer cancel()

	if err := chromedp.Do(ctx, chromedp.Navigate(DuckDuckGoURL), chromedp.WaitVisible("body")); err != nil {
		return out, err
	}

	queryJSON := must.Eval(json.Marshal(s.Query))

	evalCode := fmt.Sprintf(`
	(async () => {
    const url = "https://html.duckduckgo.com/html/";
    const data = new URLSearchParams({q: %s});

    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
        },
        body: data
    });

    if (!response.ok) {
        throw new Error("HTTP status: " + response.status);
    }

    return await response.text();
	})();
	`, queryJSON)

	html, err := chromedp.Run(ctx,
		chromedp.Evaluate[string](
			evalCode,
			chromedp.EvalAwaitPromise,
		),
	)
	if err != nil {
		return out, err
	}
	if err := s.validateDDGResponse(html); err != nil {
		return out, err
	}

	out.Results = parseDDG(html)

	return out, nil
}

func (s *Scraper) searchDDG() (Output, error) {
	out := Output{Query: s.Query}
	data := url.Values{}
	data.Set("q", s.Query)
	log.Debug.Printf("searching query %q", s.Query)

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
	client.Timeout = s.Opts.ReqTimeoutSec

	resp, err := client.Do(req)
	if err != nil {
		return out, fmt.Errorf("send request to duckduckgo API: %v", err)
	}
	defer resp.Body.Close()

	respBody := string(must.Eval(io.ReadAll(resp.Body)))

	if err := s.validateDDGResponse(respBody); err != nil {
		return out, err
	}

	out.Results = parseDDG(respBody)

	return out, nil
}

func parseDDG(resp string) []Result {
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

func (s *Scraper) validateDDGResponse(resp string) error {
	if strings.Contains(resp, DDGNoResultsPrefix+s.Query) {
		return errors.New(DDGNoResultsPrefix + s.Query)
	}
	if strings.Contains(resp, "Unfortunately, bots use DuckDuckGo too.") {
		return CaptchaErr
	}
	return nil
}
