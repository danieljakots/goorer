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
	Amount   float64
	Date     time.Time
	With     string
	Category string
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

func readMonthlyFile(monthlyFilePath string) (map[string][]moneyExchange, error) {
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

func printSummary(date time.Time, entries map[string][]moneyExchange) {
	fmt.Println("summary")
	var spending float64
	for _, entry := range entries["spendings"] {
		spending = spending + entry.Amount
	}
	var earning float64
	for _, entry := range entries["earnings"] {
		earning = earning + entry.Amount
	}
	fmt.Println(earning)
	fmt.Println(spending)

}

func printEarnings() {
	fmt.Println("earnings")
}

func printSpendings() {
	fmt.Println("spendings")
}

func main() {
	mode, date, err := cli()
	if err != nil {
		log.Fatal("Couldn't parse cli: ", err)
	}

	categories, err := readCategoriesFile(dataPath + "categories.yml")
	if err != nil {
		log.Fatal("Couldn't parse categories file: ", err)
	}
	fmt.Println(categories)

	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	e := make(map[string][]moneyExchange)
	for _, file := range files {
		fmt.Println(file.Name())
		if file.Name() == "categories.yml" {
			continue
		}
		// XXX use FS proper join
		fileEntries, err := readMonthlyFile(dataPath + file.Name())
		if err != nil {
			log.Fatal("Couldn't parse records file: ", err)
		}
		for _, spending := range fileEntries["spendings"] {
			e["spendings"] = append(e["spendings"], spending)
		}
		for _, earning := range fileEntries["earnings"] {
			e["earnings"] = append(e["earnings"], earning)
		}
	}
	fmt.Println(e)

	// Populate the Category field for each spendings entry
	if mode != "summary" {
		for n := range e["spendings"] {
			e["spendings"][n].Category = categories[e["spendings"][n].With]
		}
	}

	switch mode {
	case "summary":
		printSummary(date, e)
	case "earnings":
		printEarnings()
	case "spendings":
		printSpendings()
	default:
		log.Fatal("How did you end up here pal?")
	}
}
