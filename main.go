package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/scraper"
	"errors"
	"time"
)

func main() {
	cmdi := NewCmdi()
	log.Setup(cmdi.verbose)
	scraper.SetupLogger(cmdi.verbose)
	cmdi.Validate()
	cmdi.Execute()
}

func (c *Cmdi) Execute() {
	scrpr := &scraper.Scraper{
		Opts: scraper.Options{
			OutputMode: c.mapOutputMode(),
		},
	}

	if c.out.file != "" {
		defer scrpr.SetupLogTarget(c.out.file)()
	}

	ctx, cancel := scrpr.NewChromeDPContext(!c.hideBrowser)
	defer cancel()
	if c.org != "" {
		_, err := scrpr.SearchDDGHeadless(dork.Postman(c.org), ctx)
		if err != nil {
			log.Error.Fatalf("'scraper.SearchDDGHeadless': %v'", err)
		}
	}
	for _, q := range c.readOrgsAsDorks() {
		_, err := scrpr.SearchDDGHeadless(q, ctx)
		if err != nil {
			if errors.Is(err, scraper.CaptchaErr) {
				log.Warn.Printf("failed to search %q. captcha error", q)
				continue
			}
			log.Error.Fatalf("'scraper.SearchDDGHeadless: %v'", err)
		}
		time.Sleep(time.Duration(c.delayMS) * time.Millisecond)
	}
}
