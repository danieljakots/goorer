package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

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
	fmt.Println("usage:", os.Args[0], "[-h] {summary, earnings, spendings} "+
		"path/to/data")
	fmt.Println("      ", os.Args[0], "each subcommand accepts a -date YYYY[-MM]")
	os.Exit(1)
}

func cli() (string, dateFilter, string, error) {
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
		return "", dateFilter{}, "", err
	}

	var date dateFilter
	var err error
	var dataPath string
	switch os.Args[1] {
	case "summary":
		summary.Parse(os.Args[2:])
		date, err = parseCliDate(*dateSummary)
		if len(summary.Args()) == 1 {
			dataPath = summary.Args()[0]
		} else {
			printHelp()
		}
	case "earnings":
		earnings.Parse(os.Args[2:])
		date, err = parseCliDate(*dateEarnings)
		if len(earnings.Args()) == 1 {
			dataPath = earnings.Args()[0]
		} else {
			printHelp()
		}
	case "spendings":
		spendings.Parse(os.Args[2:])
		date, err = parseCliDate(*dateSpendings)
		if len(spendings.Args()) == 1 {
			dataPath = spendings.Args()[0]
		} else {
			printHelp()
		}
	case "-h", "-help", "--help":
		printHelp()
	default:
		fmt.Println("Wrong subcommands")
		printHelp()
	}
	return os.Args[1], date, dataPath, err
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

func calcSummary(date dateFilter, entries map[string][]moneyExchange) (float64,
	float64, float64) {
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

	return earningSum, spendingSum, delta

}

func printSummary(earningSum, spendingSum, delta float64) {
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

func calcEarnings(date dateFilter, e map[string][]moneyExchange) (map[string]float64,
	[]float64) {
	earnings := make(map[string]float64)
	reverseEarnings := make(map[float64]string)
	for _, entry := range e["earnings"] {
		if acceptDate(date, entry.Date) {
			earnings[entry.With] += entry.Amount
		}
	}
	order := make([]float64, len(earnings))
	// Hopes And Prayers that there won't be conflict(s)
	i := 0
	for source, amount := range earnings {
		reverseEarnings[amount] = source
		order[i] = amount
		i++
	}
	if len(earnings) != len(reverseEarnings) {
		log.Fatal("The sums of entries from two differents source are the " +
			"same, and somehow, that's a problem.")
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(order)))

	return earnings, order
}

func printEarnings(earnings map[string]float64, order []float64) {
	if len(earnings) == 0 {
		fmt.Println("No money was earnt for that period")
	}
	for source, amount := range earnings {
		fmt.Printf("From %v: we earnt $%.2f\n", source, amount)
	}
}

func calcSpendings(date dateFilter, e map[string][]moneyExchange) (map[string]float64,
	[]float64) {
	spendings := make(map[string]float64)
	reverseSpendings := make(map[float64]string)
	for _, entry := range e["spendings"] {
		if acceptDate(date, entry.Date) {
			spendings[entry.Category] += entry.Amount
		}
	}
	order := make([]float64, len(spendings))
	// Hopes And Prayers that there won't be conflict(s)
	i := 0
	for source, amount := range spendings {
		reverseSpendings[amount] = source
		order[i] = amount
		i++
	}
	if len(spendings) != len(reverseSpendings) {
		log.Fatal("The sums of entries from two differents source are the " +
			"same, and somehow, that's a problem.")
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(order)))

	return spendings, order
}

func printSpendings(spendings map[string]float64, order []float64) {
	for source, amount := range spendings {
		fmt.Printf("For %v: we spent $%.2f\n", source, amount)
	}
}

func main() {
	mode, date, dataPath, err := cli()
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
		printSummary(calcSummary(date, e))
	case "earnings":
		printEarnings(calcEarnings(date, e))
	case "spendings":
		printSpendings(calcSpendings(date, e))
	default:
		log.Fatal("How did you end up here pal?")
	}
}
