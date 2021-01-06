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

	if !reflect.DeepEqual(categories, shouldBe) {
		t.Errorf("readCategoriesFile() failed: got %v, wanted %v",
			categories, shouldBe)
	}
}

func TestReadMonthlyFile(t *testing.T) {
	entries, err := readMonthlyFile("testdata/december-20.yml")
	if err != nil {
		t.Error("readCategoriesFile() failed")
		t.Fatal(err)
	}
	shouldBe := make(map[string][]moneyExchange, 3)
	date, err := time.Parse("2006-01-02", "2020-12-25")
	if err != nil {
		t.Error("time.Parse in readCategoriesFile() failed")
		t.Fatal(err)
	}
	shouldBe["earnings"] = append(shouldBe["earnings"],
		moneyExchange{4321.0, date, "Company", ""})

	date, err = time.Parse("2006-01-02", "2020-12-01")
	if err != nil {
		t.Error("time.Parse in readCategoriesFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{1234.0, date, "rent", ""})

	date, err = time.Parse("2006-01-02", "2020-12-25")
	if err != nil {
		t.Error("time.Parse in readCategoriesFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{42.24, date, "cat food shop", ""})

	if !reflect.DeepEqual(entries, shouldBe) {
		t.Errorf("readCategoriesFile() failed: got %v, wanted %v",
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
