package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/must"
	"apihunter/scraper"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

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
