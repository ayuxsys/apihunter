package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/must"
	"apihunter/scraper"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

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

func main() {
	cmdi := NewCmdi()
	log.Setup(cmdi.Verbose)
	if err := cmdi.Validate(); err != nil {
		log.Error.Fatal(err)
	}

	if cmdi.OrgList != "" {
		file := must.Eval(os.Open(cmdi.OrgList))
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			org := scanner.Text()
			if org == "" {
				continue
			}
			cmdi.Run(org)
			if cmdi.Delay != 0 {
				log.Info.Printf("sleeping for %d seconds", cmdi.Delay)
				time.Sleep(time.Duration(cmdi.Delay) * time.Second)
			}
		}
		if scanner.Err() != nil {
			log.Error.Fatalf("bufio: %v", scanner.Err())
		}
	}

	cmdi.Run(cmdi.Org)
}

func NewCmdi() *Cmdi {
	cmdi := &Cmdi{}

	pflag.IntVarP(&cmdi.Timeout, "timeout", "t", 5, "timeout in seconds")
	pflag.StringVarP(&cmdi.Org, "org", "o", "", "organisation name to hunt api collection for")
	pflag.IntVarP(&cmdi.Delay, "delay", "d", 0, "delay timeout in seconds to wait for while multiple scraping orgs")
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

func (c *Cmdi) Run(org string) {
	log.Info.Printf("searching api docs for organisation %q", org)
	d := dork.New(org)

	log.Info.Println("sending request to duckduckgo API")
	resp, err := scraper.SearchDDG(d.Postman(), time.Duration(c.Timeout)*time.Second)
	if err != nil {
		log.Error.Fatal(err)
	}

	log.Info.Println("parsing response")
	res := scraper.ParseDDG(resp)

	if len(res) != 0 {
		log.Success.Printf("received %d results", len(res))
		for _, r := range res {
			fmt.Println(r.Href)
		}
	} else {
		log.Error.Printf("no results found for query: %q", d.Postman())
		return
	}

	if c.OutputFile != "" {
		file, err := os.OpenFile(c.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Error.Fatalf("error opening output file: %v", err)
		}
		defer file.Close()

		if _, err := file.Write(must.Eval(json.Marshal(Results{Org: org, Results: res}))); err != nil {
			log.Error.Fatalf("error while writing to %q: %v", c.OutputFile, err)
		}
		file.Write([]byte("\n"))
	}
}
