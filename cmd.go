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
	Timeout    int
	Org        string
	OrgList    string
	Verbose    bool
	OutputFile string
	Delay      int
}

func NewCmdi() *Cmdi {
	cmdi := &Cmdi{}

	pflag.IntVarP(&cmdi.Timeout, "timeout", "t", 5, "timeout in seconds")
	pflag.StringVarP(&cmdi.Org, "org", "o", "", "organisation name to hunt api collection for")
	pflag.IntVarP(&cmdi.Delay, "delay", "d", 5, "delay timeout in seconds to wait for while multiple scraping orgs")
	pflag.StringVarP(&cmdi.OrgList, "list", "l", "", "file containing a list of organisations name to hunt api collection for")
	pflag.BoolVarP(&cmdi.Verbose, "verbose", "v", false, "enable verbose logging")
	pflag.StringVarP(&cmdi.OutputFile, "output", "j", "", "path to the json output file")

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
