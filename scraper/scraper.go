package scraper

import (
	"errors"
	"time"
)

type Scraper struct {
	Query   string
	Queries []string
	Opts    Options
}

type Options struct {
	Headless       bool
	ReqTimeoutSec  time.Duration
	DelayBtwnReqMS time.Duration
	OutputMode     OutputMode
}

type OutputMode int

const (
	JsonLOp OutputMode = iota + 1
	MarkDownOp
)

type Output struct {
	Query   string   `json:"query"`
	Results []Result `json:"results"`
}

type Result struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

func (s *Scraper) Validate() error {
	if len(s.Queries) == 0 && s.Query == "" {
		return errors.New("no input provided")
	}
	return nil
}
