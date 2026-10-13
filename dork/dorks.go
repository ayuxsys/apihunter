package dork

import "fmt"

type Dork struct {
	Org string
}

func New(org string) *Dork {
	return &Dork{Org: org}
}

func (d *Dork) Postman() string {
	return fmt.Sprintf("site:postman.com %q", d.Org)
}
