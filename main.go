package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/scraper"
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
	scraper := &scraper.Scraper{
		Opts: scraper.Options{
			OutputMode: c.mapOutputMode(),
		},
	}

	if c.out.file != "" {
		defer scraper.SetupLogTarget(c.out.file)()
	}

	//if c.headless {
	ctx, cancel := scraper.NewChromeDPContext(!c.hideBrowser)
	defer cancel()
	if c.org != "" {
		_, err := scraper.SearchDDGHeadless(dork.Postman(c.org), ctx)
		if err != nil {
			log.Error.Fatalf("'scraper.SearchDDGHeadless': %v'", err)
		}
	}
	for _, q := range c.readOrgsAsDorks() {
		_, err := scraper.SearchDDGHeadless(q, ctx)
		if err != nil {
			log.Error.Fatalf("'scraper.SearchDDGHeadless: %v'", err)
		}
		time.Sleep(time.Duration(c.delayMS) * time.Millisecond)
	}
	//}
}
