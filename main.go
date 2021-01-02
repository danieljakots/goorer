package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const dataPath = "testdata/"

type moneyExchange struct {
	Amount float64
	Date time.Time
	With string
}

func readCategoriesFile(categoriesFilePath string) (map[string]string, error) {
	yamlFile, err := ioutil.ReadFile(categoriesFilePath)
	if err != nil {
		return nil, err
	}
	categories := make(map[string]string)
	err = yaml.Unmarshal(yamlFile, categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func readMonthlyFile(monthlyFilePath string) (map[string][]moneyExchange, error){
	m := make(map[string][]moneyExchange)

	yamlFile, err := ioutil.ReadFile(monthlyFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return nil, err
	}
	return m, nil

}

func parseArgDate(date string) (time.Time, error) {
	var timeFormat string
	if len(date) == 4 {
		timeFormat = "2006"
	} else if len(date) == 7 {
		timeFormat = "2006-01"
	}
	return time.Parse(timeFormat, date)
}

func cli() (string, time.Time, error) {
	summary := flag.NewFlagSet("summary", flag.ExitOnError)
	dateSummary := summary.String("date", "",
		"Focus on entries at given date (YYYY[-MM] format")
	earnings := flag.NewFlagSet("earnings", flag.ExitOnError)
	dateEarnings := earnings.String("date", "",
		"Focus on earnings at given date (YYYY[-MM] format")
	spendings := flag.NewFlagSet("spendings", flag.ExitOnError)
	dateSpendings := spendings.String("date", "",
		"Focus on spendings at given date (YYYY[-MM] format")

	if len(os.Args) < 2 {
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}

	var date time.Time
	var err error
	switch os.Args[1] {
	case "summary":
		fmt.Println(os.Args[2:])
		summary.Parse(os.Args[2:])
		date, err = parseArgDate(*dateSummary)
	case "earnings":
		earnings.Parse(os.Args[2:])
		date, err = parseArgDate(*dateEarnings)
	case "spendings":
		spendings.Parse(os.Args[2:])
		date, err = parseArgDate(*dateSpendings)
	default:
		fmt.Println("Wrong subcommands")
		os.Exit(1)
	}
	return os.Args[1], date, err
}

func main() {
	_, _, err :=cli()
	if err != nil {
		log.Fatal("Couldn't parse cli: ", err)
	}

	_, err = readCategoriesFile(dataPath + "categories.yml")
	if err != nil {
		log.Fatal("Couldn't parse categories file: ", err)
	}
	// fmt.Println(categories)

	entries, err := readMonthlyFile(dataPath + "december-20.yml")
	if err != nil {
		log.Fatal("Couldn't parse records file: ", err)
	}
	fmt.Println(entries)

	var sum float64
	for _, spending := range(entries["spendings"]) {
		sum = sum + spending.Amount
	}
	fmt.Println(sum)
}
