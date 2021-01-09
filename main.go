package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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

type kvlist []kv

type kv struct {
	key        string
	value      float64
	percentage float64
}

func (kvl kvlist) Len() int {
	return len(kvl)
}
func (kvl kvlist) Swap(i, j int) {
	kvl[i], kvl[j] = kvl[j], kvl[i]
}
func (kvl kvlist) Less(i, j int) bool {
	return kvl[i].value < kvl[j].value
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

func readAllMonthlyFiles(files []os.FileInfo, dataPath string) (
	map[string][]moneyExchange, error) {
	e := make(map[string][]moneyExchange)
	for _, file := range files {
		if file.Name() == "categories.yml" {
			continue
		}
		fileEntries, err := readMonthlyFile(path.Join(dataPath, file.Name()))
		if err != nil {
			crafterErr := fmt.Sprintf("Couldn't parse records file %v: %v",
				file.Name(), err)
			return nil, errors.New(crafterErr)
		}
		for _, spending := range fileEntries["spendings"] {
			e["spendings"] = append(e["spendings"], spending)
		}
		for _, earning := range fileEntries["earnings"] {
			e["earnings"] = append(e["earnings"], earning)
		}
	}
	return e, nil
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
		"path/to/data\n")
	fmt.Println("Each subcommand accepts a --date YYYY[-MM] to filter on a " +
		"subset of entries")
	fmt.Println("The spendings subcommand accept a -d/--details. This prints " +
		"spendings\nwithout using categories.")
	os.Exit(1)
}

func cli() (string, dateFilter, string, bool, error) {
	summary := flag.NewFlagSet("summary", flag.ExitOnError)
	dateSummary := summary.String("date", "",
		"Focus on entries at given date (YYYY[-MM] format")
	earnings := flag.NewFlagSet("earnings", flag.ExitOnError)
	dateEarnings := earnings.String("date", "",
		"Focus on earnings at given date (YYYY[-MM] format")
	spendings := flag.NewFlagSet("spendings", flag.ExitOnError)
	dateSpendings := spendings.String("date", "",
		"Focus on spendings at given date (YYYY[-MM] format")
	dSpendings := spendings.Bool("d", false, "Details mode")
	detailsSpendings := spendings.Bool("details", false, "Details mode")

	if len(os.Args) < 2 {
		err := errors.New("You need to pick a sucommand. " +
			"Available subcommands are summary, earnings, and spendings")
		return "", dateFilter{}, "", false, err
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

	var details bool
	if *dSpendings {
		details = true
	}
	if *detailsSpendings {
		details = true
	}

	return os.Args[1], date, dataPath, details, err
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
		if !acceptDate(date, entry.Date) {
			continue
		}
		spendingSum = spendingSum + entry.Amount
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

func calcEarnings(date dateFilter, e map[string][]moneyExchange) []kv {
	earningSum, _, _ := calcSummary(date, e)
	earnings := make([]kv, 0)
OUTER:
	for _, entry := range e["earnings"] {
		if !acceptDate(date, entry.Date) {
			continue
		}

		for n := range earnings {
			if earnings[n].key != entry.With {
				continue
			}
			earnings[n].value += entry.Amount
			continue OUTER
		}
		earnings = append(earnings, kv{entry.With, entry.Amount, 0})
	}
	sort.Sort(sort.Reverse(kvlist(earnings)))

	for n := range earnings {
		earnings[n].percentage = earnings[n].value / earningSum * 100
	}

	return earnings
}

func printIngs(data []kv) {
	if len(data) == 0 {
		fmt.Println("No money was earnt for that period")
	}
	for n := range data {
		percentage := fmt.Sprintf("%.2f%%", data[n].percentage)
		value := fmt.Sprintf("$%.2f,", data[n].value)
		fmt.Printf("From %-25s: we earnt %-11s this is %6s\n",
			data[n].key, value, percentage)
	}
}

func calcSpendings(date dateFilter, e map[string][]moneyExchange, details bool) []kv {
	_, spendingSum, _ := calcSummary(date, e)
	spendings := make([]kv, 0)
OUTER:
	for _, entry := range e["spendings"] {
		if !acceptDate(date, entry.Date) {
			continue
		}
		if details {
			for n := range spendings {
				if spendings[n].key != entry.With {
					continue
				}
				spendings[n].value += entry.Amount
				continue OUTER
			}
			spendings = append(spendings, kv{entry.With, entry.Amount, 0})
		} else {
			for n := range spendings {
				if spendings[n].key != entry.Category {
					continue
				}
				spendings[n].value += entry.Amount
				continue OUTER
			}
			spendings = append(spendings, kv{entry.Category,
				entry.Amount, 0})
		}
	}

	for n := range spendings {
		spendings[n].percentage = spendings[n].value / spendingSum * 100
	}

	sort.Sort(sort.Reverse(kvlist(spendings)))
	return spendings
}

func main() {
	mode, date, dataPath, details, err := cli()
	if err != nil {
		log.Print("Couldn't parse cli: ", err)
		printHelp()
	}

	categories, err := readCategoriesFile(path.Join(dataPath, "categories.yml"))
	if err != nil {
		log.Fatal("Couldn't parse categories file: ", err)
	}

	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	e, err := readAllMonthlyFiles(files, dataPath)
	if err != nil {
		log.Fatal(err)
	}

	// Populate the Category field for each spendings entry
	if mode != "summary" {
		for n := range e["spendings"] {
			if cat, ok := categories[e["spendings"][n].With]; ok {
				e["spendings"][n].Category = cat
				continue
			}
			log.Fatal("Couldn't find category for ",
				e["spendings"][n].With)
		}
	}

	switch mode {
	case "summary":
		printSummary(calcSummary(date, e))
	case "earnings":
		printIngs(calcEarnings(date, e))
	case "spendings":
		printIngs(calcSpendings(date, e, details))
	default:
		log.Fatal("How did you end up here pal?")
	}
}
