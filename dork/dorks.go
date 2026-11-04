package dork

import "fmt"

func Postman(org string) string {
	if org == "" {
		return ""
	}
	return fmt.Sprintf("site:postman.com %q", org)
}
