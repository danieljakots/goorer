package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const dataPath = "testdata/"

type dateFilter struct {
	date      time.Time
	precision string
}

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

func parseCliDate(dateCli string) (dateFilter, error) {
	var timeFormat, precision string
	if len(dateCli) == 4 {
		timeFormat = "2006"
		precision = "year"
	} else if len(dateCli) == 7 {
		timeFormat = "2006-01"
		precision = "month"
	} else if len(dateCli) == 0 {
		precision = "null"
	}

	date, err := time.Parse(timeFormat, dateCli)
	if err != nil {
		return dateFilter{}, err
	}

	return dateFilter{date, precision}, nil
}

func printHelp() {
	fmt.Println("usage:", os.Args[0], "[-h] {summary, earnings, spendings} ")
	fmt.Println("      ", os.Args[0], "each subcommand accepts a -date YYYY[-MM]")
	os.Exit(1)
}

func cli() (string, dateFilter, error) {
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
		err := errors.New("You need to pick a sucommand. " +
			"Available subcommands are summary, earnings, and spendings")
		return "", dateFilter{}, err
	}

	var date dateFilter
	var err error
	switch os.Args[1] {
	case "summary":
		summary.Parse(os.Args[2:])
		date, err = parseCliDate(*dateSummary)
	case "earnings":
		earnings.Parse(os.Args[2:])
		date, err = parseCliDate(*dateEarnings)
	case "spendings":
		spendings.Parse(os.Args[2:])
		date, err = parseCliDate(*dateSpendings)
	case "-h", "-help", "--help":
		printHelp()
	default:
		fmt.Println("Wrong subcommands")
		printHelp()
	}
	return os.Args[1], date, err
}

func acceptDate(dateCli dateFilter, dateEntry time.Time) bool {
	if dateCli.precision == "year" {
		return dateCli.date.Year() == dateEntry.Year()
	} else if dateCli.precision == "month" {
		return dateCli.date.Year() == dateEntry.Year() &&
			dateCli.date.Month() == dateEntry.Month()
	} else if dateCli.precision == "null" {
		return true
	}
	return false // should be unreachable tho
}

func printSummary(date dateFilter, entries map[string][]moneyExchange) {
	var spendingSum float64
	for _, entry := range entries["spendings"] {
		if acceptDate(date, entry.Date) {
			spendingSum = spendingSum + entry.Amount
		}
	}
	var earningSum float64
	for _, entry := range entries["earnings"] {
		if acceptDate(date, entry.Date) {
			earningSum = earningSum + entry.Amount
		}
	}
	delta := earningSum - spendingSum

	fmt.Printf("You earnt $%.2f\n", earningSum)
	fmt.Printf("You spent $%.2f\n", spendingSum)
	if delta > 0 {
		fmt.Printf("You saved $%.2f\n", delta)
	} else {
		fmt.Printf("You overspent $%.2f\n", -delta)
	}
	if earningSum > 0 {
		fmt.Printf("You spent %.2f%% of your earnings\n",
			100*spendingSum/earningSum)
	}

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
		log.Print("Couldn't parse cli: ", err)
		printHelp()
	}

	categories, err := readCategoriesFile(dataPath + "categories.yml")
	if err != nil {
		log.Fatal("Couldn't parse categories file: ", err)
	}

	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	e := make(map[string][]moneyExchange)
	for _, file := range files {
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
