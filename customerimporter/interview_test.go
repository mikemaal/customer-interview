package customerimporter

import (
	"fmt"
	"testing"
)

func TestSortAndGenResults(t *testing.T) {

	var mapToSort = map[string]int{"hotmail.com": 5, "gmail.com": 3, "twitter.com": 2}

	domains := sortAndGenResults(mapToSort)

	var expected = "[{gmail.com 3} {hotmail.com 5} {twitter.com 2}]"
	var got = fmt.Sprint(domains)
	if got != expected {
		t.Errorf("Expecting "+expected+" but got %s", got)
	}

}

func TestExtractDomainAndUser(t *testing.T) {

	var user = "jacinto"
	var domain = "gmail.com"

	var email = user + "@" + domain

	domainRes, userRes, err := extractDomainAndUser(email)

	if err != nil {
		t.Errorf("Expecting no errors but got %s", err.Error())
	}

	if domainRes != domain {
		t.Errorf("Expecting "+domain+" but got %s", domainRes)
	}
	if userRes != user {
		t.Errorf("Expecting "+user+" but got %s", userRes)
	}

}

func TestDomainsParserCsv(t *testing.T) {

	res, err := DomainsParserCsv("../customers_test.csv", "../logs.txt", 1)
	if err != nil {
		t.Errorf("Expecting no errors but got %s", err.Error())
	}

	var expected = "[{cyberchimps.com 1} {github.io 2}]"
	var got = fmt.Sprint(res)
	if got != expected {
		t.Errorf("Expecting "+expected+" but got %s", got)
	}

}
