package main

import (
	"apihunter/dork"
	"apihunter/log"
	"apihunter/must"
	"apihunter/scraper"
	"bufio"
	"errors"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
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
		sleep := c.deferSleepTimeSec()
		log.Info.Printf("sleep %fs", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func (c *Cmdi) deferSleepTimeSec() float64 {
	if c.randomDelay == "" {
		return float64(c.delaySec)
	}
	itm := strings.Split(c.randomDelay, ",")
	min, max := must.Eval(strconv.Atoi(itm[0])), must.Eval(strconv.Atoi(itm[1]))
	return randFloat(min, max)
}

func randFloat(min, max int) float64 {
	return float64(min) + rand.Float64()*float64(max-min)
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
