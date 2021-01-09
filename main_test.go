package main

import (
	"reflect"
	"testing"
	"time"
)

func TestReadCategoriesFile(t *testing.T) {
	categories, err := readCategoriesFile("testdata/categories.yml")
	if err != nil {
		t.Error("readCategoriesFile() failed")
		t.Fatal(err)
	}
	shouldBe := make(map[string]string, 2)
	shouldBe["cat food shop"] = "cat"
	shouldBe["rent"] = "home"
	shouldBe["saq"] = "wine"

	if !reflect.DeepEqual(categories, shouldBe) {
		t.Errorf("readCategoriesFile() failed: got %v, wanted %v",
			categories, shouldBe)
	}
}

func TestReadMonthlyFile(t *testing.T) {
	entries, err := readMonthlyFile("testdata/december-20.yml")
	if err != nil {
		t.Error("readMonthlyFile() failed")
		t.Fatal(err)
	}
	shouldBe := make(map[string][]moneyExchange, 3)
	date, err := time.Parse("2006-01-02", "2020-12-25")
	if err != nil {
		t.Error("time.Parse in readMonthlyFile() failed")
		t.Fatal(err)
	}
	shouldBe["earnings"] = append(shouldBe["earnings"],
		moneyExchange{4321.0, date, "Company", ""})
	shouldBe["earnings"] = append(shouldBe["earnings"],
		moneyExchange{5.0, date, "Santa Claus", ""})

	date, err = time.Parse("2006-01-02", "2020-12-01")
	if err != nil {
		t.Error("time.Parse in readMonthlyFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{1234.0, date, "rent", ""})

	date, err = time.Parse("2006-01-02", "2020-12-12")
	if err != nil {
		t.Error("time.Parse in readMonthlyFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{13.37, date, "cat food shop", ""})

	date, err = time.Parse("2006-01-02", "2020-12-21")
	if err != nil {
		t.Error("time.Parse in readMonthlyFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{73.31, date, "saq", ""})

	date, err = time.Parse("2006-01-02", "2020-12-25")
	if err != nil {
		t.Error("time.Parse in readMonthlyFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{42.24, date, "cat food shop", ""})

	if !reflect.DeepEqual(entries, shouldBe) {
		t.Errorf("readMonthlyFile() failed: got\n%v\nwanted\n%v",
			entries, shouldBe)
	}
}

func TestParseArgDate(t *testing.T) {
	year := "2020"
	shouldBeYear, err := time.Parse("2006", year)
	if err != nil {
		t.Error("time.Parse in TestParseArgDate() failed")
		t.Fatal(err)
	}

	result, err := parseCliDate(year)
	if err != nil {
		t.Error("parseArgDate(year) gave an error")
		t.Fatal(err)
	}
	if shouldBeYear != result.date {
		t.Errorf("parseArgDate(year) result is unexpected: got %v, wanted %v",
			result, shouldBeYear)
	}

	yearMonth := "2020-12"
	shouldBeYearMonth, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		t.Error("time.Parse in TestParseArgDate() failed")
		t.Fatal(err)
	}

	result, err = parseCliDate(yearMonth)
	if err != nil {
		t.Error("parseArgDate(yearMonth) gave an error")
		t.Fatal(err)
	}
	if shouldBeYearMonth != result.date {
		t.Error("parseArgDate(yearMonth) result is unexpected:")
		t.Errorf("got %v, wanted %v", result, shouldBeYearMonth)
	}
}

func TestCalcSummary(t *testing.T) {
	shouldBeEarning := 4321.00
	shouldBeSpending := 1362.9
	shouldBeDelta := 2963.08
	date := dateFilter{time.Now(), "null"}
	entries, _ := readMonthlyFile("testdata/december-20.yml")
	e, s, d := calcSummary(date, entries)
	if e != shouldBeEarning && s != shouldBeSpending && d != shouldBeDelta {
		t.Error("calcSummary() result is unexpected:")
		t.Errorf("got %v, wanted %v", e, shouldBeEarning)
		t.Errorf("got %v, wanted %v", s, shouldBeSpending)
		t.Fatalf("got %v, wanted %v", d, shouldBeDelta)
	}
}

func TestCalcEarnings(t *testing.T) {
	shouldBeEarnings := make(map[string]float64)
	shouldBeEarnings["Company"] = 4321
	shouldBeEarnings["Santa Claus"] = 5

	date := dateFilter{time.Now(), "null"}
	entries, _ := readMonthlyFile("testdata/december-20.yml")
	earnings := calcEarnings(date, entries)
	if !reflect.DeepEqual(earnings, shouldBeEarnings) {
		t.Error("calcEarnings() earnings result is unexpected:")
		t.Errorf("got %v, wanted %v", earnings, shouldBeEarnings)
	}
}

func TestCalcSpendings(t *testing.T) {
	date := dateFilter{time.Now(), "null"}
	e, _ := readMonthlyFile("testdata/december-20.yml")
	categories, err := readCategoriesFile("testdata/categories.yml")
	if err != nil {
		t.Fatal("Couldn't parse categories file: ", err)
	}
	// Populate the Category field for each spendings entry
	for n := range e["spendings"] {
		e["spendings"][n].Category = categories[e["spendings"][n].With]
	}

	// without details
	shouldBeSpendings := make(map[string]float64)
	shouldBeSpendings["home"] = 1234
	shouldBeSpendings["cat"] = 55.61
	shouldBeSpendings["wine"] = 73.31
	spendings := calcSpendings(date, e, false)
	if !reflect.DeepEqual(spendings, shouldBeSpendings) {
		t.Error("calcSpendings() no details spendings result is unexpected:")
		t.Errorf("got %v, wanted %v", spendings, shouldBeSpendings)
	}

	// with details
	shouldBeSpendings = make(map[string]float64)
	shouldBeSpendings["rent"] = 1234
	shouldBeSpendings["cat food shop"] = 55.61
	shouldBeSpendings["saq"] = 73.31

	spendings = calcSpendings(date, e, true)
	if !reflect.DeepEqual(spendings, shouldBeSpendings) {
		t.Error("calcSpendings() w/ details spendings result is unexpected:")
		t.Errorf("got %v, wanted %v", spendings, shouldBeSpendings)
	}
}
