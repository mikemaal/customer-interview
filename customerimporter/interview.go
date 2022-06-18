// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"sort"

	emailaddress "github.com/mcnijman/go-emailaddress"
)

type Domain struct {
	Name           string
	CountCostumers int
}

func DomainsParserCsv(csvPath string, logsPath string, skipLines int) ([]Domain, error) {

	if csvPath == logsPath {
		return nil, errors.New("Log path should be different from csvPath")
	}

	var warningLogger, infoLogger, errorLogger *log.Logger
	var activeLogs bool

	if logsPath != "" {
		activeLogs = true
		file, err := os.OpenFile(logsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
		errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

		infoLogger.Println("Starting the read...")
	}

	f, err := os.Open(csvPath)
	if err != nil {
		if activeLogs {
			errorLogger.Println("Error open csv:", err.Error())
		}
		return nil, err
	}
	defer f.Close()

	mapDomainCount := make(map[string]int, 0)
	mapDomainUsers := make(map[string]map[string]bool)
	var line int
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if activeLogs {
				warningLogger.Println("Error reading line ", line, " with content:", rec)
			}
			//return nil, err
		}
		if skipLines <= 0 && err == nil && len(rec) > 3 {
			var domain, user string
			domain, user, err = extractDomainAndUser(rec[2])
			if err == nil {
				usersMap, exist := mapDomainUsers[domain]
				if !exist {
					usersMap = make(map[string]bool, 0)
					usersMap[user] = true
					mapDomainCount[domain] = 1
					mapDomainUsers[domain] = usersMap
				} else {
					_, existUser := usersMap[user]
					if !existUser {
						usersMap[user] = true
						mapDomainCount[domain]++
						mapDomainUsers[domain] = usersMap
					} else {
						if activeLogs {
							warningLogger.Println("Repeated email in line ", line, " with content:", rec)
						}
					}
				}
			} else {
				if activeLogs {
					warningLogger.Println("Bad email from line ", line, " with content:", rec)
				}
			}
		} else {
			skipLines--
		}
		line++
	}

	if activeLogs {
		infoLogger.Println("Finish the read and starting sort...")
	}

	domainsResult := sortAndGenResults(mapDomainCount)

	if activeLogs {
		infoLogger.Println("Finish sort...")
	}

	return domainsResult, err

}

func sortAndGenResults(mapDomains map[string]int) []Domain {
	resultDomains := make([]Domain, len(mapDomains))
	var i = 0
	for domainName, costumersCount := range mapDomains {
		var d Domain
		d.Name = domainName
		d.CountCostumers = costumersCount
		resultDomains[i] = d
		i++
	}

	sort.Slice(resultDomains, func(i, j int) bool {
		return resultDomains[i].Name < resultDomains[j].Name
	})

	return resultDomains
}

func extractDomainAndUser(emailStr string) (domain, user string, err error) {

	var emailParsed *emailaddress.EmailAddress
	emailParsed, err = emailaddress.Parse(emailStr)
	if err != nil {
		return "", "", err
	}
	return emailParsed.Domain, emailParsed.LocalPart, err

}
