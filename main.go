package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/scraper"
	"bufio"
	"os"
	"strings"
	"time"
)

func main() {
	cmdi := NewCmdi()
	log.Setup(cmdi.Verbose)

	if err := cmdi.Validate(); err != nil {
		log.Error.Fatal(err)
	}

	if err := cmdi.Execute(); err != nil {
		log.Error.Fatal(err)
	}
}

func (c *Cmdi) Execute() error {
	scraper := &scraper.Scraper{
		Query:   dork.Postman(c.Org),
		Queries: c.readOrgsAsDorks(),
		Opts: scraper.Options{
			Headless:       c.Headless,
			ReqTimeoutSec:  time.Duration(c.Timeout) * time.Second,
			OutputMode:     c.mapOutputMode(),
			DelayBtwnReqMS: time.Duration(c.Delay) * time.Millisecond,
		},
	}

	if c.out.file != "" {
		defer scraper.SetupLogTarget(c.out.file)()
	}

	return scraper.DuckDuckGo()
}

func (c *Cmdi) readOrgsAsDorks() []string {
	if c.OrgList == "" {
		return nil
	}
	file, err := os.Open(c.OrgList)
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
