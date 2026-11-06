package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/scraper"
	"bufio"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type Cmdi struct {
	out outputOpts
	//	timeout     int
	org     string
	orgList string
	verbose bool
	delayMS int
	// headless    bool
	hideBrowser bool
}

type outputOpts struct {
	file     string
	jsonl    bool
	markdown bool
}

func NewCmdi() *Cmdi {
	cmdi := &Cmdi{}

	//pflag.IntVarP(&cmdi.timeout, "timeout", "t", 5, "timeout in seconds")
	pflag.StringVarP(&cmdi.org, "org", "c", "", "organisation name to hunt api collection for")
	pflag.IntVarP(&cmdi.delayMS, "delay", "d", 2000, "delay timeout in milliseconds to wait for while multiple scraping orgs")
	pflag.StringVarP(&cmdi.orgList, "list", "l", "", "file containing a list of organisations name to hunt api collection for")
	pflag.BoolVarP(&cmdi.verbose, "verbose", "v", false, "enable verbose logging")
	//pflag.BoolVarP(&cmdi.headless, "headless", "b", false, "use an headless browser while scraping")
	pflag.BoolVarP(&cmdi.hideBrowser, "hide-browser", "i", false, "hide browser GUI while scraping (leads to captcha err)")
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

func (c *Cmdi) Validate() {
	if c.org == "" && c.orgList == "" {
		log.Error.Fatal("target organisation name(s) is required")
	}
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

func (c *Cmdi) readOrgsAsDorks() []string {
	if c.orgList == "" {
		return nil
	}
	file, err := os.Open(c.orgList)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var out []string
	for scanner.Scan() {
		org := strings.TrimSpace(scanner.Text())
		out = append(out, dork.Postman(org))
	}

	if err := scanner.Err(); err != nil {
		log.Error.Fatal("error reading organisation list file:", err)
	}
	return out
}
