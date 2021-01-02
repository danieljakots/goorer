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
		moneyExchange{4321.0, date, "Company"})

	date, err = time.Parse("2006-01-02", "2020-12-01")
	if err != nil {
		t.Error("time.Parse in readCategoriesFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{1234.0, date, "Rent"})

	date, err = time.Parse("2006-01-02", "2020-12-25")
	if err != nil {
		t.Error("time.Parse in readCategoriesFile() failed")
		t.Fatal(err)
	}
	shouldBe["spendings"] = append(shouldBe["spendings"],
		moneyExchange{42.24, date, "cat food shop"})

	if !reflect.DeepEqual(entries, shouldBe) {
		t.Errorf("readCategoriesFile() failed: got %v, wanted %v",
			entries, shouldBe)
	}
}
