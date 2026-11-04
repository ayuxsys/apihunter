package main

import (
	"apihunter/scraper"
	"errors"

	"github.com/spf13/pflag"
)

type Results struct {
	Org     string           `json:"org"`
	Results []scraper.Result `json:"results"`
}

type Cmdi struct {
	out      outputOpts
	Timeout  int
	Org      string
	OrgList  string
	Verbose  bool
	Delay    int
	Headless bool
}

type outputOpts struct {
	file     string
	jsonl    bool
	markdown bool
}

func NewCmdi() *Cmdi {
	cmdi := &Cmdi{}

	pflag.IntVarP(&cmdi.Timeout, "timeout", "t", 5, "timeout in seconds")
	pflag.StringVarP(&cmdi.Org, "org", "c", "", "organisation name to hunt api collection for")
	pflag.IntVarP(&cmdi.Delay, "delay", "d", 2000, "delay timeout in milliseconds to wait for while multiple scraping orgs")
	pflag.StringVarP(&cmdi.OrgList, "list", "l", "", "file containing a list of organisations name to hunt api collection for")
	pflag.BoolVarP(&cmdi.Verbose, "verbose", "v", false, "enable verbose logging")
	pflag.BoolVarP(&cmdi.Headless, "headless", "b", false, "use an headless browser while scraping")
	pflag.StringVarP(&cmdi.out.file, "output", "o", "", "path to the output file (note: can only use one output mode at a time)")
	pflag.BoolVarP(&cmdi.out.jsonl, "json", "j", false, "enable jsonl output")
	pflag.BoolVarP(&cmdi.out.markdown, "md", "m", false, "enable markdown output")

	pflag.Usage = func() {
		println("Usage: apihunter [flags]")
		pflag.PrintDefaults()
	}
	pflag.Parse()

	return cmdi
}

func (c *Cmdi) Validate() error {
	if c.Org == "" && c.OrgList == "" {
		return errors.New("target organisation name(s) is required")
	}
	return nil
}

func (c *Cmdi) mapOutputMode() scraper.OutputMode {
	if c.out.jsonl {
		return scraper.JsonLOp
	}
	if c.out.markdown {
		return scraper.MarkDownOp
	}
	return scraper.MarkDownOp // default (todo: add pretty output)
}
