package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/must"
	"apihunter/scraper"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type Results struct {
	Org     string           `json:"org"`
	Results []scraper.Result `json:"results"`
}

func main() {
	timeout := pflag.IntP("timeout", "t", 5, "timeout in seconds")
	org := pflag.StringP("org", "o", "", "organisation name to hunt api collection for")
	verbose := pflag.BoolP("verbose", "v", false, "enable verbose logging")
	outputFile := pflag.StringP("output", "j", "", "path to the json output file")

	pflag.Usage = func() {
		println("Usage: apihunter [flags]")
		pflag.PrintDefaults()
	}
	pflag.Parse()
	log.Setup(*verbose)

	if *org == "" {
		pflag.Usage()
		return
	}

	log.Info.Printf("searching api docs for organisation %q", *org)
	d := dork.New(*org)

	log.Info.Println("sending request to duckduckgo API")
	resp, err := scraper.SearchDDG(d.Postman(), time.Duration(*timeout)*time.Second)
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

	if *outputFile != "" {
		file, err := os.OpenFile(*outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Error.Fatalf("error opening output file: %v", err)
		}
		defer file.Close()

		if _, err := file.Write(must.Eval(json.Marshal(Results{Org: *org, Results: res}))); err != nil {
			log.Error.Fatalf("error while writing to %q: %v", *outputFile, err)
		}
		file.Write([]byte("\n"))
	}
}
