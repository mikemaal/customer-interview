package main

import (
	"fmt"

	"github.com/mikemaal/customer-interview/customerimporter"
)

func main() {

	//example
	res, err := customerimporter.DomainsParserCsv("customers.csv", "logs.txt", 1)
	if err == nil {
		fmt.Println(res)
	} else {
		fmt.Println("error:", err.Error())
	}

}
